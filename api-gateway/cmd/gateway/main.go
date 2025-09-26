package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/Lemper29/api-gateway/internal/config"
	"github.com/Lemper29/api-gateway/internal/logger"
	pb "github.com/Lemper29/auction/gen/auction"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func loggingMiddleware(next http.Handler) http.Handler {
	appLogger := logger.New(config.Envs.Env, config.Envs.LogLevel)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		requestLogger := appLogger.With(
			"method", r.Method,
			"path", r.URL.Path,
			"remote_addr", r.RemoteAddr,
		)

		requestLogger.Info("Request started")

		rw := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		next.ServeHTTP(rw, r)

		duration := time.Since(start)
		requestLogger.Info("Request completed",
			"status", rw.statusCode,
			"duration", duration.String(),
		)
	})
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func main() {
	ctx := context.Background()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	gwMux := runtime.NewServeMux()
	err := pb.RegisterAuctionServiceHandlerFromEndpoint(
		ctx,
		gwMux,
		config.Envs.AddressAuctionService,
		opts,
	)
	if err != nil {
		log.Fatalf("Failed to register gRPC gateway: %v", err)
	}

	loggingMux := loggingMiddleware(gwMux)

	log.Println("Starting server on :" + config.Envs.PortApiGatewayService)
	log.Fatal(http.ListenAndServe(":"+config.Envs.PortApiGatewayService, loggingMux))
}
