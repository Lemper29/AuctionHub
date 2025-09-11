include .env

.PHONY: gen
gen:
	protoc --go_out=. --go-grpc_out=. proto/auction.proto
	
.PHONY: run
run:
	go run ./auction-service/cmd/server

.PHONY: migrate-up
migrate-up:
	migrate -database "postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable" -path ./auction-service/internal/storage/migrations up

.PHONY: migrate-down
migrate-down:
	migrate -database "postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable" -path ./auction-service/internal/storage/migrations down
