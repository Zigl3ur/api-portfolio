package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Zigl3ur/api-portfolio/internal/handlers"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
)

func main() {

	app := fiber.New(fiber.Config{
		TrustedProxies:          []string{"192.168.1.188"},
		EnableTrustedProxyCheck: true,
		ProxyHeader:             "Cf-Connecting-Ip",
	})

	port, ok := os.LookupEnv("PORT")
	if !ok {
		fmt.Println("env PORT not set defaulting to 8080")
		port = "8080"
	}
	apiKey, ok := os.LookupEnv("LASTFM_API_KEY")
	if !ok {
		log.Fatal("env LASTFM_API_KEY not set")
	}
	webhook, ok := os.LookupEnv("DISCORD_WEBHOOK")
	if !ok {
		log.Fatal("env DISCORD_WEBHOOK not set")
	}

	app.Use(cors.New(cors.Config{
		AllowOrigins: "https://eden.douru.fr, http://localhost:3000",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	app.Get("/music", func(c *fiber.Ctx) error {
		return handlers.MusicHandler(c, apiKey)
	})

	messageLimiter := limiter.New(limiter.Config{
		Max:        2,
		Expiration: time.Hour * 24 * 7, // 1 week
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error": "Rate limit reached. Try again later.",
			})
		},
	})

	app.Post("/message", messageLimiter, func(c *fiber.Ctx) error {
		return handlers.MessageHandler(c, webhook)
	})

	app.Use(func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusNotFound)
	})

	log.Fatal(app.Listen(fmt.Sprintf(":%s", port)))
}
