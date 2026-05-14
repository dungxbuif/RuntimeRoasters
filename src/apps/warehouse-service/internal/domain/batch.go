package domain

import (
	"time"
)

type BatchStatus string

const (
	BatchStatusReceived   BatchStatus = "RECEIVED"
	BatchStatusProcessing BatchStatus = "PROCESSING"
	BatchStatusStocked    BatchStatus = "STOCKED"
)

type ProductionBatch struct {
	ID                string      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	HarvestID         string      `gorm:"type:uuid;not null;index"`
	BatchID           string      `gorm:"uniqueIndex;not null"`
	Status            BatchStatus `gorm:"size:20;not null;default:'RECEIVED'"`
	CoffeeType        string      `gorm:"size:20;not null"`
	OriginCode        string      `gorm:"size:10;not null"`
	IntakeWeight      float64     `gorm:"type:decimal(10,2);not null"`
	TotalInputWeight  float64     `gorm:"type:decimal(10,2);default:0"` // Sum of RoastRun inputs
	TotalOutputWeight float64     `gorm:"type:decimal(10,2);default:0"` // Sum of RoastRun outputs
	WeightLossPercent float64     `gorm:"type:decimal(5,2);default:0"`
	QualityFlag       string      `gorm:"size:20;default:'NORMAL'"`
	IntakeNote        string      `gorm:"type:text"`
	RoastRuns         []RoastRun  `gorm:"foreignKey:BatchID;references:BatchID"`
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

type RoastRun struct {
	ID           string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	BatchID      string    `gorm:"index;not null"`
	RunNumber    int       `gorm:"not null"`
	InputWeight  float64   `gorm:"type:decimal(10,2);not null"`
	OutputWeight float64   `gorm:"type:decimal(10,2);not null"`
	LossPercent  float64   `gorm:"type:decimal(5,2);not null"`
	RoasterName  string    `gorm:"size:50"`
	CreatedAt    time.Time
}

type InboxEvent struct {
	ID          string `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	MessageID   string `gorm:"uniqueIndex;not null"`
	ProcessedAt time.Time
}
