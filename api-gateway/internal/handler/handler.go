package handler

import (
	"log"
	"net/http"

	"github.com/Lemper29/api-gateway/internal/utils"
	"github.com/Lemper29/api-gateway/pkg/models"
	pb "github.com/Lemper29/auction/gen/auction"
	"github.com/gorilla/mux"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Handler struct {
	auctionClient pb.AuctionServiceClient
}

func NewHandler() *Handler {
	// Подключаемся к auction-service
	conn, err := grpc.Dial("localhost:8080",
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to auction service: %v", err)
	}

	return &Handler{
		auctionClient: pb.NewAuctionServiceClient(conn),
	}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/create", h.CreateLot).Methods("POST")
	router.HandleFunc("/{id}", h.GetLot).Methods("GET")
	router.HandleFunc("/{id}/bids", h.PlaceBid).Methods("POST")
}

func (h *Handler) CreateLot(w http.ResponseWriter, r *http.Request) {
	var payload models.CreateLotRequest

	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	grpcReq := &pb.CreateLotRequest{
		Name:           payload.Name,
		Description:    payload.Description,
		StartPrice:     payload.Start_price,
		DurationMinute: payload.Duration_minute,
	}

	res, err := h.auctionClient.CreateLot(r.Context(), grpcReq)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, res)
}

func (h *Handler) GetLot(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	grpcReq := &pb.GetLotRequest{
		LotId: id,
	}

	res, err := h.auctionClient.GetLot(r.Context(), grpcReq)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, res)
}

func (h *Handler) PlaceBid(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var payload models.PlaceBidRequest
	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	grpcReq := &pb.PlaceBidRequest{
		LotId:  id,
		UserId: payload.User_id,
		Amount: payload.Amount,
	}

	res, err := h.auctionClient.PlaceBid(r.Context(), grpcReq)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, res)
}
