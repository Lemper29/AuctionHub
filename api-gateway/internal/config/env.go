package config

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	PublicHost            string
	PortApiGatewayService string
	PortAuctionService    string
	AddressAuctionService string
	Env                   string
	LogLevel              slog.Level
}

var Envs = initConfig()

func initConfig() Config {
	godotenv.Load(".env")

	portAuctionService := getEnv("PORT_AUCTION_SERVICE", "8080")
	portApiGatewayService := getEnv("PORT_API_GATEWAY_SERVICE", "8081")
	publicHost := getEnv("PUBLIC_HOST", "localhost")

	env := getEnv("APP_ENV", "development")
	logLevelStr := getEnv("LOG_LEVEL", "debug")

	var logLevel slog.Level
	switch logLevelStr {
	case "debug":
		logLevel = slog.LevelDebug
	case "info":
		logLevel = slog.LevelInfo
	case "warn":
		logLevel = slog.LevelWarn
	case "error":
		logLevel = slog.LevelError
	default:
		logLevel = slog.LevelDebug
	}

	return Config{
		PublicHost:            publicHost,
		PortApiGatewayService: portApiGatewayService,
		PortAuctionService:    portAuctionService,
		AddressAuctionService: fmt.Sprintf("%s:%s", publicHost, portAuctionService),
		Env:                   env,
		LogLevel:              logLevel,
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
