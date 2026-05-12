package domain

import (
	"time"

	"github.com/dungxbuif/RuntimeRoasters/pkg/errs"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type Farm struct {
	ID         uint64    `gorm:"primaryKey;autoIncrement"`
	Name       string    `gorm:"size:255;not null"`
	Location   string    `gorm:"type:text;not null"`
	Area       float64   `gorm:"type:decimal(10,2);not null"`
	CoffeeType string    `gorm:"size:100"`
	OwnerID    string    `gorm:"type:uuid;not null"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type TypeInfo struct {
	Code string
	Name string
}

var AllowedLocations = []TypeInfo{
	{Code: "CAU_DAT", Name: "Cầu Đất, Đà Lạt"},
	{Code: "BUON_MA_THUOT", Name: "Buôn Ma Thuột, Đắk Lắk"},
	{Code: "PLEIKU", Name: "Pleiku, Gia Lai"},
	{Code: "GIA_NGHIA", Name: "Gia Nghĩa, Đắk Nông"},
	{Code: "KON_TUM", Name: "Kon Tum"},
}

var AllowedCoffeeTypes = []TypeInfo{
	{Code: "ARABICA", Name: "Arabica"},
	{Code: "ROBUSTA", Name: "Robusta"},
	{Code: "CHERRY", Name: "Cherry"},
	{Code: "CULI", Name: "Culi"},
}

// Extract codes for validation
func getAllowedLocationCodes() []interface{} {
	codes := make([]interface{}, len(AllowedLocations))
	for i, v := range AllowedLocations {
		codes[i] = v.Code
	}
	return codes
}

func getAllowedCoffeeCodes() []interface{} {
	codes := make([]interface{}, len(AllowedCoffeeTypes))
	for i, v := range AllowedCoffeeTypes {
		codes[i] = v.Code
	}
	return codes
}

func (f *Farm) Validate() error {
	err := validation.ValidateStruct(f,
		validation.Field(&f.Name, validation.Required, validation.Length(1, 255)),
		validation.Field(&f.OwnerID, validation.Required),
		validation.Field(&f.Area, validation.Required, validation.Min(0.01)),
		validation.Field(&f.Location, validation.Required, validation.In(getAllowedLocationCodes()...)),
		validation.Field(&f.CoffeeType, validation.Required, validation.In(getAllowedCoffeeCodes()...)),
	)

	if err != nil {
		if valErrs, ok := err.(validation.Errors); ok {
			fieldErrs := make(map[string]error)
			for field, e := range valErrs {
				fieldErrs[field] = e
			}
			return errs.MapValidationErrors(fieldErrs)
		}
		return errs.ErrValidation
	}

	return nil
}
