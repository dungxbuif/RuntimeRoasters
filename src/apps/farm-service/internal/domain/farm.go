package domain

import (
	"time"

	"github.com/dungxbuif/RuntimeRoasters/pkg/errs"
)

type Farm struct {
	ID         string
	Name       string
	Location   string
	Area       float64
	CoffeeType string
	OwnerID    string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

var AllowedLocations = []string{
	"Cầu Đất, Đà Lạt",
	"Buôn Ma Thuột, Đắk Lắk",
	"Pleiku, Gia Lai",
	"Gia Nghĩa, Đắk Nông",
	"Kon Tum",
}

func (f *Farm) Validate() error {
	if f.Name == "" {
		return errs.ErrValidation
	}
	if f.OwnerID == "" {
		return errs.ErrValidation
	}

	isValidLocation := false
	for _, loc := range AllowedLocations {
		if f.Location == loc {
			isValidLocation = true
			break
		}
	}
	if !isValidLocation {
		return errs.ErrValidation
	}

	return nil
}
