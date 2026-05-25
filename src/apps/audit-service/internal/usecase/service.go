package usecase

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"time"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	cassandrastore "RuntimeRoasters/apps/audit-service/internal/cassandra"
	"RuntimeRoasters/apps/audit-service/internal/domain"
	"RuntimeRoasters/pkg/base/identity"
	"RuntimeRoasters/pkg/events"
	"RuntimeRoasters/pkg/kafka"
	"RuntimeRoasters/pkg/logger"
	"github.com/google/uuid"
	kafkago "github.com/segmentio/kafka-go"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Service struct {
	db        *gorm.DB
	cassandra *cassandrastore.Store
}

func NewService(db *gorm.DB, cassandra *cassandrastore.Store) *Service {
	return &Service{db: db, cassandra: cassandra}
}

func (s *Service) HandleEvent(ctx context.Context, msg kafkago.Message) error {
	messageID := kafka.MessageID(msg)
	cloudEvent, err := events.ParseCloudEvent(msg.Value)
	if err != nil {
		return err
	}
	partitionKey := partitionKey(cloudEvent)
	var existing domain.AuditLog
	if err := s.db.WithContext(ctx).Where("message_id = ?", messageID).Take(&existing).Error; err == nil {
		return nil
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	var prev domain.AuditLog
	previousHash := ""
	if err := s.db.WithContext(ctx).Where("partition_key = ?", partitionKey).Order("occurred_at DESC, created_at DESC").Take(&prev).Error; err == nil {
		previousHash = prev.CurrentHash
	}
	currentHash := hash(previousHash, msg.Topic, string(msg.Value))
	log := domain.AuditLog{
		ID:           uuid.NewString(),
		PartitionKey: partitionKey,
		MessageID:    messageID,
		Topic:        cloudEvent.Type(),
		StoreID:      events.ExtensionString(cloudEvent, "storeid"),
		Payload:      string(msg.Value),
		PreviousHash: previousHash,
		CurrentHash:  currentHash,
		OccurredAt:   time.Now(),
		CreatedAt:    time.Now(),
	}
	if s.cassandra != nil {
		if err := s.cassandra.Insert(ctx, log); err != nil {
			logger.GetLogger().Warn("failed to write audit log to cassandra", zap.String("message_id", messageID), zap.Error(err))
		}
	}
	return s.db.WithContext(ctx).Clauses(clause.OnConflict{DoNothing: true}).Create(&log).Error
}

func (s *Service) ListByPartition(ctx context.Context, partition string) ([]domain.AuditLog, error) {
	if s.cassandra != nil {
		logs, err := s.cassandra.ListByPartition(ctx, partition)
		if err == nil && len(logs) > 0 {
			return filterLogs(ctx, logs), nil
		}
		if err != nil {
			logger.GetLogger().Warn("failed to read audit logs from cassandra", zap.String("partition", partition), zap.Error(err))
		}
	}
	var logs []domain.AuditLog
	err := s.db.WithContext(ctx).Where("partition_key = ?", partition).Order("occurred_at ASC").Find(&logs).Error
	return filterLogs(ctx, logs), err
}

func partitionKey(event cloudevents.Event) string {
	for _, key := range []string{"orderid", "shipmentid", "harvestid", "batchid", "storeid", "farmid", "warehouseid", "driverid", "vehicleid"} {
		if value := events.ExtensionString(event, key); value != "" {
			return value
		}
	}
	data := events.DataMap(event)
	for _, key := range []string{"order_id", "shipment_id", "harvest_id", "batch_id", "store_id", "farm_id", "warehouse_id", "driver_id", "vehicle_id"} {
		if value, ok := data[key].(string); ok && value != "" {
			return value
		}
	}
	return event.ID()
}

func filterLogs(ctx context.Context, logs []domain.AuditLog) []domain.AuditLog {
	claims, ok := identity.FromContext(ctx)
	if !ok {
		return logs
	}
	filtered := logs[:0]
	for _, log := range logs {
		if claims.CanAccessStore(log.StoreID) {
			filtered = append(filtered, log)
		}
	}
	return filtered
}

func hash(previous string, topic string, payload string) string {
	sum := sha256.Sum256([]byte(previous + "|" + topic + "|" + payload))
	return hex.EncodeToString(sum[:])
}
