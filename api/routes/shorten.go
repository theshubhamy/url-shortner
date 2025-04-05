package routes

import (
	"strconv"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/theshubhamy/urlshortner/config"
	"github.com/theshubhamy/urlshortner/database"
	"github.com/theshubhamy/urlshortner/helpers"
)

type request struct {
	Url         string        `json:"url"`
	CustomShort string        `json:"customShort"`
	Expirey     time.Duration `json:"expirey"`
}
type response struct {
	Url            string        `json:"url"`
	CustomShort    string        `json:"customShort"`
	Expirey        time.Duration `json:"expirey"`
	XRateLimitRest time.Duration `json:"XRateLimitRest"`
	XRateRemaining int           `json:"XRateRemaining"`
}

func ShortenURL(ctx *fiber.Ctx) error {
	body := new(request)

	err := ctx.BodyParser(&body)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Can't parse Json request."})
	}
	// implement rate limiting
	// redis client
	rdbRate := database.CreateClient(1)
	defer rdbRate.Close()
	value, err := rdbRate.Get(database.Ctx, ctx.IP()).Result()

	if err == redis.Nil {
		_ = rdbRate.Set(database.Ctx, ctx.IP(), config.ApiQuota, 30*60*time.Second).Err()
	} else {

		// value, _ = rdbRate.Get(database.Ctx, ctx.IP()).Result()
		valInt, _ := strconv.Atoi(value)
		if valInt <= 0 {
			limit, _ := rdbRate.TTL(database.Ctx, ctx.IP()).Result()
			return ctx.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{"error": "Rate limit exceeded", "rate_limit_rest": limit / time.Nanosecond / time.Minute})
		}
	}

	//check Request Input is an real url

	if !govalidator.IsURL(body.Url) {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid Url!"})
	}

	// check domain error
	if !helpers.RemoveDomainError(body.Url) {
		return ctx.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{"error": "Resquest can't be processed!"})

	}

	// enforse https,ssl
	body.Url = helpers.EnforseHTTP(body.Url)
	// custom url

	var id string
	if body.CustomShort == "" {
		id = uuid.NewString()[:6]
	} else {
		id = body.CustomShort
	}
	rdbURL := database.CreateClient(0)
	defer rdbURL.Close()
	val, _ := rdbURL.Get(database.Ctx, id).Result()
	if val != "" {
		return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "URL custom short already in use!"})
	}
	if body.Expirey == 0 {
		body.Expirey = 24
	}
	err = rdbURL.Set(database.Ctx, id, body.Url, body.Expirey*3600*time.Second).Err()

	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Server error!"})
	}
	res := response{
		Url:            body.Url,
		CustomShort:    "",
		Expirey:        body.Expirey,
		XRateLimitRest: 10,
		XRateRemaining: 30,
	}

	rdbRate.Decr(database.Ctx, ctx.IP())

	value, _ = rdbRate.Get(database.Ctx, ctx.IP()).Result()
	res.XRateRemaining, _ = strconv.Atoi(value)
	ttl, _ := rdbRate.TTL(database.Ctx, ctx.IP()).Result()
	res.XRateLimitRest = ttl / time.Nanosecond / time.Minute
	res.CustomShort = config.Domain + "/" + id

	return ctx.Status(fiber.StatusOK).JSON(res)
}
