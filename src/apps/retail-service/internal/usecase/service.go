package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"RuntimeRoasters/apps/retail-service/internal/domain"
	retailseed "RuntimeRoasters/apps/retail-service/internal/seed"
	"RuntimeRoasters/pkg/base/identity"
	"RuntimeRoasters/pkg/events"
	"RuntimeRoasters/pkg/kafka"
	"RuntimeRoasters/pkg/logger"
	"github.com/google/uuid"
	kafkago "github.com/segmentio/kafka-go"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type CreateOrderRequest struct {
	StoreID       string                   `json:"store_id"`
	Items         []events.RetailOrderItem `json:"items"`
	TotalAmount   float64                  `json:"total_amount"`
	PaymentMethod string                   `json:"payment_method"`
}

type Service struct {
	db                *gorm.DB
	producer          kafka.Producer
	orderCreatedTopic string
}

func NewService(db *gorm.DB, producer kafka.Producer, orderCreatedTopic string) *Service {
	return &Service{db: db, producer: producer, orderCreatedTopic: orderCreatedTopic}
}

func (s *Service) SeedStores(ctx context.Context, stores []domain.Store) error {
	if len(stores) == 0 {
		return errors.New("retail store seed is empty")
	}
	for _, store := range stores {
		if err := s.db.WithContext(ctx).Save(&store).Error; err != nil {
			return err
		}
	}
	return nil
}

func (s *Service) GetSystemStatus(ctx context.Context) (*domain.SystemStatus, error) {
	var storeCount int64
	s.db.WithContext(ctx).Model(&domain.Store{}).Count(&storeCount)

	seeded := storeCount > 0

	return &domain.SystemStatus{
		Seeded:      seeded,
		ServiceName: "retail-service",
		RecordCounts: map[string]int64{
			"stores": storeCount,
		},
	}, nil
}

func (s *Service) SeedData(ctx context.Context, force bool, usersMap map[string]string) (*domain.SeedResult, error) {
	status, err := s.GetSystemStatus(ctx)
	if err != nil {
		return nil, err
	}

	if status.Seeded && !force {
		return &domain.SeedResult{
			Success: true,
			Message: "System already seeded",
		}, nil
	}

	stores, err := retailseed.LoadStores()
	if err != nil {
		return nil, err
	}

	// Enrich stores with ManagerID from usersMap
	for i := range stores {
		if id, ok := usersMap[stores[i].ManagerEmail]; ok {
			stores[i].ManagerID = id
		} else {
			logger.GetLogger().Warn("Manager email not found in users_map during seeding", 
				zap.String("email", stores[i].ManagerEmail), 
				zap.String("store", stores[i].Name))
		}
	}

	if err := s.SeedStores(ctx, stores); err != nil {
		return nil, err
	}

	return &domain.SeedResult{
		Success:        true,
		Message:        "Successfully seeded retail data",
		RecordsCreated: int64(len(stores)),
	}, nil
}

func (s *Service) CreateStore(ctx context.Context, store domain.Store) (domain.Store, error) {
	if store.ID == "" {
		store.ID = uuid.NewString()
	}
	now := time.Now()
	store.CreatedAt = now
	store.UpdatedAt = now
	
	err := s.db.WithContext(ctx).Where("name = ? AND city = ?", store.Name, store.City).FirstOrCreate(&store).Error
	return store, err
}

func (s *Service) ListStores(ctx context.Context) ([]domain.Store, error) {
	var stores []domain.Store
	_, storeIDs, allStores := identity.StoreScopeFromContext(ctx)
	query := s.db.WithContext(ctx).Order("city, name")
	if !allStores {
		if len(storeIDs) == 0 {
			return stores, nil
		}
		query = query.Where("id IN ?", storeIDs)
	}
	err := query.Find(&stores).Error
	return stores, err
}

func (s *Service) CreateOrder(ctx context.Context, req CreateOrderRequest, idempotencyKey string) (*domain.Order, error) {
	if req.StoreID == "" {
		return nil, errors.New("store_id is required")
	}
	if len(req.Items) == 0 {
		return nil, errors.New("at least one item is required")
	}
	if claims, ok := identity.FromContext(ctx); ok && !claims.CanAccessStore(req.StoreID) {
		return nil, errors.New("store access denied")
	}
	if idempotencyKey == "" {
		idempotencyKey = uuid.NewString()
	}

	var existing domain.Order
	if err := s.db.WithContext(ctx).Where("idempotency_key = ?", idempotencyKey).First(&existing).Error; err == nil {
		return &existing, nil
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	itemsJSON, err := json.Marshal(req.Items)
	if err != nil {
		return nil, err
	}

	orderID := uuid.NewString()
	event := events.RetailOrderCreated{
		EventID:       uuid.NewString(),
		OrderID:       orderID,
		StoreID:       req.StoreID,
		Items:         req.Items,
		TotalAmount:   req.TotalAmount,
		PaymentMethod: req.PaymentMethod,
		OccurredAt:    time.Now(),
	}
	cloudEvent, err := events.NewCloudEvent(ctx, s.orderCreatedTopic, events.SourceRetailService, fmt.Sprintf("orders/%s", orderID), event, events.Metadata{
		EventID:       event.EventID,
		CorrelationID: orderID,
		OccurredAt:    event.OccurredAt,
		OrderID:       orderID,
		StoreID:       req.StoreID,
	})
	if err != nil {
		return nil, err
	}
	payload, err := json.Marshal(cloudEvent)
	if err != nil {
		return nil, err
	}
	traceHeaders := propagation.MapCarrier{}
	otel.GetTextMapPropagator().Inject(ctx, traceHeaders)

	order := &domain.Order{
		ID:             orderID,
		StoreID:        req.StoreID,
		Items:          string(itemsJSON),
		TotalAmount:    req.TotalAmount,
		Status:         domain.OrderStatusPending,
		IdempotencyKey: idempotencyKey,
	}
	outbox := &domain.OutboxEvent{
		ID:          uuid.NewString(),
		EventType:   s.orderCreatedTopic,
		Topic:       s.orderCreatedTopic,
		Key:         orderID,
		Payload:     string(payload),
		TraceParent: traceHeaders.Get("traceparent"),
		TraceState:  traceHeaders.Get("tracestate"),
		Status:      domain.OutboxStatusPending,
	}

	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(order).Error; err != nil {
			return err
		}
		return tx.Create(outbox).Error
	})
	return order, err
}

func (s *Service) GetOrder(ctx context.Context, id string) (*domain.Order, error) {
	var order domain.Order
	if err := s.db.WithContext(ctx).Preload("Store").Where("id = ?", id).First(&order).Error; err != nil {
		return nil, err
	}
	if claims, ok := identity.FromContext(ctx); ok && !claims.CanAccessStore(order.StoreID) {
		return nil, gorm.ErrRecordNotFound
	}
	return &order, nil
}

func (s *Service) ProcessOutbox(ctx context.Context, limit int) error {
	var outbox []domain.OutboxEvent
	if err := s.db.WithContext(ctx).Where("status = ?", domain.OutboxStatusPending).Order("created_at ASC").Limit(limit).Find(&outbox).Error; err != nil {
		return err
	}
	for _, event := range outbox {
		eventCtx := contextFromOutboxTrace(ctx, event)
		if err := s.producer.Publish(eventCtx, event.Topic, event.Key, json.RawMessage(event.Payload)); err != nil {
			logger.FromContext(eventCtx).Warn("failed to publish retail outbox event", zap.Error(err), zap.String("event_id", event.ID))
			_ = s.db.WithContext(ctx).Model(&event).Update("status", domain.OutboxStatusFailed).Error
			continue
		}
		if err := s.db.WithContext(ctx).Model(&event).Update("status", domain.OutboxStatusCompleted).Error; err != nil {
			return err
		}
	}
	return nil
}

func contextFromOutboxTrace(ctx context.Context, event domain.OutboxEvent) context.Context {
	carrier := propagation.MapCarrier{}
	if event.TraceParent != "" {
		carrier.Set("traceparent", event.TraceParent)
	}
	if event.TraceState != "" {
		carrier.Set("tracestate", event.TraceState)
	}
	if len(carrier) == 0 {
		return ctx
	}
	return otel.GetTextMapPropagator().Extract(ctx, carrier)
}

func (s *Service) HandleSagaEvent(ctx context.Context, msg kafkago.Message) error {
	messageID := kafka.MessageID(msg)
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var existing domain.InboxEvent
		if err := tx.Where("message_id = ?", messageID).Take(&existing).Error; err == nil {
			return nil
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		status, orderID, err := orderTransition(msg.Topic, msg.Value)
		if err != nil {
			return err
		}
		if orderID != "" {
			if err := tx.Model(&domain.Order{}).Where("id = ?", orderID).Update("status", status).Error; err != nil {
				return err
			}
		}
		return tx.Create(&domain.InboxEvent{
			ID:          uuid.NewString(),
			MessageID:   messageID,
			EventType:   msg.Topic,
			ProcessedAt: time.Now(),
		}).Error
	})
}

func orderTransition(topic string, payload []byte) (domain.OrderStatus, string, error) {
	switch topic {
	case events.TopicWarehouseStockReserved:
		cloudEvent, err := events.ParseCloudEvent(payload)
		if err != nil {
			return "", "", err
		}
		event, err := events.DataAs[events.WarehouseStockReserved](cloudEvent)
		if err != nil {
			return "", "", err
		}
		return domain.OrderStatusPreparing, event.OrderID, nil
	case events.TopicPaymentIntentCreated:
		cloudEvent, err := events.ParseCloudEvent(payload)
		if err != nil {
			return "", "", err
		}
		event, err := events.DataAs[events.PaymentIntentCreated](cloudEvent)
		if err != nil {
			return "", "", err
		}
		return domain.OrderStatusPending, event.OrderID, nil
	case events.TopicPaymentCompleted, events.TopicPaymentSimulatedCompleted:
		cloudEvent, err := events.ParseCloudEvent(payload)
		if err != nil {
			return "", "", err
		}
		event, err := events.DataAs[events.PaymentCompleted](cloudEvent)
		if err != nil {
			return "", "", err
		}
		return domain.OrderStatusPending, event.OrderID, nil
	case events.TopicPaymentFailed:
		cloudEvent, err := events.ParseCloudEvent(payload)
		if err != nil {
			return "", "", err
		}
		event, err := events.DataAs[events.PaymentFailed](cloudEvent)
		if err != nil {
			return "", "", err
		}
		return domain.OrderStatusRejected, event.OrderID, nil
	case events.TopicPaymentRefunded:
		cloudEvent, err := events.ParseCloudEvent(payload)
		if err != nil {
			return "", "", err
		}
		event, err := events.DataAs[events.PaymentRefunded](cloudEvent)
		if err != nil {
			return "", "", err
		}
		return domain.OrderStatusRejected, event.OrderID, nil
	case events.TopicWarehouseStockReservationFailed:
		cloudEvent, err := events.ParseCloudEvent(payload)
		if err != nil {
			return "", "", err
		}
		event, err := events.DataAs[events.WarehouseStockReservationFailed](cloudEvent)
		if err != nil {
			return "", "", err
		}
		return domain.OrderStatusRejected, event.OrderID, nil
	case events.TopicLogisticsDeliveryAssigned:
		cloudEvent, err := events.ParseCloudEvent(payload)
		if err != nil {
			return "", "", err
		}
		event, err := events.DataAs[events.LogisticsShipmentAssigned](cloudEvent)
		if err != nil {
			return "", "", err
		}
		return domain.OrderStatusShipping, event.OrderID, nil
	case events.TopicLogisticsDeliveryCompleted:
		cloudEvent, err := events.ParseCloudEvent(payload)
		if err != nil {
			return "", "", err
		}
		event, err := events.DataAs[events.LogisticsShipmentDelivered](cloudEvent)
		if err != nil {
			return "", "", err
		}
		return domain.OrderStatusCompleted, event.OrderID, nil
	case events.TopicLogisticsDriverReturnedToBase:
		cloudEvent, err := events.ParseCloudEvent(payload)
		if err != nil {
			return "", "", err
		}
		event, err := events.DataAs[events.LogisticsDeliveryStatusChanged](cloudEvent)
		if err != nil {
			return "", "", err
		}
		return domain.OrderStatusCompleted, event.OrderID, nil
	default:
		return "", "", nil
	}
}
