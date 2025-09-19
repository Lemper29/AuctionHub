package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	PublicHost            string
	PortApiGatewayService string
	PortAuctionService    string
	AddressAuctionService string
}

var Envs = initConfig()

func initConfig() Config {
	godotenv.Load()

	portAuctionService := getEnv("PORT_AUCTION_SERVICE", "8080")
	portApiGatewayService := getEnv("PORT_API_GATEWAY_SERVICE", "8081")
	publicHost := getEnv("PUBLIC_HOST", "localhost")

	return Config{
		PublicHost:            publicHost,
		PortApiGatewayService: portApiGatewayService,
		AddressAuctionService: fmt.Sprintf("%s:%s", publicHost, portAuctionService),
		PortAuctionService:    portAuctionService,
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
