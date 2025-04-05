package main

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/theshubhamy/urlshortner/config"
	"github.com/theshubhamy/urlshortner/routes"
)

func main() {
	fmt.Println("Welcome to the URL Shortener API!")
	config.LoadConfig()

	app := fiber.New()
	app.Use(logger.New())
	setupRoutes(app)
	log.Fatal(app.Listen(config.Port))

}

func setupRoutes(app *fiber.App) {
	app.Get("/:url", routes.ResolveURL)
	app.Post("/api/V1", routes.ShortenURL)
}
