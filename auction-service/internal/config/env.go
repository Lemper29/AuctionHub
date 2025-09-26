package config

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	PublicHost         string
	PortAuctionService string
	DBUser             string
	DBPassword         string
	DBAddress          string
	DBName             string
	DBPort             string
	DSN                string
	LogLevel           slog.Level
}

var Envs = InitConfig()

func InitConfig() *Config {
	godotenv.Load(".env")

	dbHost := getEnv("DB_HOST", "127.0.0.1")
	dbPort := getEnv("DB_PORT", "5432")
	dbUser := getEnv("DB_USER", "postgres")
	dbPassword := getEnv("DB_PASSWORD", "0000")
	dbName := getEnv("DB_NAME", "postgres")

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

	return &Config{
		PublicHost:         getEnv("PUBLIC_HOST", "http://localhost"),
		PortAuctionService: getEnv("PORT_AUCTION_SERVICE", "8080"),
		DBUser:             dbUser,
		DBPassword:         dbPassword,
		DBAddress:          fmt.Sprintf("%s:%s", dbHost, dbPort),
		DBName:             dbName,
		DBPort:             dbPort,
		DSN: fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			dbHost, dbPort, dbUser, dbPassword, dbName),
		LogLevel: logLevel,
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
