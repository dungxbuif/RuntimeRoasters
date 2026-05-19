package domain

import (
	"time"
)

type ShipmentStatus string

const (
	ShipmentStatusPending   ShipmentStatus = "PENDING"
	ShipmentStatusAssigned  ShipmentStatus = "ASSIGNED"
	ShipmentStatusPickedUp  ShipmentStatus = "PICKED_UP"
	ShipmentStatusInTransit ShipmentStatus = "IN_TRANSIT"
	ShipmentStatusDelivered ShipmentStatus = "DELIVERED"
	ShipmentStatusFailed    ShipmentStatus = "FAILED"
	ShipmentStatusCancelled ShipmentStatus = "CANCELLED"
)

type Shipment struct {
	ID                 string         `gorm:"primaryKey" json:"id"`
	OrderID            string         `gorm:"index" json:"order_id"`
	BatchID            string         `gorm:"index" json:"batch_id"`
	OriginWarehouseID  string         `gorm:"size:80" json:"origin_warehouse_id"`
	DestinationStoreID string         `gorm:"index" json:"destination_store_id"`
	DriverID           string         `gorm:"index" json:"driver_id"`
	Status             ShipmentStatus `gorm:"not null;index" json:"status"`
	DestinationAddress string         `json:"destination_address"`
	CreatedAt          time.Time      `json:"created_at"`
	UpdatedAt          time.Time      `json:"updated_at"`
	DeliveredAt        *time.Time     `json:"delivered_at"`
}
