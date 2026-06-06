package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
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
	"gorm.io/gorm/clause"
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
	counts := map[string]int64{}
	models := []struct {
		name  string
		model interface{}
	}{
		{"stores", &domain.Store{}},
		{"menus", &domain.Menu{}},
		{"menu_items", &domain.MenuItem{}},
		{"inventory_lots", &domain.InventoryLot{}},
		{"sales", &domain.Sale{}},
		{"sale_items", &domain.SaleItem{}},
		{"stock_movements", &domain.StockMovement{}},
		{"store_menu_inventories", &domain.StoreMenuInventory{}},
	}
	for _, entry := range models {
		var count int64
		if err := s.db.WithContext(ctx).Model(entry.model).Count(&count).Error; err != nil {
			return nil, err
		}
		counts[entry.name] = count
	}

	seeded := counts["stores"] == 5 &&
		counts["menus"] >= 1 &&
		counts["menu_items"] == 42 &&
		counts["inventory_lots"] >= 10 &&
		counts["sale_items"] >= 20

	return &domain.SystemStatus{
		Seeded:       seeded,
		ServiceName:  "retail-service",
		RecordCounts: counts,
	}, nil
}

func (s *Service) SeedData(ctx context.Context, force bool, usersMap map[string]string) (*domain.SeedResult, error) {
	stores, err := retailseed.LoadStores()
	if err != nil {
		return nil, err
	}
	menu, menuItems, err := retailseed.LoadMenu()
	if err != nil {
		return nil, err
	}
	demo, err := retailseed.BuildDemoData(stores, menuItems)
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

	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := upsertSeedRows(tx, stores); err != nil {
			return fmt.Errorf("seed stores: %w", err)
		}
		if err := upsertSeedRows(tx, []domain.Menu{menu}); err != nil {
			return fmt.Errorf("seed menus: %w", err)
		}
		if err := upsertSeedRows(tx, menuItems); err != nil {
			return fmt.Errorf("seed menu items: %w", err)
		}
		if err := insertSeedRows(tx, demo.Lots); err != nil {
			return fmt.Errorf("seed inventory lots: %w", err)
		}
		if err := insertSeedRows(tx, demo.Sales); err != nil {
			return fmt.Errorf("seed sales: %w", err)
		}
		if err := insertSeedRows(tx, demo.SaleItems); err != nil {
			return fmt.Errorf("seed sale items: %w", err)
		}
		if err := insertSeedRows(tx, demo.Movements); err != nil {
			return fmt.Errorf("seed stock movements: %w", err)
		}
		if err := rebuildLotBalances(ctx, tx); err != nil {
			return fmt.Errorf("rebuild inventory lot balances: %w", err)
		}
		for _, store := range stores {
			if err := rebuildStoreMenuInventory(ctx, tx, store.ID); err != nil {
				return fmt.Errorf("rebuild store menu inventory for %s: %w", store.Code, err)
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	records := int64(len(stores) + 1 + len(menuItems) + len(demo.Lots) + len(demo.Sales) + len(demo.SaleItems) + len(demo.Movements))
	return &domain.SeedResult{
		Success:        true,
		Message:        "Successfully reconciled retail master and demo data",
		RecordsCreated: records,
	}, nil
}

func upsertSeedRows[T any](tx *gorm.DB, rows []T) error {
	if len(rows) == 0 {
		return nil
	}
	return tx.Omit(clause.Associations).Clauses(clause.OnConflict{UpdateAll: true}).Create(&rows).Error
}

func insertSeedRows[T any](tx *gorm.DB, rows []T) error {
	if len(rows) == 0 {
		return nil
	}
	return tx.Omit(clause.Associations).Clauses(clause.OnConflict{DoNothing: true}).Create(&rows).Error
}

func rebuildLotBalances(ctx context.Context, tx *gorm.DB) error {
	var lots []domain.InventoryLot
	if err := tx.WithContext(ctx).Find(&lots).Error; err != nil {
		return err
	}
	type lotBalance struct {
		InventoryLotID string
		Balance        float64
	}
	var balances []lotBalance
	if err := tx.WithContext(ctx).Model(&domain.StockMovement{}).
		Select("inventory_lot_id, SUM(quantity_delta) AS balance").
		Group("inventory_lot_id").
		Scan(&balances).Error; err != nil {
		return err
	}
	byLot := make(map[string]float64, len(balances))
	for _, balance := range balances {
		byLot[balance.InventoryLotID] = balance.Balance
	}
	for _, lot := range lots {
		balance := byLot[lot.ID]
		if balance < 0 {
			return fmt.Errorf("inventory lot %s has negative ledger balance %.3f", lot.ID, balance)
		}
		updates := map[string]interface{}{"available_quantity": balance}
		if lot.Status != domain.InventoryLotStatusBlocked {
			status := domain.InventoryLotStatusAvailable
			if balance == 0 {
				status = domain.InventoryLotStatusDepleted
			}
			updates["status"] = status
		}
		if err := tx.WithContext(ctx).Model(&domain.InventoryLot{}).
			Where("id = ?", lot.ID).Updates(updates).Error; err != nil {
			return err
		}
	}
	return nil
}

func rebuildStoreMenuInventory(ctx context.Context, tx *gorm.DB, storeID string) error {
	var items []domain.MenuItem
	if err := tx.WithContext(ctx).
		Joins("JOIN menus ON menus.id = menu_items.menu_id").
		Where("menu_items.active = ? AND menus.status = ?", true, domain.MenuStatusActive).
		Find(&items).Error; err != nil {
		return err
	}

	var lots []domain.InventoryLot
	if err := tx.WithContext(ctx).
		Where("store_id = ? AND status = ?", storeID, domain.InventoryLotStatusAvailable).
		Find(&lots).Error; err != nil {
		return err
	}

	now := time.Now()
	for _, item := range items {
		var availableUnits int64
		var activeLotCount int64
		for _, lot := range lots {
			if lot.StockSKU != item.StockSKU || lot.AvailableQuantity <= 0 {
				continue
			}
			if lot.ExpiresAt != nil && !lot.ExpiresAt.After(now) {
				continue
			}
			units := int64(math.Floor(lot.AvailableQuantity / item.ConsumptionQuantity))
			if units <= 0 {
				continue
			}
			availableUnits += units
			activeLotCount++
		}
		row := domain.StoreMenuInventory{
			StoreID: storeID, MenuItemID: item.ID, AvailableUnits: availableUnits,
			ActiveLotCount: activeLotCount, Version: 1, UpdatedAt: now,
		}
		if err := tx.WithContext(ctx).Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "store_id"}, {Name: "menu_item_id"}},
			DoUpdates: clause.Assignments(map[string]interface{}{
				"available_units":  availableUnits,
				"active_lot_count": activeLotCount,
				"version":          gorm.Expr("store_menu_inventories.version + 1"),
				"updated_at":       now,
			}),
		}).Create(&row).Error; err != nil {
			return err
		}
	}
	return nil
}

func (s *Service) CreateStore(ctx context.Context, store domain.Store) (domain.Store, error) {
	if store.ID == "" {
		store.ID = uuid.NewString()
	}
	if store.Code == "" {
		codeID := uuid.NewSHA1(uuid.NameSpaceOID, []byte(store.Name+"|"+store.City)).String()
		store.Code = "STORE-" + codeID[:8]
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

func (s *Service) ConfirmOrderReceipt(ctx context.Context, id string) error {
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var order domain.Order
		if err := tx.Where("id = ?", id).Take(&order).Error; err != nil {
			return err
		}

		if order.Status != domain.OrderStatusDelivered {
			return fmt.Errorf("order %s cannot be confirmed from status %s", id, order.Status)
		}

		return tx.Model(&order).Update("status", domain.OrderStatusCompleted).Error
	})
}

func (s *Service) ListOrders(ctx context.Context) ([]domain.Order, error) {
	var orders []domain.Order
	_, storeIDs, allStores := identity.StoreScopeFromContext(ctx)
	query := s.db.WithContext(ctx).Order("created_at DESC")
	if !allStores {
		if len(storeIDs) == 0 {
			return orders, nil
		}
		query = query.Where("store_id IN ?", storeIDs)
	}
	err := query.Find(&orders).Error
	return orders, err
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
	case events.TopicPaymentIntentCreated:
		cloudEvent, err := events.ParseCloudEvent(payload)
		if err != nil {
			return "", "", err
		}
		event, err := events.DataAs[events.PaymentIntentCreated](cloudEvent)
		if err != nil {
			return "", "", err
		}
		return domain.OrderStatusPaymentPending, event.OrderID, nil

	case events.TopicPaymentCompleted, events.TopicPaymentSimulatedCompleted:
		cloudEvent, err := events.ParseCloudEvent(payload)
		if err != nil {
			return "", "", err
		}
		event, err := events.DataAs[events.PaymentCompleted](cloudEvent)
		if err != nil {
			return "", "", err
		}
		return domain.OrderStatusPaymentCompleted, event.OrderID, nil

	case events.TopicWarehouseStockReserved:
		cloudEvent, err := events.ParseCloudEvent(payload)
		if err != nil {
			return "", "", err
		}
		event, err := events.DataAs[events.WarehouseStockReserved](cloudEvent)
		if err != nil {
			return "", "", err
		}
		return domain.OrderStatusReserved, event.OrderID, nil

	case events.TopicWarehouseDispatchRequested:
		cloudEvent, err := events.ParseCloudEvent(payload)
		if err != nil {
			return "", "", err
		}
		event, err := events.DataAs[events.WarehouseDispatchRequested](cloudEvent)
		if err != nil {
			return "", "", err
		}
		return domain.OrderStatusDispatchRequested, event.OrderID, nil

	case events.TopicLogisticsDeliveryAssigned, events.TopicLogisticsDeliveryDeparted:
		cloudEvent, err := events.ParseCloudEvent(payload)
		if err != nil {
			return "", "", err
		}
		// Try parsing as assigned first, then status changed
		assigned, err := events.DataAs[events.LogisticsShipmentAssigned](cloudEvent)
		if err == nil {
			return domain.OrderStatusShipping, assigned.OrderID, nil
		}
		statusChanged, err := events.DataAs[events.LogisticsDeliveryStatusChanged](cloudEvent)
		if err == nil {
			return domain.OrderStatusShipping, statusChanged.OrderID, nil
		}
		return "", "", fmt.Errorf("failed to parse logistics event on topic %s", topic)

	case events.TopicLogisticsDeliveryDriverConfirmed:
		cloudEvent, err := events.ParseCloudEvent(payload)
		if err != nil {
			return "", "", err
		}
		event, err := events.DataAs[events.LogisticsDeliveryStatusChanged](cloudEvent)
		if err != nil {
			return "", "", err
		}
		return domain.OrderStatusDelivered, event.OrderID, nil

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

	case events.TopicPaymentFailed, events.TopicPaymentRefunded:
		cloudEvent, err := events.ParseCloudEvent(payload)
		if err != nil {
			return "", "", err
		}
		event, err := events.DataAs[events.PaymentFailed](cloudEvent)
		if err != nil {
			return "", "", err
		}
		return domain.OrderStatusFailed, event.OrderID, nil

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

	default:
		return "", "", nil
	}
}
