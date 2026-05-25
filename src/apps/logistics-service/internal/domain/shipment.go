package domain

import (
	"time"
)

type ShipmentStatus string
type ShipmentType string
type ShipmentLeg string

const (
	ShipmentTypeFarmPickup     ShipmentType = "FARM_PICKUP"
	ShipmentTypeRetailDelivery ShipmentType = "RETAIL_DELIVERY"
	ShipmentTypeReturnToBase   ShipmentType = "RETURN_TO_BASE"

	ShipmentLegOutbound ShipmentLeg = "OUTBOUND"
	ShipmentLegReturn   ShipmentLeg = "RETURN"

	ShipmentStatusPending              ShipmentStatus = "PENDING"
	ShipmentStatusAssigned             ShipmentStatus = "ASSIGNED"
	ShipmentStatusPickedUp             ShipmentStatus = "PICKED_UP"
	ShipmentStatusInTransit            ShipmentStatus = "IN_TRANSIT"
	ShipmentStatusDelivered            ShipmentStatus = "DELIVERED"
	ShipmentStatusFailed               ShipmentStatus = "FAILED"
	ShipmentStatusCancelled            ShipmentStatus = "CANCELLED"
	ShipmentStatusInTransitToFarm      ShipmentStatus = "IN_TRANSIT_TO_FARM"
	ShipmentStatusArrivedAtFarm        ShipmentStatus = "ARRIVED_AT_FARM"
	ShipmentStatusReturningToWarehouse ShipmentStatus = "RETURNING_TO_WAREHOUSE"
	ShipmentStatusArrivedWarehouse     ShipmentStatus = "ARRIVED_WAREHOUSE"
	ShipmentStatusInTransitToStore     ShipmentStatus = "IN_TRANSIT_TO_STORE"
	ShipmentStatusArrivedAtStore       ShipmentStatus = "ARRIVED_AT_STORE"
	ShipmentStatusReturningToBase      ShipmentStatus = "RETURNING_TO_BASE"
	ShipmentStatusReturnedToBase       ShipmentStatus = "RETURNED_TO_BASE"
	ShipmentStatusCompleted            ShipmentStatus = "COMPLETED"
)

type Shipment struct {
	ID                    string         `gorm:"primaryKey" json:"id"`
	Type                  ShipmentType   `gorm:"size:32;not null;index;default:'RETAIL_DELIVERY'" json:"type"`
	Status                ShipmentStatus `gorm:"not null;index" json:"status"`
	CurrentLeg            ShipmentLeg    `gorm:"size:24;index" json:"current_leg"`
	RouteID               string         `gorm:"size:80;index" json:"route_id"`
	OriginLocationID      string         `gorm:"size:80;index" json:"origin_location_id"`
	DestinationLocationID string         `gorm:"size:80;index" json:"destination_location_id"`
	FarmID                string         `gorm:"size:80;index" json:"farm_id"`
	HarvestID             string         `gorm:"size:80;index" json:"harvest_id"`
	WarehouseID           string         `gorm:"size:80;index" json:"warehouse_id"`
	OrderID               string         `gorm:"index" json:"order_id"`
	BatchID               string         `gorm:"index" json:"batch_id"`
	OriginWarehouseID     string         `gorm:"size:80" json:"origin_warehouse_id"`
	DestinationStoreID    string         `gorm:"index" json:"destination_store_id"`
	StoreID               string         `gorm:"size:80;index" json:"store_id"`
	DriverID              string         `gorm:"index" json:"driver_id"`
	VehicleID             string         `gorm:"size:80;index" json:"vehicle_id"`
	DestinationAddress    string         `json:"destination_address"`
	CreatedAt             time.Time      `json:"created_at"`
	UpdatedAt             time.Time      `json:"updated_at"`
	AssignedAt            *time.Time     `json:"assigned_at"`
	DepartedAt            *time.Time     `json:"departed_at"`
	ArrivedAtOriginAt     *time.Time     `json:"arrived_at_origin_at"`
	LoadedAt              *time.Time     `json:"loaded_at"`
	DepartedOriginAt      *time.Time     `json:"departed_origin_at"`
	ArrivedDestinationAt  *time.Time     `json:"arrived_destination_at"`
	DeliveredAt           *time.Time     `json:"delivered_at"`
	ReturnStartedAt       *time.Time     `json:"return_started_at"`
	ReturnedAt            *time.Time     `json:"returned_at"`
}
