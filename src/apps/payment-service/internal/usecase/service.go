package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"RuntimeRoasters/apps/payment-service/internal/domain"
	"RuntimeRoasters/apps/payment-service/internal/provider"
	"RuntimeRoasters/pkg/base/identity"
	"RuntimeRoasters/pkg/events"
	"RuntimeRoasters/pkg/kafka"
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
	cloudEvent, err := events.ParseCloudEvent(msg.Value)
	if err != nil {
		return err
	}
	event, err := events.DataAs[events.RetailOrderCreated](cloudEvent)
	if err != nil {
		return err
	}
	correlationID := events.ExtensionString(cloudEvent, "correlationid")
	if correlationID == "" {
		correlationID = event.OrderID
	}
	traceID := events.ExtensionString(cloudEvent, "traceid")

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
			if err := s.publishPaymentEvent(ctx, s.intentCreatedTopic, event.OrderID, createdEvent, correlationID, cloudEvent.ID(), traceID); err != nil {
				return err
			}
		}
		if completedEvent != nil {
			if err := s.publishPaymentEvent(ctx, s.paymentCompletedTopic, event.OrderID, completedEvent, correlationID, cloudEvent.ID(), traceID); err != nil {
				return err
			}
		}
		if failedEvent != nil {
			if err := s.publishPaymentEvent(ctx, s.paymentFailedTopic, event.OrderID, failedEvent, correlationID, cloudEvent.ID(), traceID); err != nil {
				return err
			}
		}
		return markProcessed(tx, messageID, msg.Topic)
	})
}

func (s *Service) HandleCompensationEvent(ctx context.Context, msg kafkago.Message) error {
	messageID := kafka.MessageID(msg)
	cloudEvent, err := events.ParseCloudEvent(msg.Value)
	if err != nil {
		return err
	}
	failed, err := events.DataAs[events.WarehouseStockReservationFailed](cloudEvent)
	if err != nil {
		return err
	}
	correlationID := events.ExtensionString(cloudEvent, "correlationid")
	if correlationID == "" {
		correlationID = failed.OrderID
	}
	traceID := events.ExtensionString(cloudEvent, "traceid")

	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if processed, err := alreadyProcessed(tx, messageID); processed || err != nil {
			return err
		}
		refunded, err := s.refund(ctx, tx, failed.OrderID, failed.Reason)
		if err != nil {
			return err
		}
		if refunded != nil {
			if err := s.publishPaymentEvent(ctx, s.paymentRefundedTopic, failed.OrderID, refunded, correlationID, cloudEvent.ID(), traceID); err != nil {
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
			return s.publishPaymentEvent(ctx, events.TopicPaymentCompleted, payment.OrderID, event, payment.OrderID, req.EventID, "")
		case domain.PaymentStatusFailed:
			if err := tx.Model(&payment).Update("status", domain.PaymentStatusFailed).Error; err != nil {
				return err
			}
			return s.publishPaymentEvent(ctx, s.paymentFailedTopic, payment.OrderID, failedEvent(payment, req.Reason), payment.OrderID, req.EventID, "")
		default:
			return fmt.Errorf("unsupported simulated webhook status %s", req.Status)
		}
	})
}

func (s *Service) publishPaymentEvent(ctx context.Context, topic string, key string, payload interface{}, correlationID string, causationID string, traceID string) error {
	metadata := paymentMetadata(payload)
	if metadata.EventID == "" {
		return errors.New("payment event id is required")
	}
	if correlationID == "" {
		correlationID = metadata.OrderID
	}
	metadata.CorrelationID = correlationID
	metadata.CausationID = causationID
	metadata.TraceID = traceID
	cloudEvent, err := events.NewCloudEvent(ctx, topic, events.SourcePaymentService, fmt.Sprintf("orders/%s", metadata.OrderID), payload, metadata)
	if err != nil {
		return err
	}
	return s.producer.Publish(ctx, topic, key, cloudEvent)
}

func paymentMetadata(payload interface{}) events.Metadata {
	switch event := payload.(type) {
	case *events.PaymentIntentCreated:
		return events.Metadata{EventID: event.EventID, OccurredAt: event.OccurredAt, OrderID: event.OrderID, PaymentID: event.PaymentID, StoreID: event.StoreID}
	case *events.PaymentCompleted:
		return events.Metadata{EventID: event.EventID, OccurredAt: event.OccurredAt, OrderID: event.OrderID, PaymentID: event.PaymentID, StoreID: event.StoreID}
	case *events.PaymentFailed:
		return events.Metadata{EventID: event.EventID, OccurredAt: event.OccurredAt, OrderID: event.OrderID, PaymentID: event.PaymentID, StoreID: event.StoreID}
	case *events.PaymentRefunded:
		return events.Metadata{EventID: event.EventID, OccurredAt: event.OccurredAt, OrderID: event.OrderID, PaymentID: event.PaymentID, StoreID: event.StoreID}
	default:
		return events.Metadata{}
	}
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
