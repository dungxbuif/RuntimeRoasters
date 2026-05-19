package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/dungxbuif/RuntimeRoasters/apps/retail-service/internal/domain"
	"github.com/dungxbuif/RuntimeRoasters/pkg/base/identity"
	"github.com/dungxbuif/RuntimeRoasters/pkg/events"
	"github.com/dungxbuif/RuntimeRoasters/pkg/kafka"
	"github.com/dungxbuif/RuntimeRoasters/pkg/logger"
	"github.com/google/uuid"
	kafkago "github.com/segmentio/kafka-go"
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

func (s *Service) SeedStores(ctx context.Context) error {
	stores := []domain.Store{
		{ID: "11111111-1111-1111-1111-111111111101", Name: "Hoan Kiem Store", City: "Hanoi", Address: "2 Ly Thai To", ManagerEmail: "mgr.hn.hoankiem@runtimeroasters.com", Status: domain.StoreStatusActive},
		{ID: "11111111-1111-1111-1111-111111111102", Name: "Cau Giay Store", City: "Hanoi", Address: "102 Tran Thai Tong", ManagerEmail: "mgr.hn.caugiay@runtimeroasters.com", Status: domain.StoreStatusActive},
		{ID: "11111111-1111-1111-1111-111111111103", Name: "District 1 Store", City: "Ho Chi Minh City", Address: "45 Le Thanh Ton", ManagerEmail: "mgr.hcm.q1@runtimeroasters.com", Status: domain.StoreStatusActive},
		{ID: "11111111-1111-1111-1111-111111111104", Name: "District 7 Store", City: "Ho Chi Minh City", Address: "Phu My Hung", ManagerEmail: "mgr.hcm.q7@runtimeroasters.com", Status: domain.StoreStatusActive},
		{ID: "11111111-1111-1111-1111-111111111105", Name: "Hai Chau Store", City: "Da Nang", Address: "15 Bach Dang", ManagerEmail: "mgr.dn.haichau@runtimeroasters.com", Status: domain.StoreStatusActive},
	}
	for _, store := range stores {
		if err := s.db.WithContext(ctx).Where("id = ?", store.ID).FirstOrCreate(&store).Error; err != nil {
			return err
		}
	}
	return nil
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
	payload, err := json.Marshal(event)
	if err != nil {
		return nil, err
	}

	order := &domain.Order{
		ID:             orderID,
		StoreID:        req.StoreID,
		Items:          string(itemsJSON),
		TotalAmount:    req.TotalAmount,
		Status:         domain.OrderStatusPending,
		IdempotencyKey: idempotencyKey,
	}
	outbox := &domain.OutboxEvent{
		ID:        uuid.NewString(),
		EventType: s.orderCreatedTopic,
		Topic:     s.orderCreatedTopic,
		Key:       orderID,
		Payload:   string(payload),
		Status:    domain.OutboxStatusPending,
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
		if err := s.producer.Publish(ctx, event.Topic, event.Key, json.RawMessage(event.Payload)); err != nil {
			logger.FromContext(ctx).Warn("failed to publish retail outbox event", zap.Error(err), zap.String("event_id", event.ID))
			_ = s.db.WithContext(ctx).Model(&event).Update("status", domain.OutboxStatusFailed).Error
			continue
		}
		if err := s.db.WithContext(ctx).Model(&event).Update("status", domain.OutboxStatusCompleted).Error; err != nil {
			return err
		}
	}
	return nil
}

func (s *Service) HandleSagaEvent(ctx context.Context, msg kafkago.Message) error {
	messageID := fmt.Sprintf("%s-%d-%d", msg.Topic, msg.Partition, msg.Offset)
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
		var event events.WarehouseStockReserved
		if err := json.Unmarshal(payload, &event); err != nil {
			return "", "", err
		}
		return domain.OrderStatusPreparing, event.OrderID, nil
	case events.TopicPaymentIntentCreated:
		var event events.PaymentIntentCreated
		if err := json.Unmarshal(payload, &event); err != nil {
			return "", "", err
		}
		return domain.OrderStatusPending, event.OrderID, nil
	case events.TopicPaymentCompleted:
		var event events.PaymentCompleted
		if err := json.Unmarshal(payload, &event); err != nil {
			return "", "", err
		}
		return domain.OrderStatusPending, event.OrderID, nil
	case events.TopicPaymentFailed:
		var event events.PaymentFailed
		if err := json.Unmarshal(payload, &event); err != nil {
			return "", "", err
		}
		return domain.OrderStatusRejected, event.OrderID, nil
	case events.TopicPaymentRefunded:
		var event events.PaymentRefunded
		if err := json.Unmarshal(payload, &event); err != nil {
			return "", "", err
		}
		return domain.OrderStatusRejected, event.OrderID, nil
	case events.TopicWarehouseStockReservationFailed:
		var event events.WarehouseStockReservationFailed
		if err := json.Unmarshal(payload, &event); err != nil {
			return "", "", err
		}
		return domain.OrderStatusRejected, event.OrderID, nil
	case events.TopicLogisticsShipmentAssigned:
		var event events.LogisticsShipmentAssigned
		if err := json.Unmarshal(payload, &event); err != nil {
			return "", "", err
		}
		return domain.OrderStatusShipping, event.OrderID, nil
	case events.TopicLogisticsShipmentDelivered:
		var event events.LogisticsShipmentDelivered
		if err := json.Unmarshal(payload, &event); err != nil {
			return "", "", err
		}
		return domain.OrderStatusCompleted, event.OrderID, nil
	default:
		return "", "", nil
	}
}
