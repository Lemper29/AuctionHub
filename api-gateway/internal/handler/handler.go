package handler

import (
	"io"
	"log"
	"log/slog"
	"net/http"

	"github.com/Lemper29/api-gateway/internal/config"
	"github.com/Lemper29/api-gateway/internal/logger"
	"github.com/Lemper29/api-gateway/internal/utils"
	"github.com/Lemper29/api-gateway/pkg/models"
	pb "github.com/Lemper29/auction/gen/auction"
	"github.com/gorilla/mux"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Handler struct {
	auctionClient pb.AuctionServiceClient
	logger        *slog.Logger
}

func NewHandler() *Handler {
	conn, err := grpc.Dial(config.Envs.AddressAuctionService,
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to auction service: %v", err)
	}

	appLogger := logger.New(config.Envs.Env, config.Envs.LogLevel)
	serverLogger := appLogger.With(
		"service", "api-gateway",
		"component", "http-handler",
	)

	return &Handler{
		auctionClient: pb.NewAuctionServiceClient(conn),
		logger:        serverLogger,
	}
}

func (h *Handler) CreateLot(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	h.logger.InfoContext(ctx, "CreateLot request started",
		"method", r.Method,
		"path", r.URL.Path,
		"remote_addr", r.RemoteAddr,
	)

	var payload models.CreateLotRequest

	if err := utils.ParseJSON(r, &payload); err != nil {
		h.logger.WarnContext(ctx, "Invalid JSON payload", "error", err.Error())
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	h.logger.InfoContext(ctx, "Creating lot via API",
		"name", payload.Name,
		"start_price", payload.StartPrice,
		"duration_minute", payload.DurationMinute,
	)

	grpcReq := &pb.CreateLotRequest{
		Name:           payload.Name,
		Description:    payload.Description,
		StartPrice:     payload.StartPrice,
		DurationMinute: payload.DurationMinute,
	}

	res, err := h.auctionClient.CreateLot(ctx, grpcReq)
	if err != nil {
		h.logger.ErrorContext(ctx, "Failed to create lot via gRPC",
			"error", err.Error(),
			"name", payload.Name,
		)
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	h.logger.InfoContext(ctx, "Lot created successfully",
		"lot_id", res.Lot.Id,
		"name", res.Lot.Name,
		"status", res.Lot.Status,
	)
	utils.WriteJSON(w, http.StatusCreated, res)
}

func (h *Handler) GetLot(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	id := vars["lot_id"]

	h.logger.InfoContext(ctx, "GetLot request",
		"method", r.Method,
		"path", r.URL.Path,
		"lot_id", id,
		"remote_addr", r.RemoteAddr,
	)

	grpcReq := &pb.GetLotRequest{
		LotId: id,
	}

	res, err := h.auctionClient.GetLot(ctx, grpcReq)
	if err != nil {
		h.logger.ErrorContext(ctx, "Failed to get lot via gRPC",
			"lot_id", id,
			"error", err.Error(),
		)
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	h.logger.DebugContext(ctx, "Lot retrieved successfully",
		"lot_id", id,
		"current_price", res.Lot.CurrentPrice,
		"status", res.Lot.Status,
	)
	utils.WriteJSON(w, http.StatusOK, res)
}

func (h *Handler) PlaceBid(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	id := vars["lot_id"]

	h.logger.InfoContext(ctx, "PlaceBid request started",
		"method", r.Method,
		"path", r.URL.Path,
		"lot_id", id,
		"remote_addr", r.RemoteAddr,
	)

	var payload models.PlaceBidRequest
	if err := utils.ParseJSON(r, &payload); err != nil {
		h.logger.WarnContext(ctx, "Invalid JSON payload for PlaceBid",
			"lot_id", id,
			"error", err.Error(),
		)
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	h.logger.InfoContext(ctx, "Processing bid",
		"lot_id", id,
		"user_id", payload.User_id,
		"amount", payload.Amount,
	)

	grpcReq := &pb.PlaceBidRequest{
		LotId:  id,
		UserId: payload.User_id,
		Amount: payload.Amount,
	}

	res, err := h.auctionClient.PlaceBid(ctx, grpcReq)
	if err != nil {
		h.logger.ErrorContext(ctx, "Failed to place bid via gRPC",
			"lot_id", id,
			"user_id", payload.User_id,
			"amount", payload.Amount,
			"error", err.Error(),
		)
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	if res.Success {
		h.logger.InfoContext(ctx, "Bid accepted",
			"lot_id", id,
			"user_id", payload.User_id,
			"amount", payload.Amount,
			"new_price", res.UpdatedLot.CurrentPrice,
			"winner", res.UpdatedLot.CurrentWinner,
		)
	} else {
		h.logger.WarnContext(ctx, "Bid rejected",
			"lot_id", id,
			"user_id", payload.User_id,
			"amount", payload.Amount,
			"reason", res.Message,
			"current_price", res.UpdatedLot.CurrentPrice,
		)
	}

	utils.WriteJSON(w, http.StatusOK, res)
}

func (h *Handler) SubscribeToLot(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	id := vars["lot_id"]

	h.logger.InfoContext(ctx, "SubscribeToLot request started",
		"method", r.Method,
		"path", r.URL.Path,
		"lot_id", id,
		"remote_addr", r.RemoteAddr,
	)

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	grpcReq := &pb.SubscribeToLotRequest{
		LotId: id,
	}

	stream, err := h.auctionClient.SubscribeToLot(ctx, grpcReq)
	if err != nil {
		h.logger.ErrorContext(ctx, "Failed to subscribe to lot",
			"lot_id", id,
			"error", err.Error(),
		)
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	h.logger.InfoContext(ctx, "Subscription established", "lot_id", id)

	messageCount := 0
	flusher, ok := w.(http.Flusher)
	if !ok {
		h.logger.ErrorContext(ctx, "Streaming not supported")
		utils.WriteError(w, http.StatusInternalServerError,
			http.ErrNotSupported)
		return
	}

	for {
		subscribeToLot, err := stream.Recv()
		if err == io.EOF {
			h.logger.InfoContext(ctx, "Subscription ended by server",
				"lot_id", id,
				"total_messages", messageCount,
			)
			break
		}
		if err != nil {
			h.logger.ErrorContext(ctx, "Error receiving stream message",
				"lot_id", id,
				"error", err.Error(),
				"message_count", messageCount,
			)
			break
		}

		if subscribeToLot.Lot.Status == "COMPLETED" {
			h.logger.InfoContext(ctx, "Auction completed via subscription",
				"lot_id", id,
				"winner", subscribeToLot.Lot.CurrentWinner,
				"final_price", subscribeToLot.Lot.CurrentPrice,
			)
		}

		messageCount++
		if messageCount%10 == 0 {
			h.logger.DebugContext(ctx, "Sending subscription update",
				"lot_id", id,
				"message_count", messageCount,
				"current_price", subscribeToLot.Lot.CurrentPrice,
			)
		}

		if err := utils.WriteJSON(w, http.StatusOK, subscribeToLot); err != nil {
			h.logger.ErrorContext(ctx, "Failed to write JSON response",
				"lot_id", id,
				"error", err.Error(),
			)
			break
		}
		flusher.Flush()
	}

	h.logger.InfoContext(ctx, "Subscription finished",
		"lot_id", id,
		"total_messages_sent", messageCount,
	)
}
