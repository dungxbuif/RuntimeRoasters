package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/dungxbuif/RuntimeRoasters/apps/payment-service/internal/domain"
	"github.com/dungxbuif/RuntimeRoasters/apps/payment-service/internal/provider"
	"github.com/dungxbuif/RuntimeRoasters/pkg/base/identity"
	"github.com/dungxbuif/RuntimeRoasters/pkg/events"
	"github.com/dungxbuif/RuntimeRoasters/pkg/kafka"
	"github.com/google/uuid"
	kafkago "github.com/segmentio/kafka-go"
	"gorm.io/gorm"
)

type Service struct {
	db                    *gorm.DB
	producer              kafka.Producer
	providers             *provider.Factory
	defaultProvider       string
	defaultCurrency       string
	backfillOrders        bool
	startedAt             time.Time
	intentCreatedTopic    string
	paymentCompletedTopic string
	paymentFailedTopic    string
	paymentRefundedTopic  string
}

func NewService(db *gorm.DB, producer kafka.Producer, providers *provider.Factory, defaultProvider string, defaultCurrency string, backfillOrders bool, intentCreatedTopic string, paymentCompletedTopic string, paymentFailedTopic string, paymentRefundedTopic string) *Service {
	return &Service{
		db:                    db,
		producer:              producer,
		providers:             providers,
		defaultProvider:       defaultProvider,
		defaultCurrency:       defaultCurrency,
		backfillOrders:        backfillOrders,
		startedAt:             time.Now(),
		intentCreatedTopic:    intentCreatedTopic,
		paymentCompletedTopic: paymentCompletedTopic,
		paymentFailedTopic:    paymentFailedTopic,
		paymentRefundedTopic:  paymentRefundedTopic,
	}
}

func (s *Service) HandleOrderCreated(ctx context.Context, msg kafkago.Message) error {
	messageID := kafka.MessageID(msg)
	var event events.RetailOrderCreated
	if err := json.Unmarshal(msg.Value, &event); err != nil {
		return err
	}

	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if processed, err := alreadyProcessed(tx, messageID); processed || err != nil {
			return err
		}
		if !s.backfillOrders && event.OccurredAt.Before(s.startedAt.Add(-5*time.Second)) {
			return markProcessed(tx, messageID, msg.Topic)
		}
		payment, createdEvent, completedEvent, failedEvent, err := s.createPayment(ctx, tx, event)
		if err != nil {
			return err
		}
		_ = payment
		if createdEvent != nil {
			if err := s.producer.Publish(ctx, s.intentCreatedTopic, event.OrderID, createdEvent); err != nil {
				return err
			}
		}
		if completedEvent != nil {
			if err := s.producer.Publish(ctx, s.paymentCompletedTopic, event.OrderID, completedEvent); err != nil {
				return err
			}
		}
		if failedEvent != nil {
			if err := s.producer.Publish(ctx, s.paymentFailedTopic, event.OrderID, failedEvent); err != nil {
				return err
			}
		}
		return markProcessed(tx, messageID, msg.Topic)
	})
}

func (s *Service) HandleCompensationEvent(ctx context.Context, msg kafkago.Message) error {
	messageID := kafka.MessageID(msg)
	var failed events.WarehouseStockReservationFailed
	if err := json.Unmarshal(msg.Value, &failed); err != nil {
		return err
	}

	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if processed, err := alreadyProcessed(tx, messageID); processed || err != nil {
			return err
		}
		refunded, err := s.refund(ctx, tx, failed.OrderID, failed.Reason)
		if err != nil {
			return err
		}
		if refunded != nil {
			if err := s.producer.Publish(ctx, s.paymentRefundedTopic, failed.OrderID, refunded); err != nil {
				return err
			}
		}
		return markProcessed(tx, messageID, msg.Topic)
	})
}

func (s *Service) ListPayments(ctx context.Context) ([]domain.Payment, error) {
	var payments []domain.Payment
	_, storeIDs, allStores := identity.StoreScopeFromContext(ctx)
	query := s.db.WithContext(ctx).Order("created_at DESC")
	if !allStores {
		if len(storeIDs) == 0 {
			return payments, nil
		}
		query = query.Where("store_id IN ?", storeIDs)
	}
	err := query.Find(&payments).Error
	return payments, err
}

func (s *Service) GetPaymentByOrder(ctx context.Context, orderID string) (*domain.Payment, error) {
	var payment domain.Payment
	if err := s.db.WithContext(ctx).Where("order_id = ?", orderID).Take(&payment).Error; err != nil {
		return nil, err
	}
	if claims, ok := identity.FromContext(ctx); ok && !claims.CanAccessStore(payment.StoreID) {
		return nil, gorm.ErrRecordNotFound
	}
	return &payment, nil
}

func (s *Service) SimulateWebhook(ctx context.Context, providerName string, payload []byte, timestamp string, signature string) error {
	gateway, err := s.providers.Get(providerName)
	if err != nil {
		return err
	}
	if err := gateway.VerifyWebhook(payload, timestamp, signature); err != nil {
		return err
	}
	var req struct {
		EventID     string `json:"event_id"`
		ProviderRef string `json:"provider_ref"`
		Status      string `json:"status"`
		Reason      string `json:"reason"`
	}
	if err := json.Unmarshal(payload, &req); err != nil {
		return err
	}
	if req.EventID == "" {
		return errors.New("event_id is required")
	}
	if req.ProviderRef == "" {
		return errors.New("provider_ref is required")
	}
	status := strings.ToUpper(req.Status)
	if status == "" {
		status = string(domain.PaymentStatusSucceeded)
	}

	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var existingEvent domain.WebhookEvent
		if err := tx.Where("provider = ? AND event_id = ?", gateway.Name(), req.EventID).Take(&existingEvent).Error; err == nil {
			return nil
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		var payment domain.Payment
		if err := tx.Where("provider_ref = ?", req.ProviderRef).Take(&payment).Error; err != nil {
			return err
		}
		if err := tx.Create(&domain.WebhookEvent{
			ID:          uuid.NewString(),
			Provider:    gateway.Name(),
			EventID:     req.EventID,
			ProviderRef: req.ProviderRef,
			Status:      status,
			Reason:      req.Reason,
			ProcessedAt: time.Now(),
		}).Error; err != nil {
			return err
		}
		switch domain.PaymentStatus(status) {
		case domain.PaymentStatusSucceeded:
			if err := tx.Model(&payment).Update("status", domain.PaymentStatusSucceeded).Error; err != nil {
				return err
			}
			event, err := s.completedEvent(tx, payment.ID)
			if err != nil {
				return err
			}
			return s.producer.Publish(ctx, s.paymentCompletedTopic, payment.OrderID, event)
		case domain.PaymentStatusFailed:
			if err := tx.Model(&payment).Update("status", domain.PaymentStatusFailed).Error; err != nil {
				return err
			}
			return s.producer.Publish(ctx, s.paymentFailedTopic, payment.OrderID, failedEvent(payment, req.Reason))
		default:
			return fmt.Errorf("unsupported simulated webhook status %s", req.Status)
		}
	})
}

func (s *Service) createPayment(ctx context.Context, tx *gorm.DB, event events.RetailOrderCreated) (*domain.Payment, *events.PaymentIntentCreated, *events.PaymentCompleted, *events.PaymentFailed, error) {
	var existing domain.Payment
	if err := tx.Where("order_id = ?", event.OrderID).Take(&existing).Error; err == nil {
		return &existing, nil, nil, nil, nil
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil, nil, nil, err
	}

	providerName := event.PaymentMethod
	if providerName == "" {
		providerName = s.defaultProvider
	}
	currency := s.defaultCurrency
	gateway, err := s.providers.Get(providerName)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	paymentID := uuid.NewString()
	intent, err := gateway.CreateIntent(ctx, provider.IntentRequest{PaymentID: paymentID, OrderID: event.OrderID, Amount: event.TotalAmount, Currency: currency})
	if err != nil {
		return nil, nil, nil, nil, err
	}
	itemsJSON, err := json.Marshal(event.Items)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	status := domain.PaymentStatusSucceeded
	if event.TotalAmount <= 0 {
		status = domain.PaymentStatusFailed
	}
	payment := &domain.Payment{
		ID:          paymentID,
		OrderID:     event.OrderID,
		StoreID:     event.StoreID,
		Provider:    gateway.Name(),
		ProviderRef: intent.ProviderRef,
		Items:       string(itemsJSON),
		Amount:      event.TotalAmount,
		Currency:    currency,
		Status:      status,
		CheckoutURL: intent.CheckoutURL,
		Simulated:   intent.Simulated,
	}
	if err := tx.Create(payment).Error; err != nil {
		return nil, nil, nil, nil, err
	}
	intentEvent := &events.PaymentIntentCreated{
		EventID:     uuid.NewString(),
		PaymentID:   payment.ID,
		OrderID:     payment.OrderID,
		StoreID:     payment.StoreID,
		Provider:    payment.Provider,
		ProviderRef: payment.ProviderRef,
		Items:       event.Items,
		Amount:      payment.Amount,
		Currency:    payment.Currency,
		CheckoutURL: payment.CheckoutURL,
		Simulated:   payment.Simulated,
		OccurredAt:  time.Now(),
	}
	if status == domain.PaymentStatusFailed {
		return payment, intentEvent, nil, failedEvent(*payment, "invalid amount"), nil
	}
	completed, err := s.completedEvent(tx, payment.ID)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	return payment, intentEvent, completed, nil, nil
}

func (s *Service) completedEvent(tx *gorm.DB, paymentID string) (*events.PaymentCompleted, error) {
	var payment domain.Payment
	if err := tx.Where("id = ?", paymentID).Take(&payment).Error; err != nil {
		return nil, err
	}
	var items []events.RetailOrderItem
	if err := json.Unmarshal([]byte(payment.Items), &items); err != nil {
		return nil, err
	}
	return &events.PaymentCompleted{
		EventID:     uuid.NewString(),
		PaymentID:   payment.ID,
		OrderID:     payment.OrderID,
		StoreID:     payment.StoreID,
		Provider:    payment.Provider,
		ProviderRef: payment.ProviderRef,
		Items:       items,
		Amount:      payment.Amount,
		Currency:    payment.Currency,
		OccurredAt:  time.Now(),
	}, nil
}

func (s *Service) refund(ctx context.Context, tx *gorm.DB, orderID string, reason string) (*events.PaymentRefunded, error) {
	var payment domain.Payment
	if err := tx.Where("order_id = ?", orderID).Take(&payment).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	if payment.Status == domain.PaymentStatusRefunded {
		return nil, nil
	}
	gateway, err := s.providers.Get(payment.Provider)
	if err != nil {
		return nil, err
	}
	refund, err := gateway.Refund(ctx, payment.ProviderRef, payment.Amount, payment.Currency, reason)
	if err != nil {
		return nil, err
	}
	if err := tx.Model(&payment).Updates(map[string]interface{}{"status": domain.PaymentStatusRefunded, "refund_ref": refund.RefundRef}).Error; err != nil {
		return nil, err
	}
	return &events.PaymentRefunded{
		EventID:     uuid.NewString(),
		PaymentID:   payment.ID,
		OrderID:     payment.OrderID,
		StoreID:     payment.StoreID,
		Provider:    payment.Provider,
		ProviderRef: payment.ProviderRef,
		RefundRef:   refund.RefundRef,
		Reason:      reason,
		Amount:      payment.Amount,
		Currency:    payment.Currency,
		OccurredAt:  time.Now(),
	}, nil
}

func failedEvent(payment domain.Payment, reason string) *events.PaymentFailed {
	if reason == "" {
		reason = "payment failed"
	}
	return &events.PaymentFailed{
		EventID:     uuid.NewString(),
		PaymentID:   payment.ID,
		OrderID:     payment.OrderID,
		StoreID:     payment.StoreID,
		Provider:    payment.Provider,
		ProviderRef: payment.ProviderRef,
		Reason:      reason,
		OccurredAt:  time.Now(),
	}
}

func alreadyProcessed(tx *gorm.DB, messageID string) (bool, error) {
	var existing domain.InboxEvent
	if err := tx.Where("message_id = ?", messageID).Take(&existing).Error; err == nil {
		return true, nil
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return false, err
	}
	return false, nil
}

func markProcessed(tx *gorm.DB, messageID string, eventType string) error {
	return tx.Create(&domain.InboxEvent{ID: uuid.NewString(), MessageID: messageID, EventType: eventType, ProcessedAt: time.Now()}).Error
}
