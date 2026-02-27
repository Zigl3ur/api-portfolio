package handlers

import (
	"time"

	"github.com/Zigl3ur/api-portfolio/internal/lastfm"
	"github.com/gofiber/fiber/v2"
)

type cacheData struct {
	data *lastfm.FormatedData
	time time.Time
}

var cache = &cacheData{
	data: &lastfm.FormatedData{
		IsListenning: false,
	},
	time: time.Time{},
}

func MusicHandler(c *fiber.Ctx, apiKey string) error {

	start := time.Now()
	c.Accepts("application/json")

	if !cache.time.IsZero() && time.Since(cache.time) < 30*time.Second {
		c.Set("X-Cached", "true")
		return c.JSON(cache.data)
	}

	data, err := lastfm.MusicHandler(apiKey)
	if err != nil {
		return err
	}

	cache.data = data
	cache.time = start

	return c.JSON(data)
}
