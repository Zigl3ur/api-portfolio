package main

import (
	"fmt"
	"log"

	"github.com/Zigl3ur/api-portfolio/internal/config"
	"github.com/Zigl3ur/api-portfolio/internal/handlers"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
)

func main() {
	cfg := config.Load()

	app := fiber.New(fiber.Config{
		AppName:    "api-portfolio",
		TrustProxy: true,
		TrustProxyConfig: fiber.TrustProxyConfig{
			Proxies: []string{"192.168.1.188"},
		},
		ProxyHeader: "Cf-Connecting-Ip",
	})

	allowedOrigins := []string{
		"https://eden.douru.fr",
	}

	if cfg.Env == "development" {
		allowedOrigins = append(allowedOrigins, "http://localhost:3000")
	}

	app.Use(cors.New(cors.Config{
		AllowOrigins: allowedOrigins,
		AllowHeaders: []string{"Origin", "Content-Type", "Accept"},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
	}))

	musicHandler := handlers.NewMusicHandler(cfg)
	messageHandler := handlers.NewMessageHandler(cfg)

	app.Get("/music", musicHandler.Handler)
	app.Post("/message", handlers.MessageLimiter, messageHandler.Handler)

	app.Use(func(c fiber.Ctx) error {
		return c.SendStatus(fiber.StatusNotFound)
	})

	log.Fatal(app.Listen(fmt.Sprintf(":%s", cfg.Port)))
}
