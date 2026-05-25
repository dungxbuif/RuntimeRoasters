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
	ID                string       `gorm:"primaryKey" json:"id"`
	UserID            string       `gorm:"size:160;index" json:"user_id"`
	Name              string       `gorm:"not null" json:"name"`
	Phone             string       `gorm:"not null" json:"phone"`
	VehicleID         string       `json:"vehicle_id"`
	Status            DriverStatus `gorm:"size:24;not null;index;default:'IDLE'" json:"status"`
	CurrentShipmentID string       `gorm:"size:80;index" json:"current_shipment_id"`
	IsAvailable       bool         `gorm:"default:true" json:"is_available"`
	CreatedAt         time.Time    `json:"created_at"`
	UpdatedAt         time.Time    `json:"updated_at"`
}

type Vehicle struct {
	ID               string       `gorm:"primaryKey" json:"id"`
	PlateNumber      string       `gorm:"not null;unique" json:"plate_number"`
	Type             string       `gorm:"not null" json:"type"` // e.g., TRUCK, VAN
	CapacityKG       float64      `json:"capacity_kg"`
	HomeWarehouseID  string       `gorm:"size:80;index" json:"home_warehouse_id"`
	Status           DriverStatus `gorm:"size:24;not null;index;default:'IDLE'" json:"status"`
	CurrentLatitude  float64      `json:"current_latitude"`
	CurrentLongitude float64      `json:"current_longitude"`
	CreatedAt        time.Time    `json:"created_at"`
	UpdatedAt        time.Time    `json:"updated_at"`
}
