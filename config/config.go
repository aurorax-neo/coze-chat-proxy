package config

import (
	"github.com/joho/godotenv"
	"os"
)

type Config struct {
	LogLevel  string
	BotConfig string
	AuthToken string
}

var CONFIG *Config

func init() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	CONFIG = &Config{}
	// LOG_LEVEL
	CONFIG.LogLevel = os.Getenv("LOG_LEVEL")
	if CONFIG.LogLevel == "" {
		CONFIG.LogLevel = "info"
	}
	// BOT_CONFIG
	CONFIG.BotConfig = os.Getenv("BOT_CONFIG")
	if CONFIG.BotConfig == "" {
		CONFIG.BotConfig = "bot.json"
	}
	// AUTH_TOKEN
	CONFIG.AuthToken = os.Getenv("AUTH_TOKEN")
	if CONFIG.AuthToken == "" {
		panic("AUTH_TOKEN is not set")
	}
}
