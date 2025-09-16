package main

import (
	"log"

	"github.com/Lemper29/auction-service/internal/config"
	"github.com/Lemper29/auction-service/internal/server"
	"github.com/Lemper29/auction-service/internal/storage/db"
	"gorm.io/driver/postgres"
)

func main() {
	dsn := postgres.Config{
		DSN:                  config.Envs.DSN,
		PreferSimpleProtocol: true,
	}

	storage, err := db.NewPostgresDB(dsn)
	if err != nil {
		log.Fatalf("Database err: %v", err)
	}

	serve := server.NewGrpcServer(":8080", storage)

	if err := serve.Start(); err != nil {
		log.Fatalf("Server err: %v", err)
	}
}
