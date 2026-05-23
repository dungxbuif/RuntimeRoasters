package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/dungxbuif/RuntimeRoasters/apps/trace-service/internal/domain"
	"github.com/dungxbuif/RuntimeRoasters/apps/trace-service/internal/search"
	"github.com/dungxbuif/RuntimeRoasters/pkg/base/identity"
	"github.com/dungxbuif/RuntimeRoasters/pkg/kafka"
	"github.com/dungxbuif/RuntimeRoasters/pkg/logger"
	"github.com/google/uuid"
	kafkago "github.com/segmentio/kafka-go"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Service struct {
	db     *gorm.DB
	search *search.ElasticsearchClient
}

type TraceReadModel struct {
	EntityID   string                 `json:"entity_id"`
	EntityType string                 `json:"entity_type"`
	StoreID    string                 `json:"store_id,omitempty"`
	Status     string                 `json:"status"`
	Order      map[string]interface{} `json:"order,omitempty"`
	Payment    map[string]interface{} `json:"payment,omitempty"`
	Shipment   map[string]interface{} `json:"shipment,omitempty"`
	TraceIDs   []string               `json:"trace_ids,omitempty"`
	Timeline   []TimelineEntry        `json:"timeline"`
	Events     []domain.TraceEvent    `json:"events"`
	UpdatedAt  time.Time              `json:"updated_at"`
}

type TimelineEntry struct {
	Topic       string                 `json:"topic"`
	Status      string                 `json:"status"`
	Title       string                 `json:"title"`
	Description string                 `json:"description"`
	Timestamp   time.Time              `json:"timestamp"`
	Payload     map[string]interface{} `json:"payload,omitempty"`
}

func NewService(db *gorm.DB, searchClient *search.ElasticsearchClient) *Service {
	return &Service{db: db, search: searchClient}
}

func (s *Service) HandleEvent(ctx context.Context, msg kafkago.Message) error {
	messageID := kafka.MessageID(msg)
	ids := extractIDs(msg.Value)
	traceID := ""
	if spanContext := trace.SpanContextFromContext(ctx); spanContext.HasTraceID() {
		traceID = spanContext.TraceID().String()
	}
	event := domain.TraceEvent{
		ID:         uuid.NewString(),
		MessageID:  messageID,
		Topic:      msg.Topic,
		BatchID:    ids["batch_id"],
		OrderID:    ids["order_id"],
		ShipmentID: ids["shipment_id"],
		StoreID:    ids["store_id"],
		TraceID:    traceID,
		Payload:    string(msg.Value),
		OccurredAt: time.Now(),
	}

	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Clauses(clause.OnConflict{DoNothing: true}).Create(&event).Error; err != nil {
			return err
		}
		for _, key := range []string{"batch_id", "order_id", "shipment_id"} {
			if ids[key] == "" {
				continue
			}
			if err := s.rebuildDocument(ctx, tx, key, ids[key]); err != nil {
				return err
			}
		}
		return nil
	})
}

func (s *Service) GetTrace(ctx context.Context, entityID string) ([]domain.TraceEvent, error) {
	var events []domain.TraceEvent
	err := s.db.WithContext(ctx).
		Where("batch_id = ? OR order_id = ? OR shipment_id = ?", entityID, entityID, entityID).
		Order("occurred_at ASC, created_at ASC").
		Find(&events).Error
	if err != nil {
		return nil, err
	}
	if claims, ok := identity.FromContext(ctx); ok {
		filtered := events[:0]
		for _, event := range events {
			if claims.CanAccessStore(event.StoreID) {
				filtered = append(filtered, event)
			}
		}
		events = filtered
	}
	return events, err
}

func (s *Service) GetTraceDocument(ctx context.Context, entityID string) (map[string]interface{}, bool, error) {
	if s.search != nil {
		doc, found, err := s.search.GetTraceDocument(ctx, entityID)
		if err == nil && found {
			if claims, ok := identity.FromContext(ctx); ok && !claims.CanAccessStore(stringValue(doc, "store_id")) {
				return nil, false, nil
			}
			return doc, true, nil
		}
		if err != nil {
			logger.GetLogger().Warn("elasticsearch trace document lookup failed", zap.String("entity_id", entityID), zap.Error(err))
		}
	}

	var doc domain.TraceDocument
	err := s.db.WithContext(ctx).Where("entity_id = ?", entityID).First(&doc).Error
	if err == nil {
		if claims, ok := identity.FromContext(ctx); ok && !claims.CanAccessStore(doc.StoreID) {
			return nil, false, nil
		}
		var parsed map[string]interface{}
		if unmarshalErr := json.Unmarshal([]byte(doc.Document), &parsed); unmarshalErr != nil {
			return nil, false, unmarshalErr
		}
		return parsed, true, nil
	}
	if err == gorm.ErrRecordNotFound {
		return nil, false, nil
	}
	return nil, false, err
}

func (s *Service) rebuildDocument(ctx context.Context, tx *gorm.DB, entityType string, entityID string) error {
	var events []domain.TraceEvent
	if err := tx.WithContext(ctx).Where(entityType+" = ?", entityID).Order("occurred_at ASC").Find(&events).Error; err != nil {
		return err
	}
	readModel := buildReadModel(entityID, entityType, events)
	payload, err := json.Marshal(readModel)
	if err != nil {
		return err
	}
	traceIDsJSON, err := json.Marshal(readModel.TraceIDs)
	if err != nil {
		return err
	}
	doc := domain.TraceDocument{ID: uuid.NewString(), EntityID: entityID, EntityType: entityType, StoreID: readModel.StoreID, TraceIDs: string(traceIDsJSON), Document: string(payload), UpdatedAt: time.Now()}
	if err := tx.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "entity_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"entity_type", "store_id", "trace_ids", "document", "updated_at"}),
	}).Create(&doc).Error; err != nil {
		return err
	}
	if s.search != nil {
		if err := s.search.UpsertTraceDocument(ctx, search.TraceDocument{
			EntityID:   entityID,
			EntityType: entityType,
			Status:     readModel.Status,
			Order:      readModel.Order,
			Payment:    readModel.Payment,
			Shipment:   readModel.Shipment,
			TraceIDs:   readModel.TraceIDs,
			Timeline:   readModel.Timeline,
			Events:     events,
			UpdatedAt:  readModel.UpdatedAt,
		}); err != nil {
			logger.GetLogger().Warn("elasticsearch trace document upsert failed", zap.String("entity_id", entityID), zap.Error(err))
		}
	}
	return nil
}

func buildReadModel(entityID string, entityType string, events []domain.TraceEvent) TraceReadModel {
	model := TraceReadModel{
		EntityID:   entityID,
		EntityType: entityType,
		Status:     "PENDING",
		Timeline:   make([]TimelineEntry, 0, len(events)),
		Events:     events,
		UpdatedAt:  time.Now(),
	}
	seenTraceIDs := map[string]struct{}{}
	for _, event := range events {
		if event.TraceID != "" {
			if _, ok := seenTraceIDs[event.TraceID]; !ok {
				seenTraceIDs[event.TraceID] = struct{}{}
				model.TraceIDs = append(model.TraceIDs, event.TraceID)
			}
		}
		payload := parsePayload(event.Payload)
		entry := timelineEntry(event, payload)
		model.Timeline = append(model.Timeline, entry)
		model.Status = entry.Status
		switch event.Topic {
		case "retail.order.created":
			model.Order = payload
		case "payment.intent.created", "payment.completed", "payment.failed", "payment.refunded":
			model.Payment = merge(model.Payment, payload)
		case "logistics.shipment.assigned", "logistics.shipment.delivered":
			model.Shipment = merge(model.Shipment, payload)
		}
	}
	model.StoreID = firstStoreID(model)
	return model
}

func firstStoreID(m TraceReadModel) string {
	for _, event := range m.Events {
		if event.StoreID != "" {
			return event.StoreID
		}
	}
	for _, payload := range []map[string]interface{}{m.Order, m.Payment, m.Shipment} {
		if storeID := stringValue(payload, "store_id"); storeID != "" {
			return storeID
		}
	}
	return ""
}

func timelineEntry(event domain.TraceEvent, payload map[string]interface{}) TimelineEntry {
	status, title := eventStatus(event.Topic)
	description := eventDescription(event.Topic, payload)
	return TimelineEntry{
		Topic:       event.Topic,
		Status:      status,
		Title:       title,
		Description: description,
		Timestamp:   event.OccurredAt,
		Payload:     payload,
	}
}

func eventStatus(topic string) (string, string) {
	switch topic {
	case "retail.order.created":
		return "ORDER_CREATED", "Order created"
	case "payment.intent.created":
		return "PAYMENT_PENDING", "Payment intent created"
	case "payment.completed":
		return "PAYMENT_COMPLETED", "Payment completed"
	case "payment.failed":
		return "REJECTED", "Payment failed"
	case "payment.refunded":
		return "REFUNDED", "Payment refunded"
	case "warehouse.stock.reserved":
		return "STOCK_RESERVED", "Stock reserved"
	case "warehouse.stock.reservation_failed":
		return "REJECTED", "Stock reservation failed"
	case "logistics.shipment.assigned":
		return "SHIPPING", "Shipment assigned"
	case "logistics.shipment.delivered":
		return "COMPLETED", "Shipment delivered"
	case "logistics.gps.updated":
		return "IN_TRANSIT", "Driver location updated"
	default:
		return "UPDATED", topic
	}
}

func eventDescription(topic string, payload map[string]interface{}) string {
	switch topic {
	case "payment.intent.created":
		return fmt.Sprintf("%s simulated checkout created for %.2f %s", stringValue(payload, "provider"), floatValue(payload, "amount"), stringValue(payload, "currency"))
	case "payment.completed":
		return fmt.Sprintf("%s payment completed for %.2f %s", stringValue(payload, "provider"), floatValue(payload, "amount"), stringValue(payload, "currency"))
	case "payment.refunded":
		return fmt.Sprintf("%s refund issued because %s", stringValue(payload, "provider"), stringValue(payload, "reason"))
	case "warehouse.stock.reserved":
		return fmt.Sprintf("Reserved %.0f units of %s", floatValue(payload, "quantity"), stringValue(payload, "sku"))
	case "warehouse.stock.reservation_failed":
		return fmt.Sprintf("Stock reservation failed: %s", stringValue(payload, "reason"))
	case "logistics.shipment.assigned":
		return fmt.Sprintf("Shipment %s assigned to driver %s", stringValue(payload, "shipment_id"), stringValue(payload, "driver_id"))
	case "logistics.shipment.delivered":
		return fmt.Sprintf("Shipment %s delivered", stringValue(payload, "shipment_id"))
	default:
		return topic
	}
}

func parsePayload(payload string) map[string]interface{} {
	var raw map[string]interface{}
	_ = json.Unmarshal([]byte(payload), &raw)
	return raw
}

func merge(base map[string]interface{}, next map[string]interface{}) map[string]interface{} {
	if base == nil {
		base = map[string]interface{}{}
	}
	for key, value := range next {
		base[key] = value
	}
	return base
}

func stringValue(payload map[string]interface{}, key string) string {
	value, _ := payload[key].(string)
	return value
}

func floatValue(payload map[string]interface{}, key string) float64 {
	switch value := payload[key].(type) {
	case float64:
		return value
	case int:
		return float64(value)
	default:
		return 0
	}
}

func extractIDs(payload []byte) map[string]string {
	var raw map[string]interface{}
	_ = json.Unmarshal(payload, &raw)
	ids := map[string]string{}
	for _, key := range []string{"batch_id", "order_id", "shipment_id", "harvest_id", "store_id"} {
		if value, ok := raw[key].(string); ok {
			if key == "harvest_id" {
				ids["batch_id"] = value
				continue
			}
			ids[key] = value
		}
	}
	return ids
}
