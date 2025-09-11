package service

import (
	"context"
	pb "github/auction/gen/proto"
	"github/auction/internal/storage"
	"github/auction/pkg/models"

	"log"
	"time"
)

type LotService struct {
	repo storage.Storage
}

func NewLotService(repo storage.Storage) *LotService {
	return &LotService{
		repo: repo,
	}
}

func (l *LotService) CreateLot(ctx context.Context, createLot *pb.CreateLotRequest) (*pb.CreateLotResponse, error) {
	lot := &models.CreateLotRequest{
		Name:            createLot.Name,
		Description:     createLot.Description,
		Start_price:     createLot.StartPrice,
		Duration_minute: createLot.DurationMinute,
	}

	createdLot, err := l.repo.CreateLot(ctx, lot)
	if err != nil {
		return nil, err
	}

	return &pb.CreateLotResponse{
		Lot: convertToPbLot(createdLot),
	}, nil
}

func (l *LotService) GetLot(ctx context.Context, getLot *pb.GetLotRequest) (*pb.GetLotResponse, error) {
	lot := &models.GetLotRequest{
		Lot_id: getLot.LotId,
	}

	res, err := l.repo.GetLot(ctx, lot)
	if err != nil {
		return nil, err
	}

	return &pb.GetLotResponse{
		Lot: convertToPbLot(&res.Lot),
	}, err
}

func (l *LotService) PlaceBid(ctx context.Context, messagePlaceBid *pb.PlaceBidRequest) (*pb.PlaceBidResponse, error) {
	mes := &models.PlaceBidRequest{
		Lot_id:  messagePlaceBid.LotId,
		User_id: messagePlaceBid.UserId,
		Amount:  messagePlaceBid.Amount,
	}

	res, err := l.repo.PlaceBid(ctx, mes)
	if err != nil {
		return nil, err
	}

	updatedLot := convertToPbLot(&res.Updated_lot)

	resMessage := &pb.PlaceBidResponse{
		Success:    res.Success,
		Message:    res.Message,
		UpdatedLot: updatedLot,
	}

	return resMessage, nil
}

func (l *LotService) SubscribeToLot(req *pb.SubscribeToLotRequest, stream pb.AuctionService_SubscribeToLotServer) error {
	ticker := time.NewTicker(time.Second * 3)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			lotResponse, err := l.GetLot(stream.Context(), &pb.GetLotRequest{LotId: req.LotId})
			if err != nil {
				return err
			}

			if lotResponse.Lot.EndTimeUnix <= time.Now().Unix() {
				if lotResponse.Lot.Status == "ACTIVE" {
					lotResponse.Lot.Status = "COMPLETED"
					if err := stream.Send(&pb.SubscribeToLotResponse{Lot: lotResponse.Lot}); err != nil {
						log.Printf("Ошибка отправки: %v", err)
					}
					return nil
				}
			}

			if err := stream.Send(&pb.SubscribeToLotResponse{Lot: lotResponse.Lot}); err != nil {
				log.Printf("Ошибка отправки: %v", err)
				return err
			}

		case <-stream.Context().Done():
			return nil
		}
	}
}

// Вспомогательная функция для конвертации
func convertToPbLot(lot *models.Lot) *pb.Lot {
	return &pb.Lot{
		Id:            lot.Id,
		Name:          lot.Name,
		Description:   lot.Description,
		StartPrice:    lot.Start_price,
		CurrentPrice:  lot.Current_price,
		CurrentWinner: lot.Current_winner,
		Status:        lot.Status,
		EndTimeUnix:   lot.End_time_unix, // Уже int64
	}
}
