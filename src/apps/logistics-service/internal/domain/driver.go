package domain

import (
	"time"
)

type DriverStatus string

const (
	DriverStatusIdle    DriverStatus = "IDLE"
	DriverStatusBusy    DriverStatus = "BUSY"
	DriverStatusOffline DriverStatus = "OFFLINE"
)

type Driver struct {
	ID          string       `gorm:"primaryKey" json:"id"`
	Name        string       `gorm:"not null" json:"name"`
	Phone       string       `gorm:"not null" json:"phone"`
	VehicleID   string       `json:"vehicle_id"`
	Status      DriverStatus `gorm:"size:24;not null;index;default:'IDLE'" json:"status"`
	IsAvailable bool         `gorm:"default:true" json:"is_available"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
}

type Vehicle struct {
	ID          string    `gorm:"primaryKey" json:"id"`
	PlateNumber string    `gorm:"not null;unique" json:"plate_number"`
	Type        string    `gorm:"not null" json:"type"` // e.g., TRUCK, VAN
	CapacityKG  float64   `json:"capacity_kg"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
