package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	PublicHost             string
	PortAuctionService     string
	DBUser                 string
	DBPassword             string
	DBAddress              string
	DBName                 string
	DBPort                 string
	DSN                    string
	JWTSecret              string
	JWTExpirationInSeconds int64
}

var Envs = initConfig()

func initConfig() Config {
	godotenv.Load()

	dbHost := getEnv("DB_HOST", "127.0.0.1")
	dbPort := getEnv("DB_PORT", "5432")
	dbUser := getEnv("DB_USER", "postgres")
	dbPassword := getEnv("DB_PASSWORD", "0000")
	dbName := getEnv("DB_NAME", "postgres")

	return Config{
		PublicHost:         getEnv("PUBLIC_HOST", "http://localhost"),
		PortAuctionService: getEnv("PORT_AUCTION_SERVICE", "8080"),
		DBUser:             dbUser,
		DBPassword:         dbPassword,
		DBAddress:          fmt.Sprintf("%s:%s", dbHost, dbPort),
		DBName:             dbName,
		DBPort:             dbPort,
		DSN: fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			dbHost, dbPort, dbUser, dbPassword, dbName),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
