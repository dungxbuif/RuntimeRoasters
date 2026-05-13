package domain

import (
	"strings"
	"time"

	"github.com/dungxbuif/RuntimeRoasters/pkg/errs"
)

// Farm represents a coffee farm entity.
type Farm struct {
	ID         uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	Name       string    `gorm:"size:255;not null" json:"name"`
	Location   string    `gorm:"type:location_enum;not null" json:"location"`
	Area       float64   `gorm:"type:decimal(10,2);not null" json:"area"`
	CoffeeType string    `gorm:"type:coffee_type_enum;not null" json:"coffee_type"`
	OwnerID    string    `gorm:"type:uuid;not null" json:"owner_id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// Simplified metadata for UI
var AllowedLocations = map[string]string{
	"CAU_DAT":       "Cầu Đất, Đà Lạt",
	"BUON_MA_THUOT": "Buôn Ma Thuột, Đắk Lắk",
	"PLEIKU":        "Pleiku, Gia Lai",
	"GIA_NGHIA":     "Gia Nghĩa, Đắk Nông",
	"KON_TUM":       "Kon Tum",
}

var AllowedCoffeeTypes = map[string]string{
	"ARABICA": "Arabica",
	"ROBUSTA": "Robusta",
	"CHERRY":  "Cherry",
	"CULI":    "Culi",
}

func (f *Farm) Validate() error {
	if strings.TrimSpace(f.Name) == "" {
		return errs.ErrValidation
	}
	if _, ok := AllowedLocations[f.Location]; !ok {
		return errs.ErrValidation
	}
	if _, ok := AllowedCoffeeTypes[f.CoffeeType]; !ok {
		return errs.ErrValidation
	}
	if f.Area < 0.1 {
		return errs.ErrValidation
	}
	return nil
}
