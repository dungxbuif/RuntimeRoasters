package domain

import (
	"time"
)

type LocationType string

const (
	LocationTypeFarm     LocationType = "FARM"
	LocationTypeRoastery LocationType = "ROASTERY"
	LocationTypeRetailer LocationType = "RETAILER"
)

type Location struct {
	ID        string       `gorm:"primaryKey" json:"id"`
	Name      string       `gorm:"not null" json:"name"`
	Type      LocationType `gorm:"not null" json:"type"`
	Lat       float64      `gorm:"not null" json:"lat"`
	Lng       float64      `gorm:"not null" json:"lng"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
}
