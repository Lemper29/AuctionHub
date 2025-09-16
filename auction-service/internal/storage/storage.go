package storage

import (
	"context"

	"github.com/Lemper29/auction-service/pkg/models"
)

type Storage interface {
	CreateLot(ctx context.Context, req *models.CreateLotRequest) (*models.Lot, error)
	GetLot(ctx context.Context, req *models.GetLotRequest) (*models.GetLotResponse, error)
	PlaceBid(ctx context.Context, req *models.PlaceBidRequest) (*models.PlaceBidResponse, error)
}
