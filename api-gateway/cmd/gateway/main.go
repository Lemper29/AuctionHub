package main

import (
	"context"
	"log"
	"net/http"

	pb "github.com/Lemper29/auction/gen/auction"
	"github.com/gorilla/mux"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	router := mux.NewRouter()

	ctx := context.Background()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	gwMux := runtime.NewServeMux()
	err := pb.RegisterAuctionServiceHandlerFromEndpoint(ctx, gwMux, "localhost:8080", opts)
	if err != nil {
		log.Fatalf("Failed to register gRPC gateway: %v", err)
	}

	router.PathPrefix("/").Handler(gwMux)

	log.Println("Starting server on :8081")
	log.Fatal(http.ListenAndServe(":8081", router))
}
