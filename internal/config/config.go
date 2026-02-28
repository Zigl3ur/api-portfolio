package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Env            string
	Port           string
	LastfmApiKey   string
	DiscordWebhook string
}

func Load() *Config {
	AppEnv := os.Getenv("APP_ENV")

	if AppEnv == "development" || AppEnv == "" {
		if err := godotenv.Load(".env.development"); err != nil {
			panic("Error loading .env.development file")
		}
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	if os.Getenv("LASTFM_API_KEY") == "" {
		panic("LASTFM_API_KEY environment variable is required")
	}

	if os.Getenv("DISCORD_WEBHOOK") == "" {
		panic("DISCORD_WEBHOOK environment variable is required")
	}

	return &Config{
		Env:            AppEnv,
		Port:           port,
		LastfmApiKey:   os.Getenv("LASTFM_API_KEY"),
		DiscordWebhook: os.Getenv("DISCORD_WEBHOOK"),
	}
}
