package main

import (
	"log"

	"github.com/Lemper29/auction-service/internal/config"
	"github.com/Lemper29/auction-service/internal/logger"
	"github.com/Lemper29/auction-service/internal/server"
	"github.com/Lemper29/auction-service/internal/storage/db"
	"gorm.io/driver/postgres"
)

func main() {
	appLogger := logger.New("auction-service", config.Envs.LogLevel)
	appLogger.Info("Starting auction service", "version", "1.0.0")

	dsn := postgres.Config{
		DSN:                  config.Envs.DSN,
		PreferSimpleProtocol: true,
	}

	storage, err := db.NewPostgresDB(dsn)
	if err != nil {
		appLogger.Error("Database connection failed", "error", err)
		log.Fatalf("Database err: %v", err)
	}

	appLogger.Info("Database connection established")

	serve := server.NewGrpcServer(":"+config.Envs.PortAuctionService, storage, appLogger)

	appLogger.Info("Server starting", "port", config.Envs.PortAuctionService)
	if err := serve.Start(); err != nil {
		appLogger.Error("Server failed to start", "error", err)
		log.Fatalf("Server err: %v", err)
	}
}
