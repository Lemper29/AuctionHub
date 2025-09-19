package db

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Lemper29/auction-service/internal/storage"
	"github.com/Lemper29/auction-service/pkg/models"
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
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	log.Println("Database connection established successfully")
	return &PostgresStorage{db: db}, nil
}

func (p *PostgresStorage) CreateLot(ctx context.Context, createLot *models.CreateLotRequest) (*models.Lot, error) {
	log.Printf("CreateLot request: %+v", createLot)

	id := uuid.New().String()
	endTime := time.Now().Add(time.Duration(createLot.DurationMinute) * time.Minute)

	lot := &models.Lot{
		Id:            id,
		Name:          createLot.Name,
		Description:   createLot.Description,
		StartPrice:    createLot.StartPrice,
		CurrentPrice:  createLot.StartPrice,
		CurrentWinner: "",
		Status:        "ACTIVE",
		End_time_unix: endTime.Unix(),
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	log.Printf("Lot object before save: %+v", lot)

	result := p.db.WithContext(ctx).Create(lot)
	if result.Error != nil {
		log.Printf("Error creating lot: %v", result.Error)
		return nil, result.Error
	}

	log.Printf("Rows affected: %d", result.RowsAffected)

	var savedLot models.Lot
	err := p.db.WithContext(ctx).First(&savedLot, "id = ?", id).Error
	if err != nil {
		log.Printf("Error retrieving saved lot: %v", err)
		return nil, err
	}

	log.Printf("Lot retrieved from DB: %+v", savedLot)
	return &savedLot, nil
}

func (p *PostgresStorage) GetLot(ctx context.Context, getLot *models.GetLotRequest) (*models.GetLotResponse, error) {
	log.Printf("GetLot request for ID: %s", getLot.Lot_id)

	var lot models.Lot
	err := p.db.WithContext(ctx).First(&lot, "id = ?", getLot.Lot_id).Error
	if err != nil {
		log.Printf("Error getting lot: %v", err)
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("lot not found")
		}
		return nil, err
	}

	log.Printf("Retrieved lot from DB: %+v", lot)
	return &models.GetLotResponse{Lot: lot}, nil
}

func (p *PostgresStorage) PlaceBid(ctx context.Context, placeBid *models.PlaceBidRequest) (*models.PlaceBidResponse, error) {
	log.Printf("PlaceBid request: %+v", placeBid)

	var lot models.Lot
	err := p.db.WithContext(ctx).First(&lot, "id = ?", placeBid.Lot_id).Error
	if err != nil {
		log.Printf("Lot not found: %v", err)
		return &models.PlaceBidResponse{
			Success: false,
			Message: "Лот не найден",
		}, nil
	}

	log.Printf("Current lot state: %+v", lot)

	if lot.Status != "ACTIVE" {
		return &models.PlaceBidResponse{
			Success:     false,
			Message:     "Аукцион для этого лота завершен",
			Updated_lot: lot,
		}, nil
	}

	if time.Now().Unix() > lot.End_time_unix {
		lot.Status = "COMPLETED"
		p.db.Save(&lot)
		return &models.PlaceBidResponse{
			Success:     false,
			Message:     "Аукцион завершен",
			Updated_lot: lot,
		}, nil
	}

	if placeBid.Amount <= lot.CurrentPrice {
		return &models.PlaceBidResponse{
			Success:     false,
			Message:     "Ставка должна быть выше текущей цены",
			Updated_lot: lot,
		}, nil
	}

	bid := &models.Bid{
		ID:             uuid.New().String(),
		LotId:          placeBid.Lot_id,
		UserId:         placeBid.User_id,
		Amount:         placeBid.Amount,
		Timestamp_unix: time.Now().Unix(),
		CreatedAt:      time.Now(),
	}

	if err := p.db.WithContext(ctx).Create(bid).Error; err != nil {
		log.Printf("Error creating bid: %v", err)
		return &models.PlaceBidResponse{
			Success:     false,
			Message:     "Ошибка при сохранении ставки",
			Updated_lot: lot,
		}, err
	}

	lot.CurrentPrice = placeBid.Amount
	lot.CurrentWinner = placeBid.User_id
	lot.UpdatedAt = time.Now()

	if err := p.db.WithContext(ctx).Save(&lot).Error; err != nil {
		log.Printf("Error updating lot: %v", err)
		return &models.PlaceBidResponse{
			Success:     false,
			Message:     "Ошибка при обновлении лота",
			Updated_lot: lot,
		}, err
	}

	var updatedLot models.Lot
	if err := p.db.WithContext(ctx).First(&updatedLot, "id = ?", placeBid.Lot_id).Error; err != nil {
		log.Printf("Error retrieving updated lot: %v", err)
		return &models.PlaceBidResponse{
			Success:     false,
			Message:     "Ошибка при получении обновленного лота",
			Updated_lot: lot,
		}, err
	}

	log.Printf("Updated lot after bid: %+v", updatedLot)

	return &models.PlaceBidResponse{
		Success:     true,
		Message:     "Ставка принята",
		Updated_lot: updatedLot,
	}, nil
}
