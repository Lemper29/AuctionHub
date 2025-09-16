PROTO_DIR=protos
GEN_DIR=gen

.PHONY: gen
gen: gen-grpc gen-gateway gen-openapi

.PHONY: gen-grpc
gen-grpc:
	protoc -I./$(PROTO_DIR) \
    -I./$(PROTO_DIR)/googleapis \
    --go_out=./$(GEN_DIR) --go_opt=paths=source_relative \
    --go-grpc_out=./$(GEN_DIR) --go-grpc_opt=paths=source_relative \
    --grpc-gateway_out=./$(GEN_DIR) --grpc-gateway_opt=paths=source_relative \
    --openapiv2_out=./$(GEN_DIR) \
    ./$(PROTO_DIR)/auction/auction.proto

.PHONY: gen-gateway
gen-gateway:
	protoc -I=./$(PROTO_DIR) \
    -I./$(PROTO_DIR)/googleapis \
	--grpc-gateway_out=./$(GEN_DIR) \
	--grpc-gateway_opt=logtostderr=true \
	--grpc-gateway_opt=paths=source_relative \
	--grpc-gateway_opt=generate_unbound_methods=true \
	./$(PROTO_DIR)/auction/auction.proto

.PHONY: gen-openapi
gen-openapi:
	protoc -I=./$(PROTO_DIR) \
    -I./$(PROTO_DIR)/googleapis \
	--openapiv2_out=./$(GEN_DIR) \
	--openapiv2_opt=logtostderr=true \
	--openapiv2_opt=generate_unbound_methods=true \
	./$(PROTO_DIR)/auction/auction.proto

.PHONY: run-server
run-server:
	cd auction-service && go run ./cmd/server

.PHONY: run-gateway
run-gateway:
	cd api-gateway && go run ./cmd/gateway

.PHONY: run-api
run-api:
	cd api-gateway && go run ./cmd/main

.PHONY: migrate-up
migrate-up:
	migrate -database "postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable" -path ./auction-service/internal/storage/migrations up

.PHONY: migrate-down
migrate-down:
	migrate -database "postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable" -path ./auction-service/internal/storage/migrations down

.PHONY: test
test:
	go test ./... -v

.PHONY: build
build:
	go build -o bin/auction-server ./auction-service/cmd/server
	go build -o bin/api-gateway ./api-gateway/cmd/gateway

.PHONY: docker-build
docker-build:
	docker build -t auction-service -f auction-service/Dockerfile .
	docker build -t api-gateway -f api-gateway/Dockerfile .

.PHONY: help
help:
	@echo "Available commands:"
	@echo "  make gen          - Generate all code (gRPC + Gateway + OpenAPI)"
	@echo "  make gen-grpc     - Generate only gRPC code"
	@echo "  make gen-gateway  - Generate gRPC-Gateway code"
	@echo "  make gen-openapi  - Generate OpenAPI documentation"
	@echo "  make run-server   - Run auction service (gRPC server)"
	@echo "  make run-gateway  - Run API Gateway (gRPC-Gateway)"
	@echo "  make run-api      - Run custom API Gateway"
	@echo "  make migrate-up   - Run database migrations"
	@echo "  make migrate-down - Rollback database migrations"
	@echo "  make test         - Run tests"
	@echo "  make build        - Build binaries"