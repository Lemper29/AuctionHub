package main

import (
	"log"
	"net/http"

	"github.com/Lemper29/api-gateway/internal/handler"
	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()
	subrouter := router.PathPrefix("/api/v1").Subrouter()

	handler := handler.NewHandler()
	handler.RegisterRoutes(subrouter)

	log.Println("Starting server on :8081")
	log.Fatal(http.ListenAndServe(":8081", router))
}
