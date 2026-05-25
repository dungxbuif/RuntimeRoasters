package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/dungxbuif/RuntimeRoasters/apps/logistics-service/internal/domain"
	"github.com/dungxbuif/RuntimeRoasters/pkg/base/identity"
	"github.com/dungxbuif/RuntimeRoasters/pkg/events"
	"github.com/dungxbuif/RuntimeRoasters/pkg/kafka"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	kafkago "github.com/segmentio/kafka-go"
	"gorm.io/gorm"
)

const driverLocationsKey = "drivers:locations"

type Service struct {
	db                        *gorm.DB
	rdb                       *redis.Client
	producer                  kafka.Producer
	shipmentAssignedTopic     string
	shipmentDeliveredTopic    string
	gpsUpdatedTopic           string
	defaultDestinationStoreID string
}

func NewService(db *gorm.DB, rdb *redis.Client, producer kafka.Producer, assignedTopic string, deliveredTopic string, gpsTopic string, defaultStoreID string) *Service {
	return &Service{db: db, rdb: rdb, producer: producer, shipmentAssignedTopic: assignedTopic, shipmentDeliveredTopic: deliveredTopic, gpsUpdatedTopic: gpsTopic, defaultDestinationStoreID: defaultStoreID}
}

func (s *Service) SeedDrivers(ctx context.Context) error {
	drivers := []domain.Driver{
		{ID: "22222222-2222-2222-2222-222222222201", Name: "Driver Hanoi 1", Phone: "+84010000001", Status: domain.DriverStatusIdle, IsAvailable: true},
		{ID: "22222222-2222-2222-2222-222222222202", Name: "Driver Hanoi 2", Phone: "+84010000002", Status: domain.DriverStatusIdle, IsAvailable: true},
		{ID: "22222222-2222-2222-2222-222222222203", Name: "Driver HCM 1", Phone: "+84010000003", Status: domain.DriverStatusIdle, IsAvailable: true},
		{ID: "22222222-2222-2222-2222-222222222204", Name: "Driver Da Nang 1", Phone: "+84010000004", Status: domain.DriverStatusIdle, IsAvailable: true},
		{ID: "22222222-2222-2222-2222-222222222205", Name: "Driver Backup", Phone: "+84010000005", Status: domain.DriverStatusIdle, IsAvailable: true},
	}
	for i, driver := range drivers {
		if err := s.db.WithContext(ctx).Where("id = ?", driver.ID).FirstOrCreate(&driver).Error; err != nil {
			return err
		}
		if err := s.storeDriverLocation(ctx, driver.ID, 21.0285+float64(i)*0.002, 105.8542+float64(i)*0.002); err != nil {
			return err
		}
	}
	return nil
}

func (s *Service) ListShipments(ctx context.Context) ([]domain.Shipment, error) {
	var shipments []domain.Shipment
	_, storeIDs, allStores := identity.StoreScopeFromContext(ctx)
	query := s.db.WithContext(ctx).Order("created_at DESC")
	if !allStores {
		if len(storeIDs) == 0 {
			return shipments, nil
		}
		query = query.Where("destination_store_id IN ?", storeIDs)
	}
	err := query.Find(&shipments).Error
	return shipments, err
}

func (s *Service) UpdateDriverLocation(ctx context.Context, driverID string, shipmentID string, lat float64, long float64) error {
	storeID := ""
	if claims, ok := identity.FromContext(ctx); ok && !claims.IsAdmin() {
		if shipmentID == "" {
			return errors.New("shipment_id is required")
		}
		var shipment domain.Shipment
		if err := s.db.WithContext(ctx).Where("id = ?", shipmentID).Take(&shipment).Error; err != nil {
			return err
		}
		if !claims.CanAccessStore(shipment.DestinationStoreID) {
			return gorm.ErrRecordNotFound
		}
		storeID = shipment.DestinationStoreID
	} else if shipmentID != "" {
		var shipment domain.Shipment
		if err := s.db.WithContext(ctx).Where("id = ?", shipmentID).Take(&shipment).Error; err == nil {
			storeID = shipment.DestinationStoreID
		}
	}
	if err := s.storeDriverLocation(ctx, driverID, lat, long); err != nil {
		return err
	}
	event := events.LogisticsGPSUpdated{EventID: uuid.NewString(), DriverID: driverID, ShipmentID: shipmentID, StoreID: storeID, Latitude: lat, Longitude: long, OccurredAt: time.Now()}
	cloudEvent, err := events.NewCloudEvent(ctx, s.gpsUpdatedTopic, events.SourceLogisticsService, fmt.Sprintf("drivers/%s", driverID), event, events.Metadata{
		EventID:       event.EventID,
		CorrelationID: firstNonEmpty(shipmentID, driverID),
		OccurredAt:    event.OccurredAt,
		DriverID:      driverID,
		ShipmentID:    shipmentID,
		StoreID:       storeID,
	})
	if err != nil {
		return err
	}
	return s.producer.Publish(ctx, s.gpsUpdatedTopic, driverID, cloudEvent)
}

func (s *Service) storeDriverLocation(ctx context.Context, driverID string, lat float64, long float64) error {
	if err := s.rdb.GeoAdd(ctx, driverLocationsKey, &redis.GeoLocation{Name: driverID, Longitude: long, Latitude: lat}).Err(); err != nil {
		return err
	}
	return s.rdb.Set(ctx, "driver:last_seen:"+driverID, time.Now().Format(time.RFC3339), 90*time.Second).Err()
}

func (s *Service) DeliverShipment(ctx context.Context, shipmentID string) error {
	now := time.Now()
	var shipment domain.Shipment
	if err := s.db.WithContext(ctx).Where("id = ?", shipmentID).Take(&shipment).Error; err != nil {
		return err
	}
	if claims, ok := identity.FromContext(ctx); ok && !claims.CanAccessStore(shipment.DestinationStoreID) {
		return gorm.ErrRecordNotFound
	}
	if err := s.db.WithContext(ctx).Model(&shipment).Updates(map[string]interface{}{"status": domain.ShipmentStatusDelivered, "delivered_at": &now}).Error; err != nil {
		return err
	}
	if shipment.DriverID != "" {
		_ = s.db.WithContext(ctx).Model(&domain.Driver{}).Where("id = ?", shipment.DriverID).Updates(map[string]interface{}{"status": domain.DriverStatusIdle, "is_available": true}).Error
	}
	event := events.LogisticsShipmentDelivered{EventID: uuid.NewString(), ShipmentID: shipment.ID, OrderID: shipment.OrderID, StoreID: shipment.DestinationStoreID, OccurredAt: now}
	return s.publishLogisticsEvent(ctx, s.shipmentDeliveredTopic, shipment.OrderID, event, firstNonEmpty(shipment.OrderID, shipment.ID), "", "")
}

func (s *Service) HandleWarehouseEvent(ctx context.Context, msg kafkago.Message) error {
	messageID := kafka.MessageID(msg)
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var existing domain.ProcessedKafkaMessage
		if err := tx.Where("msg_key = ?", messageID).Take(&existing).Error; err == nil {
			return nil
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		shipment, err := s.createShipmentFromEvent(ctx, tx, msg.Topic, msg.Value)
		if err != nil {
			return err
		}
		if shipment != nil {
			assigned, err := s.assignNearestDriver(ctx, tx, shipment)
			if err != nil {
				return err
			}
			if assigned != nil {
				event := events.LogisticsShipmentAssigned{EventID: uuid.NewString(), ShipmentID: assigned.ID, OrderID: assigned.OrderID, StoreID: assigned.DestinationStoreID, DriverID: assigned.DriverID, OccurredAt: time.Now()}
				correlationID := assigned.OrderID
				traceID := ""
				if cloudEvent, err := events.ParseCloudEvent(msg.Value); err == nil {
					correlationID = firstNonEmpty(events.ExtensionString(cloudEvent, "correlationid"), assigned.OrderID, assigned.ID)
					traceID = events.ExtensionString(cloudEvent, "traceid")
				}
				if err := s.publishLogisticsEvent(ctx, s.shipmentAssignedTopic, assigned.OrderID, event, correlationID, messageID, traceID); err != nil {
					return err
				}
			}
		}

		return tx.Create(&domain.ProcessedKafkaMessage{MsgKey: messageID, Topic: msg.Topic, ProcessedAt: time.Now()}).Error
	})
}

func (s *Service) createShipmentFromEvent(ctx context.Context, tx *gorm.DB, topic string, payload []byte) (*domain.Shipment, error) {
	cloudEvent, err := events.ParseCloudEvent(payload)
	if err != nil {
		return nil, err
	}
	switch topic {
	case events.TopicWarehouseStockReserved:
		event, err := events.DataAs[events.WarehouseStockReserved](cloudEvent)
		if err != nil {
			return nil, err
		}
		shipment := &domain.Shipment{ID: uuid.NewString(), OrderID: event.OrderID, OriginWarehouseID: "WAREHOUSE-001", DestinationStoreID: event.StoreID, DestinationAddress: event.StoreID, Status: domain.ShipmentStatusPending}
		return shipment, tx.Create(shipment).Error
	case events.TopicWarehouseInventoryUpdated:
		event, err := events.DataAs[events.WarehouseStockUpdated](cloudEvent)
		if err != nil {
			return nil, err
		}
		shipment := &domain.Shipment{ID: uuid.NewString(), BatchID: event.BatchID, OriginWarehouseID: "WAREHOUSE-001", DestinationStoreID: s.defaultDestinationStoreID, DestinationAddress: s.defaultDestinationStoreID, Status: domain.ShipmentStatusPending}
		return shipment, tx.Create(shipment).Error
	default:
		return nil, nil
	}
}

func (s *Service) assignNearestDriver(ctx context.Context, tx *gorm.DB, shipment *domain.Shipment) (*domain.Shipment, error) {
	// 1. Get Destination Location
	var loc domain.Location
	if err := tx.Where("id = ?", shipment.DestinationStoreID).Take(&loc).Error; err != nil {
		// Fallback to a default center point if location not found, but log it
		fmt.Printf("Warning: destination location %s not found, using default center\n", shipment.DestinationStoreID)
		loc = domain.Location{Lat: 21.0285, Lng: 105.8542} // Hanoi center
	}

	// 2. Search for nearest driver in Valkey
	locations, err := s.rdb.GeoSearch(ctx, driverLocationsKey, &redis.GeoSearchQuery{
		Longitude:  loc.Lng,
		Latitude:   loc.Lat,
		Radius:     50, // Increased radius to 50km for better simulation coverage
		RadiusUnit: "km",
		Count:      1,
		Sort:       "ASC",
	}).Result()

	if err != nil {
		return nil, err
	}

	if len(locations) == 0 {
		return nil, fmt.Errorf("no drivers available within 50km of %s", shipment.DestinationStoreID)
	}

	driverID := locations[0]

	// 3. Atomic Update: Driver Busy + Shipment Assigned
	if err := tx.Model(&domain.Driver{}).Where("id = ?", driverID).Updates(map[string]interface{}{
		"status":       domain.DriverStatusBusy,
		"is_available": false,
	}).Error; err != nil {
		return nil, err
	}

	if err := tx.Model(shipment).Updates(map[string]interface{}{
		"driver_id": driverID,
		"status":    domain.ShipmentStatusAssigned,
	}).Error; err != nil {
		return nil, err
	}

	shipment.DriverID = driverID
	shipment.Status = domain.ShipmentStatusAssigned
	return shipment, nil
}

func (s *Service) publishLogisticsEvent(ctx context.Context, topic string, key string, payload interface{}, correlationID string, causationID string, traceID string) error {
	metadata := logisticsMetadata(payload)
	if correlationID == "" {
		correlationID = firstNonEmpty(metadata.OrderID, metadata.ShipmentID, metadata.DriverID)
	}
	metadata.CorrelationID = correlationID
	metadata.CausationID = causationID
	metadata.TraceID = traceID
	cloudEvent, err := events.NewCloudEvent(ctx, topic, events.SourceLogisticsService, fmt.Sprintf("shipments/%s", metadata.ShipmentID), payload, metadata)
	if err != nil {
		return err
	}
	return s.producer.Publish(ctx, topic, key, cloudEvent)
}

func logisticsMetadata(payload interface{}) events.Metadata {
	switch event := payload.(type) {
	case events.LogisticsShipmentAssigned:
		return events.Metadata{EventID: event.EventID, OccurredAt: event.OccurredAt, OrderID: event.OrderID, StoreID: event.StoreID, ShipmentID: event.ShipmentID, DriverID: event.DriverID}
	case events.LogisticsShipmentDelivered:
		return events.Metadata{EventID: event.EventID, OccurredAt: event.OccurredAt, OrderID: event.OrderID, StoreID: event.StoreID, ShipmentID: event.ShipmentID}
	default:
		return events.Metadata{}
	}
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if value != "" {
			return value
		}
	}
	return ""
}
