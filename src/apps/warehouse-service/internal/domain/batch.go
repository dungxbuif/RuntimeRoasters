package domain

import (
	"time"
)

type BatchStatus string
type IntakeStatus string

const (
	BatchStatusDraft      BatchStatus = "DRAFT"
	BatchStatusProcessing BatchStatus = "PROCESSING"
	BatchStatusReady      BatchStatus = "READY_TO_STOCK"
	BatchStatusStocked    BatchStatus = "STOCKED"

	IntakeStatusUnassigned IntakeStatus = "UNASSIGNED"
	IntakeStatusAssigned   IntakeStatus = "ASSIGNED"
)

// Intake represents raw coffee received from a farm
type Intake struct {
	ID         string       `gorm:"type:uuid;primaryKey"`
	HarvestID  string       `gorm:"type:varchar(64);not null;index"`
	CoffeeType string       `gorm:"size:20;not null"`
	OriginCode string       `gorm:"size:10;not null"`
	Quantity   float64      `gorm:"type:decimal(10,2);not null"`
	Status     IntakeStatus `gorm:"size:20;not null;default:'UNASSIGNED'"`
	BatchID    *string      `gorm:"type:uuid;index"` // Link to ProductionBatch
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

// ProductionBatch represents a production run aggregating multiple intakes
type ProductionBatch struct {
	ID                string      `gorm:"type:uuid;primaryKey"`
	BatchID           string      `gorm:"uniqueIndex;not null"`
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
	ID           string  `gorm:"type:uuid;primaryKey"`
	BatchID      string  `gorm:"type:uuid;index;not null"`
	RunNumber    int     `gorm:"not null"`
	InputWeight  float64 `gorm:"type:decimal(10,2);not null"`
	OutputWeight float64 `gorm:"type:decimal(10,2);not null"`
	Status       string  `gorm:"size:20;not null;default:'COMPLETED'"`
	CreatedAt    time.Time
}

type InboxEvent struct {
	ID          string `gorm:"type:uuid;primaryKey"`
	MessageID   string `gorm:"uniqueIndex;not null"`
	ProcessedAt time.Time
}
