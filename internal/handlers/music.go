package handlers

import (
	"time"

	"github.com/Zigl3ur/api-portfolio/internal/config"
	"github.com/Zigl3ur/api-portfolio/internal/lastfm"
	"github.com/gofiber/fiber/v3"
)

type cacheData struct {
	data *lastfm.FormatedData
	time time.Time
}

var cache = &cacheData{
	data: &lastfm.FormatedData{
		IsListening: false,
	},
	time: time.Time{},
}

type MusicHandler struct {
	cfg *config.Config
}

func NewMusicHandler(cfg *config.Config) *MusicHandler {
	return &MusicHandler{cfg: cfg}
}

func (h *MusicHandler) Handler(c fiber.Ctx) error {
	start := time.Now()

	if !cache.time.IsZero() && time.Since(cache.time) < 30*time.Second {
		c.Set("X-Cached", "true")
		return c.JSON(cache.data)
	}

	data, err := lastfm.MusicHandler(h.cfg.LastfmApiKey)
	if err != nil {
		return err
	}

	cache.data = data
	cache.time = start

	return c.JSON(data)
}
