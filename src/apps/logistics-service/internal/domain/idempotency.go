package domain

import "time"

type ProcessedKafkaMessage struct {
	MsgKey      string    `gorm:"primaryKey"`
	Topic       string    `gorm:"not null"`
	ProcessedAt time.Time `gorm:"not null"`
}

type InboxEvent struct {
	ID          string    `gorm:"primaryKey"`
	MessageID   string    `gorm:"uniqueIndex;not null"`
	EventType   string    `gorm:"size:120;not null"`
	ProcessedAt time.Time `gorm:"not null"`
}
