package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"RuntimeRoasters/apps/trace-service/internal/domain"
	"RuntimeRoasters/apps/trace-service/internal/search"
	"RuntimeRoasters/pkg/base/identity"
	rrevents "RuntimeRoasters/pkg/events"
	"RuntimeRoasters/pkg/kafka"
	"RuntimeRoasters/pkg/logger"
	"RuntimeRoasters/pkg/topology"
	cloudevents "github.com/cloudevents/sdk-go/v2"
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
	cloudEvent, err := rrevents.ParseCloudEvent(msg.Value)
	if err != nil {
		return err
	}
	ids := extractIDs(cloudEvent)
	traceID := ids["trace_id"]
	if spanContext := trace.SpanContextFromContext(ctx); spanContext.HasTraceID() {
		traceID = firstNonEmpty(traceID, spanContext.TraceID().String())
	}
	projection := topology.ProjectEvent(cloudEvent, msg.Value)
	displayPayload, err := json.Marshal(projection.DisplayPayload)
	if err != nil {
		return err
	}
	event := domain.TraceEvent{
		ID:             uuid.NewString(),
		MessageID:      messageID,
		Topic:          cloudEvent.Type(),
		BatchID:        ids["batch_id"],
		OrderID:        ids["order_id"],
		ShipmentID:     ids["shipment_id"],
		StoreID:        ids["store_id"],
		HarvestID:      ids["harvest_id"],
		FarmID:         ids["farm_id"],
		WarehouseID:    ids["warehouse_id"],
		DriverID:       ids["driver_id"],
		VehicleID:      ids["vehicle_id"],
		TraceID:        traceID,
		FlowID:         projection.FlowID,
		NodeID:         projection.NodeID,
		EdgeID:         projection.EdgeID,
		Pattern:        projection.Pattern,
		SourceService:  projection.SourceService,
		Visibility:     projection.Visibility,
		DisplayPayload: string(displayPayload),
		Payload:        string(msg.Value),
		OccurredAt:     cloudEvent.Time(),
	}

	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Clauses(clause.OnConflict{DoNothing: true}).Create(&event).Error; err != nil {
			return err
		}
		for _, key := range []string{"batch_id", "order_id", "shipment_id", "harvest_id"} {
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

func (s *Service) GetTopologyConfig(ctx context.Context, publicOnly bool) (topology.Config, error) {
	cfg := topology.CanonicalConfig()
	rows := []struct {
		FlowID     string
		LastSeenAt time.Time
	}{}
	query := s.db.WithContext(ctx).Model(&domain.TraceEvent{}).Select("flow_id, max(occurred_at) as last_seen_at").Where("flow_id <> ''").Group("flow_id")
	if publicOnly {
		query = query.Where("visibility = ?", topology.VisibilityPublic)
	}
	if err := query.Scan(&rows).Error; err != nil {
		return cfg, err
	}
	seen := map[string]string{}
	for _, row := range rows {
		if row.FlowID == "" || row.LastSeenAt.IsZero() {
			continue
		}
		seen[row.FlowID] = row.LastSeenAt.UTC().Format(time.RFC3339)
	}
	for i := range cfg.Flows {
		if lastSeen, ok := seen[cfg.Flows[i].ID]; ok {
			cfg.Flows[i].LastSeenAt = &lastSeen
		}
	}
	return cfg, nil
}

func (s *Service) GetTopologyHistory(ctx context.Context, filter TopologyHistoryFilter) ([]topology.HistoryEntry, error) {
	var events []domain.TraceEvent
	query := s.db.WithContext(ctx).Order("occurred_at DESC, created_at DESC").Limit(filter.limit())
	if filter.PublicOnly {
		query = query.Where("visibility = ?", topology.VisibilityPublic)
	}
	if filter.FlowID != "" {
		query = query.Where("flow_id = ?", filter.FlowID)
	}
	if filter.EntityID != "" {
		query = query.Where("batch_id = ? OR order_id = ? OR shipment_id = ? OR harvest_id = ?", filter.EntityID, filter.EntityID, filter.EntityID, filter.EntityID)
	}
	if filter.Cursor != "" {
		if cursorTime, err := time.Parse(time.RFC3339Nano, filter.Cursor); err == nil {
			query = query.Where("occurred_at < ?", cursorTime)
		}
	}
	if claims, ok := identity.FromContext(ctx); ok && !filter.PublicOnly && !claims.IsAdmin() {
		switch claims.Role {
		case identity.RoleStoreMgr:
			if len(claims.StoreIDs) == 0 {
				return []topology.HistoryEntry{}, nil
			}
			query = query.Where("store_id IN ? OR store_id = ''", claims.StoreIDs)
		case identity.RoleWarehouseMgr:
			if len(claims.WarehouseIDs) == 0 {
				return []topology.HistoryEntry{}, nil
			}
			query = query.Where("warehouse_id IN ? OR warehouse_id = ''", claims.WarehouseIDs)
		case "DRIVER":
			query = query.Where("driver_id = ? OR driver_id = ''", claims.Subject)
		default:
			query = query.Where("visibility = ?", topology.VisibilityPublic)
		}
	}
	if err := query.Find(&events).Error; err != nil {
		return nil, err
	}
	history := make([]topology.HistoryEntry, 0, len(events))
	for _, event := range events {
		var display map[string]interface{}
		_ = json.Unmarshal([]byte(event.DisplayPayload), &display)
		status := ""
		if value, ok := display["status"].(string); ok {
			status = value
		}
		history = append(history, topology.HistoryEntry{
			EventID:        event.MessageID,
			Topic:          event.Topic,
			Status:         status,
			FlowID:         event.FlowID,
			NodeID:         event.NodeID,
			EdgeID:         event.EdgeID,
			Pattern:        event.Pattern,
			SourceService:  event.SourceService,
			Visibility:     event.Visibility,
			TraceID:        event.TraceID,
			OccurredAt:     event.OccurredAt,
			DisplayPayload: display,
		})
	}
	return history, nil
}

type TopologyHistoryFilter struct {
	FlowID     string
	EntityID   string
	Cursor     string
	Limit      int
	PublicOnly bool
}

func (f TopologyHistoryFilter) limit() int {
	if f.Limit <= 0 || f.Limit > 200 {
		return 80
	}
	return f.Limit
}

func (s *Service) GetTrace(ctx context.Context, entityID string) ([]domain.TraceEvent, error) {
	var events []domain.TraceEvent
	err := s.db.WithContext(ctx).
		Where("batch_id = ? OR order_id = ? OR shipment_id = ? OR harvest_id = ?", entityID, entityID, entityID, entityID).
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
		case rrevents.TopicRetailOrderCreated:
			model.Order = payload
		case rrevents.TopicPaymentIntentCreated, rrevents.TopicPaymentCompleted, rrevents.TopicPaymentSimulatedCompleted, rrevents.TopicPaymentFailed, rrevents.TopicPaymentRefunded:
			model.Payment = merge(model.Payment, payload)
		case rrevents.TopicLogisticsDeliveryAssigned, rrevents.TopicLogisticsDeliveryCompleted:
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
	case rrevents.TopicRetailOrderCreated:
		return "ORDER_CREATED", "Order created"
	case rrevents.TopicPaymentIntentCreated:
		return "PAYMENT_PENDING", "Payment intent created"
	case rrevents.TopicPaymentCompleted, rrevents.TopicPaymentSimulatedCompleted:
		return "PAYMENT_COMPLETED", "Payment completed"
	case rrevents.TopicPaymentFailed:
		return "REJECTED", "Payment failed"
	case rrevents.TopicPaymentRefunded:
		return "REFUNDED", "Payment refunded"
	case rrevents.TopicWarehouseStockReserved:
		return "STOCK_RESERVED", "Stock reserved"
	case rrevents.TopicWarehouseStockReservationFailed:
		return "REJECTED", "Stock reservation failed"
	case rrevents.TopicWarehousePickupRequested:
		return "PICKUP_REQUESTED", "Warehouse pickup requested"
	case rrevents.TopicWarehousePickupReceived:
		return "PICKUP_RECEIVED", "Warehouse pickup received"
	case rrevents.TopicWarehouseDispatchRequested:
		return "DISPATCH_REQUESTED", "Warehouse dispatch requested"
	case rrevents.TopicWarehouseIntakeCreated:
		return "INTAKE_CREATED", "Warehouse intake created"
	case rrevents.TopicLogisticsDeliveryAssigned:
		return "SHIPPING", "Shipment assigned"
	case rrevents.TopicLogisticsDeliveryDeparted:
		return "DELIVERY_DEPARTED", "Delivery departed"
	case rrevents.TopicLogisticsDeliveryArrivedAtStore:
		return "ARRIVED_AT_STORE", "Delivery arrived at store"
	case rrevents.TopicLogisticsDeliveryDriverConfirmed:
		return "DELIVERY_CONFIRMED", "Driver confirmed delivery"
	case rrevents.TopicLogisticsDeliveryCompleted:
		return "COMPLETED", "Shipment delivered"
	case rrevents.TopicLogisticsPickupAssigned:
		return "PICKUP_ASSIGNED", "Pickup assigned"
	case rrevents.TopicLogisticsPickupDeparted:
		return "PICKUP_DEPARTED", "Pickup departed"
	case rrevents.TopicLogisticsPickupArrivedAtFarm:
		return "ARRIVED_AT_FARM", "Pickup arrived at farm"
	case rrevents.TopicLogisticsPickupLoadingConfirmed:
		return "LOADING_CONFIRMED", "Pickup loading confirmed"
	case rrevents.TopicLogisticsPickupReturnStarted:
		return "PICKUP_RETURN_STARTED", "Pickup return started"
	case rrevents.TopicLogisticsPickupArrivedAtWarehouse:
		return "ARRIVED_AT_WAREHOUSE", "Pickup arrived at warehouse"
	case rrevents.TopicLogisticsPickupCompleted:
		return "PICKUP_COMPLETED", "Pickup completed"
	case rrevents.TopicLogisticsDriverReturnStarted:
		return "DRIVER_RETURN_STARTED", "Driver return started"
	case rrevents.TopicLogisticsDriverReturnCompleted, rrevents.TopicLogisticsDriverReturnedToBase:
		return "DRIVER_RETURNED", "Driver returned to base"
	case rrevents.TopicLogisticsGPSUpdated:
		return "IN_TRANSIT", "Driver location updated"
	case rrevents.TopicLogisticsShipmentStatusChanged:
		return "SHIPMENT_STATUS_CHANGED", "Shipment status changed"
	case rrevents.TopicNotificationCreated:
		return "NOTIFICATION_CREATED", "Notification created"
	case rrevents.TopicNotificationAcknowledged:
		return "NOTIFICATION_ACKNOWLEDGED", "Notification acknowledged"
	case rrevents.TopicSocketBroadcastRequested:
		return "SOCKET_BROADCAST_REQUESTED", "Socket broadcast requested"
	default:
		return "UPDATED", topic
	}
}

func eventDescription(topic string, payload map[string]interface{}) string {
	switch topic {
	case rrevents.TopicPaymentIntentCreated:
		return fmt.Sprintf("%s simulated checkout created for %.2f %s", stringValue(payload, "provider"), floatValue(payload, "amount"), stringValue(payload, "currency"))
	case rrevents.TopicPaymentCompleted, rrevents.TopicPaymentSimulatedCompleted:
		return fmt.Sprintf("%s payment completed for %.2f %s", stringValue(payload, "provider"), floatValue(payload, "amount"), stringValue(payload, "currency"))
	case rrevents.TopicPaymentRefunded:
		return fmt.Sprintf("%s refund issued because %s", stringValue(payload, "provider"), stringValue(payload, "reason"))
	case rrevents.TopicWarehouseStockReserved:
		return fmt.Sprintf("Reserved %.0f units of %s", floatValue(payload, "quantity"), stringValue(payload, "sku"))
	case rrevents.TopicWarehouseStockReservationFailed:
		return fmt.Sprintf("Stock reservation failed: %s", stringValue(payload, "reason"))
	case rrevents.TopicWarehousePickupRequested:
		return fmt.Sprintf("Pickup requested for harvest %s", stringValue(payload, "harvest_id"))
	case rrevents.TopicWarehousePickupReceived:
		return fmt.Sprintf("Pickup %s received at warehouse %s", stringValue(payload, "pickup_id"), stringValue(payload, "warehouse_id"))
	case rrevents.TopicWarehouseDispatchRequested:
		return fmt.Sprintf("Dispatch %s requested from warehouse %s", stringValue(payload, "dispatch_id"), stringValue(payload, "warehouse_id"))
	case rrevents.TopicWarehouseIntakeCreated:
		return fmt.Sprintf("Intake %s created for harvest %s", stringValue(payload, "intake_id"), stringValue(payload, "harvest_id"))
	case rrevents.TopicLogisticsDeliveryAssigned:
		return fmt.Sprintf("Shipment %s assigned to driver %s", stringValue(payload, "shipment_id"), stringValue(payload, "driver_id"))
	case rrevents.TopicLogisticsDeliveryCompleted:
		return fmt.Sprintf("Shipment %s delivered", stringValue(payload, "shipment_id"))
	case rrevents.TopicLogisticsPickupAssigned, rrevents.TopicLogisticsPickupDeparted, rrevents.TopicLogisticsPickupArrivedAtFarm, rrevents.TopicLogisticsPickupLoadingConfirmed, rrevents.TopicLogisticsPickupReturnStarted, rrevents.TopicLogisticsPickupArrivedAtWarehouse, rrevents.TopicLogisticsPickupCompleted:
		return fmt.Sprintf("Pickup shipment %s status %s", stringValue(payload, "shipment_id"), stringValue(payload, "status"))
	case rrevents.TopicLogisticsDeliveryDeparted, rrevents.TopicLogisticsDeliveryArrivedAtStore, rrevents.TopicLogisticsDeliveryDriverConfirmed, rrevents.TopicLogisticsDriverReturnStarted, rrevents.TopicLogisticsDriverReturnCompleted, rrevents.TopicLogisticsDriverReturnedToBase:
		return fmt.Sprintf("Delivery shipment %s status %s", stringValue(payload, "shipment_id"), stringValue(payload, "status"))
	case rrevents.TopicLogisticsShipmentStatusChanged:
		return fmt.Sprintf("Shipment %s status %s", stringValue(payload, "shipment_id"), stringValue(payload, "status"))
	case rrevents.TopicNotificationCreated:
		return fmt.Sprintf("Notification created: %s", stringValue(payload, "title"))
	case rrevents.TopicNotificationAcknowledged:
		return fmt.Sprintf("Notification %s acknowledged", stringValue(payload, "notification_id"))
	case rrevents.TopicSocketBroadcastRequested:
		return fmt.Sprintf("Socket broadcast requested for %s", stringValue(payload, "channel"))
	default:
		return topic
	}
}

func parsePayload(payload string) map[string]interface{} {
	cloudEvent, err := rrevents.ParseCloudEvent([]byte(payload))
	if err != nil {
		return map[string]interface{}{}
	}
	return rrevents.DataMap(cloudEvent)
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

func extractIDs(event cloudevents.Event) map[string]string {
	ids := map[string]string{
		"trace_id":     rrevents.ExtensionString(event, "traceid"),
		"batch_id":     rrevents.ExtensionString(event, "batchid"),
		"order_id":     rrevents.ExtensionString(event, "orderid"),
		"shipment_id":  rrevents.ExtensionString(event, "shipmentid"),
		"harvest_id":   rrevents.ExtensionString(event, "harvestid"),
		"store_id":     rrevents.ExtensionString(event, "storeid"),
		"farm_id":      rrevents.ExtensionString(event, "farmid"),
		"warehouse_id": rrevents.ExtensionString(event, "warehouseid"),
		"driver_id":    rrevents.ExtensionString(event, "driverid"),
		"vehicle_id":   rrevents.ExtensionString(event, "vehicleid"),
	}
	data := rrevents.DataMap(event)
	for _, key := range []string{"batch_id", "order_id", "shipment_id", "harvest_id", "store_id", "farm_id", "warehouse_id", "driver_id", "vehicle_id"} {
		if ids[key] == "" {
			ids[key] = stringValue(data, key)
		}
	}
	return ids
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if value != "" {
			return value
		}
	}
	return ""
}
