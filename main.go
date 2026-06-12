package main

import (
	"context"
	"fmt"
	"log"

	"github.com/Zigl3ur/api-portfolio/internal/cache"
	"github.com/Zigl3ur/api-portfolio/internal/config"
	"github.com/Zigl3ur/api-portfolio/internal/handlers"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
)

func main() {
	appCtx, cancel := context.WithCancel(context.Background())
	cfg := config.Load()

	cache := cache.NewCache(appCtx)

	app := fiber.New(fiber.Config{
		AppName:    "api-portfolio",
		TrustProxy: true,
		TrustProxyConfig: fiber.TrustProxyConfig{
			Proxies: []string{"192.168.1.188"},
		},
		ProxyHeader: "Cf-Connecting-Ip",
	})

	app.Hooks().OnPostShutdown(func(e error) error {
		cancel()
		return nil
	})

	allowedOrigins := []string{
		"https://eden.douru.fr",
	}

	if cfg.Env == "development" || cfg.Env == "" {
		allowedOrigins = append(allowedOrigins, "http://localhost:3000")
	}

	app.Use(cors.New(cors.Config{
		AllowOrigins: allowedOrigins,
		AllowHeaders: []string{"Origin", "Content-Type", "Accept"},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
	}))

	musicHandler := handlers.NewMusicHandler(cfg, cache)
	messageHandler := handlers.NewMessageHandler(cfg)

	musicGroup := app.Group("/music")
	musicGroup.Get("/currently-listening", musicHandler.CurrentlyListening)
	musicGroup.Get("/top-albums", musicHandler.TopAlbums)
	musicGroup.Get("/album-info", musicHandler.AlbumInfo)

	app.Post("/message", handlers.MessageLimiter, messageHandler.Handler)

	app.Use(func(c fiber.Ctx) error {
		return c.SendStatus(fiber.StatusNotFound)
	})

	log.Fatal(app.Listen(fmt.Sprintf(":%s", cfg.Port)))
}
