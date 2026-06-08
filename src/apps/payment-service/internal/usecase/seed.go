package usecase

import (
	"context"
	"encoding/json"
	"time"

	"RuntimeRoasters/apps/payment-service/internal/domain"
	"RuntimeRoasters/apps/payment-service/internal/provider"
	"RuntimeRoasters/pkg/events"
	"RuntimeRoasters/pkg/logger"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const demoSeedPrefix = "rr31-demo"

var demoSeedReferenceTime = time.Date(2026, 6, 8, 8, 0, 0, 0, time.UTC)

type demoPaymentSeed struct {
	PaymentID   string
	OrderID     string
	StoreID     string
	Provider    string
	ProviderRef string
	RefundRef   string
	Status      domain.PaymentStatus
	Amount      float64
	Currency    string
	Items       []events.RetailOrderItem
	Reason      string
	OccurredAt  time.Time
}

func (s *Service) SeedDemoData(ctx context.Context, force bool, usersMap map[string]string) (int64, error) {
	_ = usersMap
	if force {
		logger.GetLogger().Warn("payment demo seed force reset requested; RR-31 keeps normal seed non-destructive")
	}
	seeds := demoPaymentSeeds()
	itemsJSON := func(items []events.RetailOrderItem) string {
		payload, _ := json.Marshal(items)
		return string(payload)
	}
	created := int64(0)
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, seed := range seeds {
			payment := domain.Payment{
				ID:          seed.PaymentID,
				OrderID:     seed.OrderID,
				StoreID:     seed.StoreID,
				Provider:    seed.Provider,
				ProviderRef: seed.ProviderRef,
				RefundRef:   seed.RefundRef,
				Items:       itemsJSON(seed.Items),
				Amount:      seed.Amount,
				Currency:    seed.Currency,
				Status:      seed.Status,
				CheckoutURL: "https://payments.runtime-roasters.local/demo/" + seed.ProviderRef,
				Simulated:   true,
				CreatedAt:   seed.OccurredAt,
				UpdatedAt:   seed.OccurredAt,
			}
			result := tx.Clauses(clause.OnConflict{
				Columns: []clause.Column{{Name: "provider_ref"}},
				DoUpdates: clause.AssignmentColumns([]string{
					"order_id", "store_id", "provider", "refund_ref", "items", "amount", "currency", "status", "checkout_url", "simulated", "updated_at",
				}),
			}).Create(&payment)
			if result.Error != nil {
				return result.Error
			}
			created += result.RowsAffected
			if err := seedWebhookEvent(tx, seed); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return 0, err
	}
	for _, seed := range seeds {
		if err := s.publishSeedPaymentEvents(ctx, seed); err != nil {
			logger.GetLogger().Warn("payment demo seed event publish failed", zap.String("order_id", seed.OrderID), zap.Error(err))
		}
	}
	return created, nil
}

func (s *Service) DemoPaymentCount(ctx context.Context) (int64, error) {
	var count int64
	err := s.db.WithContext(ctx).Model(&domain.Payment{}).Where("provider_ref LIKE ?", demoSeedPrefix+"%").Count(&count).Error
	return count, err
}

func seedWebhookEvent(tx *gorm.DB, seed demoPaymentSeed) error {
	status := string(seed.Status)
	event := domain.WebhookEvent{
		ID:          stableUUID("webhook", seed.ProviderRef),
		Provider:    seed.Provider,
		EventID:     demoSeedPrefix + "-webhook-" + seed.ProviderRef,
		ProviderRef: seed.ProviderRef,
		Status:      status,
		Reason:      seed.Reason,
		ProcessedAt: seed.OccurredAt.Add(3 * time.Minute),
	}
	return tx.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "provider"}, {Name: "event_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"provider_ref", "status", "reason", "processed_at"}),
	}).Create(&event).Error
}

func (s *Service) publishSeedPaymentEvents(ctx context.Context, seed demoPaymentSeed) error {
	intent := events.PaymentIntentCreated{
		EventID:     demoSeedPrefix + "-intent-" + seed.OrderID,
		PaymentID:   seed.PaymentID,
		OrderID:     seed.OrderID,
		StoreID:     seed.StoreID,
		Provider:    seed.Provider,
		ProviderRef: seed.ProviderRef,
		Items:       seed.Items,
		Amount:      seed.Amount,
		Currency:    seed.Currency,
		CheckoutURL: "https://payments.runtime-roasters.local/demo/" + seed.ProviderRef,
		Simulated:   true,
		OccurredAt:  seed.OccurredAt.Add(time.Minute),
	}
	if err := s.publishPaymentEvent(ctx, s.intentCreatedTopic, seed.OrderID, &intent, seed.OrderID, demoSeedPrefix, demoSeedPrefix+"-"+seed.OrderID); err != nil {
		return err
	}
	switch seed.Status {
	case domain.PaymentStatusSucceeded:
		completed := events.PaymentCompleted{
			EventID:     demoSeedPrefix + "-completed-" + seed.OrderID,
			PaymentID:   seed.PaymentID,
			OrderID:     seed.OrderID,
			StoreID:     seed.StoreID,
			Provider:    seed.Provider,
			ProviderRef: seed.ProviderRef,
			Items:       seed.Items,
			Amount:      seed.Amount,
			Currency:    seed.Currency,
			OccurredAt:  seed.OccurredAt.Add(3 * time.Minute),
		}
		return s.publishPaymentEvent(ctx, s.paymentCompletedTopic, seed.OrderID, &completed, seed.OrderID, intent.EventID, demoSeedPrefix+"-"+seed.OrderID)
	case domain.PaymentStatusFailed:
		failed := events.PaymentFailed{
			EventID:     demoSeedPrefix + "-failed-" + seed.OrderID,
			PaymentID:   seed.PaymentID,
			OrderID:     seed.OrderID,
			StoreID:     seed.StoreID,
			Provider:    seed.Provider,
			ProviderRef: seed.ProviderRef,
			Reason:      seed.Reason,
			OccurredAt:  seed.OccurredAt.Add(3 * time.Minute),
		}
		return s.publishPaymentEvent(ctx, s.paymentFailedTopic, seed.OrderID, &failed, seed.OrderID, intent.EventID, demoSeedPrefix+"-"+seed.OrderID)
	case domain.PaymentStatusRefunded:
		refunded := events.PaymentRefunded{
			EventID:     demoSeedPrefix + "-refunded-" + seed.OrderID,
			PaymentID:   seed.PaymentID,
			OrderID:     seed.OrderID,
			StoreID:     seed.StoreID,
			Provider:    seed.Provider,
			ProviderRef: seed.ProviderRef,
			RefundRef:   seed.RefundRef,
			Reason:      seed.Reason,
			Amount:      seed.Amount,
			Currency:    seed.Currency,
			OccurredAt:  seed.OccurredAt.Add(10 * time.Minute),
		}
		return s.publishPaymentEvent(ctx, s.paymentRefundedTopic, seed.OrderID, &refunded, seed.OrderID, intent.EventID, demoSeedPrefix+"-"+seed.OrderID)
	default:
		return nil
	}
}

func demoPaymentSeeds() []demoPaymentSeed {
	return []demoPaymentSeed{
		{
			PaymentID: stableUUID("payment", "succeeded"),
			OrderID:   stableUUID("order", "succeeded"),
			StoreID:   "11111111-1111-1111-1111-111111111103",
			Provider:  provider.ProviderStripe, ProviderRef: demoSeedPrefix + "-pi-succeeded",
			Status: domain.PaymentStatusSucceeded, Amount: 1200000, Currency: "VND",
			Items:      []events.RetailOrderItem{{SKU: "SKU-AR-VN-LD-001", Quantity: 10}},
			OccurredAt: demoSeedReferenceTime.Add(-48 * time.Hour),
		},
		{
			PaymentID: stableUUID("payment", "failed"),
			OrderID:   stableUUID("order", "failed"),
			StoreID:   "11111111-1111-1111-1111-111111111102",
			Provider:  provider.ProviderStripe, ProviderRef: demoSeedPrefix + "-pi-failed",
			Status: domain.PaymentStatusFailed, Amount: 430000, Currency: "VND", Reason: "card_declined",
			Items:      []events.RetailOrderItem{{SKU: "SKU-AR-VN-LD-001", Quantity: 3}},
			OccurredAt: demoSeedReferenceTime.Add(-36 * time.Hour),
		},
		{
			PaymentID: stableUUID("payment", "refunded"),
			OrderID:   stableUUID("order", "refunded"),
			StoreID:   "11111111-1111-1111-1111-111111111101",
			Provider:  provider.ProviderStripe, ProviderRef: demoSeedPrefix + "-pi-refunded", RefundRef: demoSeedPrefix + "-re-quality",
			Status: domain.PaymentStatusRefunded, Amount: 980000, Currency: "VND", Reason: "quality_issue",
			Items:      []events.RetailOrderItem{{SKU: "SKU-RB-VN-DL-001", Quantity: 8}},
			OccurredAt: demoSeedReferenceTime.Add(-24 * time.Hour),
		},
	}
}

func stableUUID(parts ...string) string {
	return uuid.NewSHA1(uuid.NameSpaceOID, []byte(demoSeedPrefix+"|"+joinSeedParts(parts))).String()
}

func joinSeedParts(parts []string) string {
	joined := ""
	for i, part := range parts {
		if i > 0 {
			joined += "|"
		}
		joined += part
	}
	return joined
}
