package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port           string
	LastfmApiKey   string
	DiscordWebhook string
}

func Load() *Config {
	GoEnv := os.Getenv("GO_ENV")

	if GoEnv != "PROD" && GoEnv != "production" || GoEnv == "" {
		if err := godotenv.Load(".env.development"); err != nil {
			panic("Error loading .env.development file")
		}
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	return &Config{
		Port:           port,
		LastfmApiKey:   os.Getenv("LASTFM_API_KEY"),
		DiscordWebhook: os.Getenv("DISCORD_WEBHOOK"),
	}
}
