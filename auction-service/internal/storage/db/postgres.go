package db

import (
	"context"
	"github/auction/internal/storage"
	"github/auction/pkg/models"
	"time"

	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PostgresStorage struct {
	db *gorm.DB
}

func NewPostgresDB(dsn postgres.Config) (storage.Storage, error) {
	db, err := gorm.Open(postgres.New(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	return &PostgresStorage{db: db}, nil
}

func (p *PostgresStorage) CreateLot(ctx context.Context, createLot *models.CreateLotRequest) (*models.Lot, error) {
	id := uuid.New().String()

	// Вычисляем время окончания
	endTime := time.Now().Add(time.Duration(createLot.Duration_minute) * time.Minute)

	lot := &models.Lot{
		Id:             id,
		Name:           createLot.Name,
		Description:    createLot.Description,
		Start_price:    createLot.Start_price,
		Current_price:  createLot.Start_price,
		Current_winner: "",
		Status:         "ACTIVE",
		End_time_unix:  endTime.Unix(), // Сохраняем Unix timestamp
	}

	if err := p.db.WithContext(ctx).Create(lot).Error; err != nil {
		return nil, err
	}

	return lot, nil
}

func (p *PostgresStorage) GetLot(ctx context.Context, getLot *models.GetLotRequest) (*models.GetLotResponse, error) {
	var lot models.Lot
	err := p.db.WithContext(ctx).Where("id = ?", getLot.Lot_id).First(&lot).Error
	if err != nil {
		return nil, err
	}

	return &models.GetLotResponse{Lot: lot}, nil
}

func (p *PostgresStorage) PlaceBid(ctx context.Context, placeBid *models.PlaceBidRequest) (*models.PlaceBidResponse, error) {
	var lot models.Lot
	err := p.db.WithContext(ctx).Where("id = ?", placeBid.Lot_id).First(&lot).Error
	if err != nil {
		return &models.PlaceBidResponse{
			Success: false,
			Message: "Лот не найден",
		}, nil
	}

	if lot.Status != "ACTIVE" {
		return &models.PlaceBidResponse{
			Success: false,
			Message: "Аукцион для этого лота завершен",
		}, nil
	}

	if placeBid.Amount <= lot.Current_price {
		return &models.PlaceBidResponse{
			Success: false,
			Message: "Ставка должна быть выше текущей цены",
		}, nil
	}

	if lot.End_time_unix <= time.Now().Unix() {
		lot.Status = "COMPLETED"
		p.db.Save(&lot)
		return &models.PlaceBidResponse{
			Success: false,
			Message: "Аукцион завершен",
		}, nil
	}

	bidID := uuid.New().String()

	bid := &models.Bid{
		ID:             bidID, // Используем сгенерированный ID
		LotId:          placeBid.Lot_id,
		UserId:         placeBid.User_id,
		Amount:         placeBid.Amount,
		Timestamp_unix: time.Now().Unix(),
	}

	if err := p.db.WithContext(ctx).Create(bid).Error; err != nil {
		return &models.PlaceBidResponse{
			Success: false,
			Message: "Ошибка при сохранении ставки",
		}, err
	}

	lot.Current_price = placeBid.Amount
	lot.Current_winner = placeBid.User_id

	if err := p.db.WithContext(ctx).Save(&lot).Error; err != nil {
		return &models.PlaceBidResponse{
			Success: false,
			Message: "Ошибка при обновлении лота",
		}, err
	}

	return &models.PlaceBidResponse{
		Success:     true,
		Message:     "Ставка принята",
		Updated_lot: lot,
	}, nil
}
