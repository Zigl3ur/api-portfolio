package handlers

import (
	"time"

	"github.com/Zigl3ur/api-portfolio/internal/lastfm"
	"github.com/gofiber/fiber/v2"
)

type cacheData struct {
	data      *lastfm.FormatedData
	timestamp int64
}

var cache = &cacheData{
	data: &lastfm.FormatedData{
		IsListenning: false,
	},
	timestamp: 0,
}

func MusicHandler(c *fiber.Ctx, apiKey string) error {

	start := time.Now()

	if cache.timestamp != 0 && time.Since(time.Unix(cache.timestamp, 0)) < 30*time.Second {
		c.Set("X-Cached", "true")
		c.JSON(cache.data)
		return nil
	}

	data, err := lastfm.MusicHandler(apiKey)
	if err != nil {
		return err
	}

	cache.data = data
	cache.timestamp = start.Unix()

	c.JSON(data)
	return nil
}
