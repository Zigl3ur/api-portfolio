package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Zigl3ur/api-portfolio/internal/handlers"
	"github.com/gofiber/fiber/v2"
)

func main() {

	app := fiber.New()

	port, ok := os.LookupEnv("PORT")
	if !ok {
		fmt.Println("env PORT not set defaulting to 8080")
		port = "8080"
	}
	apiKey, ok := os.LookupEnv("LASTFM_API_KEY")
	if !ok {
		log.Fatal("env LASTFM_API_KEY not set")
	}

	app.Get("/music", func(c *fiber.Ctx) error {
		return handlers.MusicHandler(c, apiKey)
	})

	app.Use(func(c *fiber.Ctx) error {
		return c.SendStatus(404)
	})

	log.Fatal(app.Listen(fmt.Sprintf(":%s", port)))
}
