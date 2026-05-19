package domain

import (
	"time"
)

type Inventory struct {
	ID                string  `gorm:"type:varchar(64);primaryKey"`
	CoffeeType        string  `gorm:"size:20;not null"`
	OriginCode        string  `gorm:"size:10;not null"`
	SKU               string  `gorm:"uniqueIndex;not null"`
	AvailableQuantity float64 `gorm:"type:decimal(12,2);default:0"`
	UpdatedAt         time.Time
}
