package config

import (
	"github.com/joho/godotenv"
	"os"
)

type Config struct {
	LogLevel   string
	ServerPort string
	BotConfig  string
	AuthToken  string
}

var CONFIG *Config

func init() {
	_ = godotenv.Load()
	CONFIG = &Config{}
	// LOG_LEVEL
	CONFIG.LogLevel = os.Getenv("LOG_LEVEL")
	if CONFIG.LogLevel == "" {
		CONFIG.LogLevel = "info"
	}
	// SERVER_PORT
	CONFIG.ServerPort = os.Getenv("SERVER_PORT")
	if CONFIG.ServerPort == "" {
		CONFIG.ServerPort = "8080"
	}
	// BOT_CONFIG
	CONFIG.BotConfig = os.Getenv("BOT_CONFIG")
	if CONFIG.BotConfig == "" {
		CONFIG.BotConfig = "bot.json"
	}
	// AUTH_TOKEN
	CONFIG.AuthToken = os.Getenv("AUTH_TOKEN")
	if CONFIG.AuthToken == "" {
		CONFIG.AuthToken = "1234567890:ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	}
}
