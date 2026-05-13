package domain

import (
	"errors"
	"strings"
	"time"
)

var (
	ErrInvalidFarmName    = errors.New("farm name is required")
	ErrInvalidLocation    = errors.New("invalid or unsupported location")
	ErrInvalidFarmArea    = errors.New("area must be at least 0.1")
	ErrInvalidFarmCoffee  = errors.New("invalid or unsupported coffee type for farm")
)

type Location string
type CoffeeType string

// Farm represents a coffee farm entity.
type Farm struct {
	ID         uint64     `gorm:"primaryKey;autoIncrement" json:"id"`
	Name       string     `gorm:"size:255;not null" json:"name"`
	Location   Location   `gorm:"type:location_enum;not null" json:"location"`
	Area       float64    `gorm:"type:decimal(10,2);not null" json:"area"`
	CoffeeType CoffeeType `gorm:"type:coffee_type_enum;not null" json:"coffee_type"`
	OwnerID    string     `gorm:"type:uuid;not null" json:"owner_id"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
}

const (
	LocationCauDat       Location = "CAU_DAT"
	LocationBuonMaThuot  Location = "BUON_MA_THUOT"
	LocationPleiku       Location = "PLEIKU"
	LocationGiaNghia     Location = "GIA_NGHIA"
	LocationKonTum       Location = "KON_TUM"
)

const (
	CoffeeTypeArabica CoffeeType = "ARABICA"
	CoffeeTypeRobusta CoffeeType = "ROBUSTA"
	CoffeeTypeCherry  CoffeeType = "CHERRY"
	CoffeeTypeCuli    CoffeeType = "CULI"
)

// Simplified metadata for UI
var AllowedLocations = map[Location]string{
	LocationCauDat:       "Cầu Đất, Đà Lạt",
	LocationBuonMaThuot:  "Buôn Ma Thuột, Đắk Lắk",
	LocationPleiku:       "Pleiku, Gia Lai",
	LocationGiaNghia:     "Gia Nghĩa, Đắk Nông",
	LocationKonTum:       "Kon Tum",
}

var AllowedCoffeeTypes = map[CoffeeType]string{
	CoffeeTypeArabica: "Arabica",
	CoffeeTypeRobusta: "Robusta",
	CoffeeTypeCherry:  "Cherry",
	CoffeeTypeCuli:    "Culi",
}

func (f *Farm) Validate() error {
	if strings.TrimSpace(f.Name) == "" {
		return ErrInvalidFarmName
	}
	if _, ok := AllowedLocations[f.Location]; !ok {
		return ErrInvalidLocation
	}
	if _, ok := AllowedCoffeeTypes[f.CoffeeType]; !ok {
		return ErrInvalidFarmCoffee
	}
	if f.Area < 0.1 {
		return ErrInvalidFarmArea
	}
	return nil
}
