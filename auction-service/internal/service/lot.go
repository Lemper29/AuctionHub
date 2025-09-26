package service

import (
	"context"
	"log/slog"
	"time"

	"github.com/Lemper29/auction-service/internal/storage"
	"github.com/Lemper29/auction-service/pkg/models"
	pb "github.com/Lemper29/auction/gen/auction"
)

type LotService struct {
	repo   storage.Storage
	logger *slog.Logger
}

func NewLotService(repo storage.Storage, logger *slog.Logger) *LotService {
	return &LotService{
		repo:   repo,
		logger: logger,
	}
}

func (l *LotService) CreateLot(ctx context.Context, createLot *pb.CreateLotRequest) (*pb.CreateLotResponse, error) {
	l.logger.InfoContext(ctx, "Creating lot",
		"name", createLot.Name,
		"start_price", createLot.StartPrice,
	)

	lot := &models.CreateLotRequest{
		Name:           createLot.Name,
		Description:    createLot.Description,
		StartPrice:     createLot.StartPrice,
		DurationMinute: createLot.DurationMinute,
	}

	createdLot, err := l.repo.CreateLot(ctx, lot)
	if err != nil {
		l.logger.ErrorContext(ctx, "Failed to create lot", "error", err)
		return nil, err
	}

	l.logger.InfoContext(ctx, "Lot created successfully", "lot_id", createdLot.Id)
	return &pb.CreateLotResponse{
		Lot: convertToPbLot(createdLot),
	}, nil
}

func (l *LotService) GetLot(ctx context.Context, getLot *pb.GetLotRequest) (*pb.GetLotResponse, error) {
	l.logger.DebugContext(ctx, "Getting lot", "lot_id", getLot.LotId)

	lot := &models.GetLotRequest{
		Lot_id: getLot.LotId,
	}

	res, err := l.repo.GetLot(ctx, lot)
	if err != nil {
		l.logger.ErrorContext(ctx, "Failed to get lot", "lot_id", getLot.LotId, "error", err)
		return nil, err
	}

	l.logger.DebugContext(ctx, "Lot retrieved", "lot_id", getLot.LotId)
	return &pb.GetLotResponse{
		Lot: convertToPbLot(&res.Lot),
	}, nil
}

func (l *LotService) PlaceBid(ctx context.Context, messagePlaceBid *pb.PlaceBidRequest) (*pb.PlaceBidResponse, error) {
	l.logger.InfoContext(ctx, "Processing bid",
		"lot_id", messagePlaceBid.LotId,
		"user_id", messagePlaceBid.UserId,
		"amount", messagePlaceBid.Amount,
	)

	mes := &models.PlaceBidRequest{
		Lot_id:  messagePlaceBid.LotId,
		User_id: messagePlaceBid.UserId,
		Amount:  messagePlaceBid.Amount,
	}

	res, err := l.repo.PlaceBid(ctx, mes)
	if err != nil {
		l.logger.ErrorContext(ctx, "Failed to process bid",
			"lot_id", messagePlaceBid.LotId,
			"error", err,
		)
		return nil, err
	}

	if res.Success {
		l.logger.InfoContext(ctx, "Bid accepted",
			"lot_id", messagePlaceBid.LotId,
			"new_price", res.Updated_lot.CurrentPrice,
			"winner", res.Updated_lot.CurrentWinner,
		)
	} else {
		l.logger.WarnContext(ctx, "Bid rejected",
			"lot_id", messagePlaceBid.LotId,
			"reason", res.Message,
			"current_price", res.Updated_lot.CurrentPrice,
		)
	}

	return &pb.PlaceBidResponse{
		Success:    res.Success,
		Message:    res.Message,
		UpdatedLot: convertToPbLot(&res.Updated_lot),
	}, nil
}

func (l *LotService) SubscribeToLot(req *pb.SubscribeToLotRequest, stream pb.AuctionService_SubscribeToLotServer) error {
	l.logger.InfoContext(stream.Context(), "Starting subscription", "lot_id", req.LotId)

	ticker := time.NewTicker(time.Second * 3)
	defer ticker.Stop()

	updateCount := 0

	for {
		select {
		case <-ticker.C:
			lotResponse, err := l.GetLot(stream.Context(), &pb.GetLotRequest{LotId: req.LotId})
			if err != nil {
				l.logger.ErrorContext(stream.Context(), "Failed to get lot for subscription",
					"lot_id", req.LotId, "error", err)
				return err
			}

			if lotResponse.Lot.EndTimeUnix <= time.Now().Unix() {
				if lotResponse.Lot.Status == "ACTIVE" {
					lotResponse.Lot.Status = "COMPLETED"
					l.logger.InfoContext(stream.Context(), "Auction completed",
						"lot_id", req.LotId,
						"winner", lotResponse.Lot.CurrentWinner,
						"final_price", lotResponse.Lot.CurrentPrice,
					)

					if err := stream.Send(&pb.SubscribeToLotResponse{Lot: lotResponse.Lot}); err != nil {
						l.logger.ErrorContext(stream.Context(), "Failed to send completion update",
							"lot_id", req.LotId, "error", err)
					}
					return nil
				}
			}

			if err := stream.Send(&pb.SubscribeToLotResponse{Lot: lotResponse.Lot}); err != nil {
				l.logger.ErrorContext(stream.Context(), "Failed to send lot update",
					"lot_id", req.LotId, "error", err)
				return err
			}

			updateCount++
			l.logger.DebugContext(stream.Context(), "Sent lot update",
				"lot_id", req.LotId,
				"update_count", updateCount,
				"current_price", lotResponse.Lot.CurrentPrice,
			)

		case <-stream.Context().Done():
			l.logger.InfoContext(stream.Context(), "Subscription ended by client",
				"lot_id", req.LotId,
				"total_updates", updateCount,
			)
			return nil
		}
	}
}

func convertToPbLot(lot *models.Lot) *pb.Lot {
	if lot == nil {
		return &pb.Lot{}
	}

	return &pb.Lot{
		Id:            lot.Id,
		Name:          lot.Name,
		Description:   lot.Description,
		StartPrice:    lot.StartPrice,
		CurrentPrice:  lot.CurrentPrice,
		CurrentWinner: lot.CurrentWinner,
		Status:        lot.Status,
		EndTimeUnix:   lot.EndTimeUnix,
	}
}
