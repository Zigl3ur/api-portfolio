package handlers

import (
	"time"

	"github.com/Zigl3ur/api-portfolio/internal/cache"
	"github.com/Zigl3ur/api-portfolio/internal/config"
	"github.com/Zigl3ur/api-portfolio/internal/lastfm"
	"github.com/gofiber/fiber/v3"
)

type MusicHandler struct {
	lastfm *lastfm.LastFM

	cache *cache.Cache
}

func NewMusicHandler(cfg *config.Config, cache *cache.Cache) *MusicHandler {
	return &MusicHandler{lastfm: lastfm.NewLastFM(cfg.LastfmApiKey), cache: cache}
}

func (h *MusicHandler) CurrentlyListening(c fiber.Ctx) error {
	cachedData := h.cache.Get("currentlyListening")
	if cachedData != nil {
		c.Set("X-Cached", "true")
		return c.JSON(cachedData)
	}

	data, err := h.lastfm.GetCurrentlyPlaying()
	if err != nil {
		return err
	}

	h.cache.Set("currentlyListening", data, 30*time.Second)

	return c.JSON(data)
}

func (h *MusicHandler) TopAlbums(c fiber.Ctx) error {
	cachedData := h.cache.Get("topAlbums")
	if cachedData != nil {
		c.Set("X-Cached", "true")
		return c.JSON(cachedData)
	}

	data, err := h.lastfm.GetTopAlbums()
	if err != nil {
		return err
	}

	h.cache.Set("topAlbums", data, 1*time.Hour)

	return c.JSON(data)
}

func (h *MusicHandler) AlbumInfo(c fiber.Ctx) error {
	artist := c.Query("artist")
	album := c.Query("album")

	if artist == "" || album == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "artist and album query parameters are required",
		})
	}

	cacheKey := "albumInfo:" + artist + ":" + album
	cachedData := h.cache.Get(cacheKey)
	if cachedData != nil {
		c.Set("X-Cached", "true")
		return c.JSON(cachedData)
	}

	data, err := h.lastfm.GetAlbumInfo(artist, album)
	if err != nil {
		return err
	}

	h.cache.Set(cacheKey, data, 1*time.Hour)

	return c.JSON(data)
}
