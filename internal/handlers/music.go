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
	return &MusicHandler{lastfm: lastfm.NewLastFM(cfg.LastfmApiKey, cfg.LastfmUsername), cache: cache}
}

func (h *MusicHandler) CurrentlyListening(c fiber.Ctx) error {
	cachedData := h.cache.Get("currentlyListening")
	if cachedData != nil {
		h.cache.SetHeader(c, "currentlyListening")
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
		h.cache.SetHeader(c, "topAlbums")
		return c.JSON(cachedData)
	}

	data, err := h.lastfm.GetTopAlbums()
	if err != nil {
		return err
	}

	h.cache.Set("topAlbums", data, 1*time.Hour)

	return c.JSON(data)
}

func (h *MusicHandler) TopTracks(c fiber.Ctx) error {
	cachedData := h.cache.Get("topTracks")
	if cachedData != nil {
		h.cache.SetHeader(c, "topTracks")
		return c.JSON(cachedData)
	}

	data, err := h.lastfm.GetTopTracks()
	if err != nil {
		return err
	}

	h.cache.Set("topTracks", data, 1*time.Hour)

	return c.JSON(data)
}

func (h *MusicHandler) TopArtists(c fiber.Ctx) error {
	cachedData := h.cache.Get("topArtists")
	if cachedData != nil {
		h.cache.SetHeader(c, "topArtists")
		return c.JSON(cachedData)
	}

	data, err := h.lastfm.GetTopArtists()
	if err != nil {
		return err
	}

	h.cache.Set("topArtists", data, 1*time.Hour)

	return c.JSON(data)
}
