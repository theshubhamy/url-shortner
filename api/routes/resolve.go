package routes

import (
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/theshubhamy/urlshortner/database"
)

func ResolveURL(ctx *fiber.Ctx) error {

	url := ctx.Params("url")

	rdb := database.CreateClient(0)
	defer rdb.Close()
	value, err := rdb.Get(database.Ctx, url).Result()
	if err == redis.Nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "short url not found!"})
	} else if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Server error!"})
	}
	rInr := database.CreateClient(1)
	defer rInr.Close()
	_ = rInr.Incr(database.Ctx, "counter")
	return ctx.Redirect(value, 301)

}
