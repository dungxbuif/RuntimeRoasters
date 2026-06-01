package usecase

import (
	"time"

	"RuntimeRoasters/apps/warehouse-service/internal/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func markInbox(tx *gorm.DB, messageID string, eventType string) error {
	return tx.Create(&domain.InboxEvent{
		ID:          uuid.NewString(),
		MessageID:   messageID,
		EventType:   eventType,
		ProcessedAt: time.Now(),
	}).Error
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if value != "" {
			return value
		}
	}
	return ""
}
