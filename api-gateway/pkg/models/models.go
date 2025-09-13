package models

import (
	"time"
)

type Lot struct {
	Id             string `gorm:"primaryKey"`
	Name           string
	Description    string
	Start_price    float64
	Current_price  float64
	Current_winner string
	Status         string
	End_time_unix  int64
	CreatedAt      time.Time `gorm:"autoCreateTime"`
	UpdatedAt      time.Time `gorm:"autoUpdateTime"`
}

type Bid struct {
	ID             string `gorm:"primaryKey"`
	LotId          string
	UserId         string
	Amount         float64
	Timestamp_unix int64
}

type CreateLotRequest struct {
	Name            string
	Description     string
	Start_price     float64
	Duration_minute int64
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
