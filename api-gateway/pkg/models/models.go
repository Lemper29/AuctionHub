package models

import (
	"time"
)

type Lot struct {
	Id            string    `gorm:"primaryKey;column:id" json:"id"`
	Name          string    `gorm:"column:name" json:"name"`
	Description   string    `gorm:"column:description" json:"description"`
	StartPrice    float64   `gorm:"column:start_price" json:"startPrice"`
	CurrentPrice  float64   `gorm:"column:current_price" json:"currentPrice"`
	CurrentWinner string    `gorm:"column:current_winner" json:"currentWinner"`
	Status        string    `gorm:"column:status" json:"status"`
	End_time_unix int64     `gorm:"column:end_time_unix" json:"end_time_unix"`
	CreatedAt     time.Time `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt     time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
}

func (Lot) TableName() string {
	return "lots"
}

type Bid struct {
	ID             string    `gorm:"primaryKey;column:id" json:"id"`
	LotId          string    `gorm:"column:lot_id" json:"lotId"`
	UserId         string    `gorm:"column:user_id" json:"userId"`
	Amount         float64   `gorm:"column:amount" json:"amount"`
	Timestamp_unix int64     `gorm:"column:timestamp_unix" json:"timestamp_unix"`
	CreatedAt      time.Time `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
}

func (Bid) TableName() string {
	return "bids"
}

type CreateLotRequest struct {
	Name           string
	Description    string
	StartPrice     float64
	DurationMinute int64
}

type CreateLotResponse struct {
	Lot Lot
}

type GetLotRequest struct {
	Lot_id string
}

type GetLotResponse struct {
	Lot Lot
}

type PlaceBidRequest struct {
	Lot_id  string
	User_id string
	Amount  float64
}

type PlaceBidResponse struct {
	Success     bool
	Message     string
	Updated_lot Lot
}
