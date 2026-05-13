package domain

import (
	"errors"
	"time"
)

var (
	ErrInvalidFarmID     = errors.New("farm id is required")
	ErrInvalidQuantity   = errors.New("quantity must be greater than 0")
	ErrInvalidCoffeeType = errors.New("invalid or unsupported coffee type")
)

type HarvestStatus string

const (
	StatusNew        HarvestStatus = "NEW"
	StatusProcessing HarvestStatus = "PROCESSING"
	StatusCompleted  HarvestStatus = "COMPLETED"
)

type Harvest struct {
	ID          uint64        `gorm:"primaryKey;autoIncrement" json:"id"`
	FarmID      uint64        `gorm:"not null" json:"farm_id"`
	OwnerID     string        `gorm:"type:uuid;not null" json:"owner_id"`
	CoffeeType  CoffeeType    `gorm:"type:varchar(50);not null" json:"coffee_type"`
	Quantity    float64       `gorm:"type:decimal(10,2);not null" json:"quantity"`
	HarvestDate time.Time     `gorm:"default:CURRENT_TIMESTAMP" json:"harvest_date"`
	Status      HarvestStatus `gorm:"type:harvest_status_enum;not null;default:'NEW'" json:"status"`
	Notes       string        `gorm:"type:text" json:"notes"`
	CreatedAt   time.Time     `json:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at"`
}

func (h *Harvest) Validate() error {
	if h.FarmID == 0 {
		return ErrInvalidFarmID
	}
	if h.Quantity <= 0 {
		return ErrInvalidQuantity
	}
	switch h.CoffeeType {
	case CoffeeTypeArabica, CoffeeTypeRobusta, CoffeeTypeCherry, CoffeeTypeCuli:
	default:
		return ErrInvalidCoffeeType
	}
	return nil
}
