package domain

import (
	"time"
)

type BatchStatus string
type IntakeStatus string
type PickupRequestStatus string

const (
	BatchStatusDraft      BatchStatus = "DRAFT"
	BatchStatusProcessing BatchStatus = "PROCESSING"
	BatchStatusReady      BatchStatus = "READY_TO_STOCK"
	BatchStatusStocked    BatchStatus = "STOCKED"

	IntakeStatusUnassigned IntakeStatus = "UNASSIGNED"
	IntakeStatusAssigned   IntakeStatus = "ASSIGNED"

	PickupRequestStatusRequested        PickupRequestStatus = "REQUESTED"
	PickupRequestStatusDispatched       PickupRequestStatus = "DISPATCHED"
	PickupRequestStatusInTransitToFarm  PickupRequestStatus = "IN_TRANSIT_TO_FARM"
	PickupRequestStatusPickedUp         PickupRequestStatus = "PICKED_UP"
	PickupRequestStatusReturning        PickupRequestStatus = "RETURNING"
	PickupRequestStatusArrivedWarehouse PickupRequestStatus = "ARRIVED_WAREHOUSE"
	PickupRequestStatusReceived         PickupRequestStatus = "RECEIVED"
	PickupRequestStatusCancelled        PickupRequestStatus = "CANCELLED"
)

// Intake represents raw coffee received from a farm
type Intake struct {
	ID          string       `gorm:"type:varchar(64);primaryKey"`
	HarvestID   string       `gorm:"type:varchar(64);not null;index"`
	PickupID    string       `gorm:"type:varchar(64);index"`
	WarehouseID string       `gorm:"type:varchar(80);index" json:"warehouse_id"`
	CoffeeType  string       `gorm:"size:20;not null"`
	OriginCode  string       `gorm:"size:10;not null"`
	Quantity    float64      `gorm:"type:decimal(10,2);not null"`
	Status      IntakeStatus `gorm:"size:20;not null;default:'UNASSIGNED'"`
	BatchID     *string      `gorm:"type:varchar(80);index"` // Link to ProductionBatch
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type PickupRequest struct {
	ID               string              `gorm:"type:varchar(64);primaryKey" json:"id"`
	HarvestID        string              `gorm:"type:varchar(64);uniqueIndex;not null" json:"harvest_id"`
	FarmID           string              `gorm:"type:varchar(80);index" json:"farm_id"`
	WarehouseID      string              `gorm:"type:varchar(80);index;not null" json:"warehouse_id"`
	OriginLocationID string              `gorm:"type:varchar(80);index" json:"origin_location_id"`
	Quantity         float64             `gorm:"type:decimal(10,2);not null" json:"quantity"`
	CoffeeType       string              `gorm:"size:20;not null" json:"coffee_type"`
	OriginCode       string              `gorm:"size:10;not null" json:"origin_code"`
	Status           PickupRequestStatus `gorm:"size:32;not null;default:'REQUESTED';index" json:"status"`
	NotificationID   string              `gorm:"type:varchar(64)" json:"notification_id,omitempty"`
	ShipmentID       string              `gorm:"type:varchar(64);index" json:"shipment_id,omitempty"`
	CreatedAt        time.Time           `json:"created_at"`
	UpdatedAt        time.Time           `json:"updated_at"`
	DispatchedAt     *time.Time          `json:"dispatched_at,omitempty"`
	ReceivedAt       *time.Time          `json:"received_at,omitempty"`
}

// ProductionBatch represents a production run aggregating multiple intakes
type ProductionBatch struct {
	ID                string      `gorm:"type:varchar(80);primaryKey"`
	BatchID           string      `gorm:"uniqueIndex;not null"`
	WarehouseID       string      `gorm:"type:varchar(80);index" json:"warehouse_id"`
	Status            BatchStatus `gorm:"size:20;not null;default:'DRAFT'"`
	Intakes           []Intake    `gorm:"foreignKey:BatchID"`
	TotalInputWeight  float64     `gorm:"type:decimal(10,2);default:0"`
	TotalOutputWeight float64     `gorm:"type:decimal(10,2);default:0"`
	WeightLossPercent float64     `gorm:"type:decimal(5,2);default:0"`
	RoastRuns         []RoastRun  `gorm:"foreignKey:BatchID"`
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

// RoastRun represents a simulated roasting cycle within a batch
type RoastRun struct {
	ID           string  `gorm:"type:varchar(64);primaryKey"`
	BatchID      string  `gorm:"type:varchar(80);index;not null"`
	RunNumber    int     `gorm:"not null"`
	InputWeight  float64 `gorm:"type:decimal(10,2);not null"`
	OutputWeight float64 `gorm:"type:decimal(10,2);not null"`
	Status       string  `gorm:"size:20;not null;default:'COMPLETED'"`
	CreatedAt    time.Time
}

func (Intake) TableName() string {
	return "intakes"
}

func (PickupRequest) TableName() string {
	return "pick_up_requests"
}

func (ProductionBatch) TableName() string {
	return "production_batches"
}

func (RoastRun) TableName() string {
	return "roast_runs"
}

type InboxEvent struct {
	ID          string `gorm:"type:varchar(64);primaryKey"`
	MessageID   string `gorm:"uniqueIndex;not null"`
	EventType   string `gorm:"size:160;index"`
	ProcessedAt time.Time
}
