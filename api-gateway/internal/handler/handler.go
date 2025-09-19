package handler

import (
	"context"
	"io"
	"log"
	"net/http"

	"github.com/Lemper29/api-gateway/internal/config"
	"github.com/Lemper29/api-gateway/internal/utils"
	"github.com/Lemper29/api-gateway/pkg/models"
	pb "github.com/Lemper29/auction/gen/auction"
	"github.com/gorilla/mux"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Handler struct {
	auctionClient pb.AuctionServiceClient
}

func NewHandler() *Handler {
	conn, err := grpc.Dial(config.Envs.AddressAuctionService,
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to auction service: %v", err)
	}

	return &Handler{
		auctionClient: pb.NewAuctionServiceClient(conn),
	}
}

func (h *Handler) RegisterRoutes(ctx context.Context, mux *runtime.ServeMux, opts []grpc.DialOption) {
	err := pb.RegisterAuctionServiceHandlerFromEndpoint(ctx, mux, config.Envs.AddressAuctionService, opts)
	if err != nil {
		log.Fatalf("Failed to register gRPC gateway: %v", err)
	}
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
		StartPrice:     payload.StartPrice,
		DurationMinute: payload.DurationMinute,
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
	id := vars["lot_id"]

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
	id := vars["lot_id"]

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

func (h *Handler) SubscribeToLot(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["lot_id"]

	grpcReq := &pb.SubscribeToLotRequest{
		LotId: id,
	}

	stream, err := h.auctionClient.SubscribeToLot(r.Context(), grpcReq)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
	}

	for {
		subscribeToLot, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("error while receiving todo: %v", err)
		}

		utils.WriteJSON(w, http.StatusOK, subscribeToLot)
	}
}
