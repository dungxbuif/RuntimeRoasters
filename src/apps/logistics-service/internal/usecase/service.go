package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"RuntimeRoasters/apps/logistics-service/internal/domain"
	logisticsseed "RuntimeRoasters/apps/logistics-service/internal/seed"
	"RuntimeRoasters/pkg/base/identity"
	"RuntimeRoasters/pkg/events"
	"RuntimeRoasters/pkg/kafka"
	"RuntimeRoasters/pkg/logger"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	kafkago "github.com/segmentio/kafka-go"
	"go.uber.org/zap"
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
	pickupAssignedTopic       string
	pickupArrivedTopic        string
	statusChangedTopic        string
	driverReturnedTopic       string
	defaultDestinationStoreID string
	defaultWarehouseID        string
}

func NewService(db *gorm.DB, rdb *redis.Client, producer kafka.Producer, assignedTopic string, deliveredTopic string, gpsTopic string, defaultStoreID string) *Service {
	return &Service{
		db:                        db,
		rdb:                       rdb,
		producer:                  producer,
		shipmentAssignedTopic:     assignedTopic,
		shipmentDeliveredTopic:    deliveredTopic,
		gpsUpdatedTopic:           gpsTopic,
		pickupAssignedTopic:       events.TopicLogisticsPickupAssigned,
		pickupArrivedTopic:        events.TopicLogisticsPickupArrivedAtWarehouse,
		statusChangedTopic:        events.TopicLogisticsShipmentStatusChanged,
		driverReturnedTopic:       events.TopicLogisticsDriverReturnedToBase,
		defaultDestinationStoreID: defaultStoreID,
		defaultWarehouseID:        "WAREHOUSE-HN-001",
	}
}

func (s *Service) SeedDrivers(ctx context.Context, vehicles []domain.Vehicle, drivers []domain.Driver) error {
	if len(vehicles) == 0 {
		return errors.New("logistics vehicle seed is empty")
	}
	if len(drivers) == 0 {
		return errors.New("logistics driver seed is empty")
	}
	for _, vehicle := range vehicles {
		if err := s.db.WithContext(ctx).Save(&vehicle).Error; err != nil {
			return err
		}
	}

	for i, driver := range drivers {
		if err := s.db.WithContext(ctx).Save(&driver).Error; err != nil {
			return err
		}
		if err := s.storeDriverLocation(ctx, driver.ID, 21.0285+float64(i)*0.002, 105.8542+float64(i)*0.002); err != nil {
			return err
		}
	}
	return nil
}

func (s *Service) SeedLocations(ctx context.Context, locations []domain.Location) error {
	if len(locations) == 0 {
		return errors.New("logistics location seed is empty")
	}
	for _, location := range locations {
		if err := s.db.WithContext(ctx).Save(&location).Error; err != nil {
			return err
		}
	}
	return nil
}

func (s *Service) GetSystemStatus(ctx context.Context) (*domain.SystemStatus, error) {
	var driverCount, locationCount, vehicleCount int64
	s.db.WithContext(ctx).Model(&domain.Driver{}).Count(&driverCount)
	s.db.WithContext(ctx).Model(&domain.Location{}).Count(&locationCount)
	s.db.WithContext(ctx).Model(&domain.Vehicle{}).Count(&vehicleCount)

	seeded := driverCount > 0 && locationCount > 0 && vehicleCount > 0

	return &domain.SystemStatus{
		Seeded:      seeded,
		ServiceName: "logistics-service",
		RecordCounts: map[string]int64{
			"drivers":   driverCount,
			"locations": locationCount,
			"vehicles":  vehicleCount,
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

	seedData, err := logisticsseed.Load()
	if err != nil {
		return nil, err
	}

	// Enrich drivers with Kratos UserID from usersMap
	for i := range seedData.Drivers {
		if id, ok := usersMap[seedData.Drivers[i].UserID]; ok {
			seedData.Drivers[i].UserID = id
		} else {
			// Fallback or warning if email not found in Kratos
			logger.GetLogger().Warn("Driver email not found in users_map during seeding",
				zap.String("email", seedData.Drivers[i].UserID),
				zap.String("driver", seedData.Drivers[i].Name))
		}
	}

	if err := s.SeedDrivers(ctx, seedData.Vehicles, seedData.Drivers); err != nil {
		return nil, err
	}
	if err := s.SeedLocations(ctx, seedData.Locations); err != nil {
		return nil, err
	}

	total := int64(len(seedData.Vehicles) + len(seedData.Drivers) + len(seedData.Locations))

	return &domain.SeedResult{
		Success:        true,
		Message:        "Successfully seeded logistics data",
		RecordsCreated: total,
	}, nil
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

func (s *Service) GetShipment(ctx context.Context, id string) (*domain.Shipment, error) {
	var shipment domain.Shipment
	if err := s.db.WithContext(ctx).Where("id = ?", id).Take(&shipment).Error; err != nil {
		return nil, err
	}
	if claims, ok := identity.FromContext(ctx); ok && !claims.IsAdmin() {
		if claims.Role == identity.RoleStoreMgr && !claims.CanAccessStore(shipment.DestinationStoreID) {
			return nil, gorm.ErrRecordNotFound
		}
		if claims.Role == identity.RoleWarehouseMgr && !claims.CanAccessWarehouse(shipment.WarehouseID) {
			return nil, gorm.ErrRecordNotFound
		}
		if claims.Role == "DRIVER" {
			driver, err := s.driverFromClaims(ctx, claims)
			if err != nil || driver.ID != shipment.DriverID {
				return nil, gorm.ErrRecordNotFound
			}
		}
	}
	return &shipment, nil
}

func (s *Service) ListLocations(ctx context.Context) ([]domain.Location, error) {
	var locations []domain.Location
	err := s.db.WithContext(ctx).Order("type, name").Find(&locations).Error
	return locations, err
}

func (s *Service) ListDrivers(ctx context.Context) ([]domain.Driver, error) {
	var drivers []domain.Driver
	err := s.db.WithContext(ctx).Order("name").Find(&drivers).Error
	return drivers, err
}

func (s *Service) ListVehicles(ctx context.Context, availableOnly bool) ([]domain.Vehicle, error) {
	var vehicles []domain.Vehicle
	query := s.db.WithContext(ctx).Order("plate_number")
	if availableOnly {
		query = query.Where("status = ?", domain.DriverStatusIdle)
	}
	err := query.Find(&vehicles).Error
	return vehicles, err
}

func (s *Service) UpdateDriverLocation(ctx context.Context, driverID string, shipmentID string, lat float64, long float64) error {
	claims, hasClaims := identity.FromContext(ctx)
	if hasClaims && claims.Role == "DRIVER" {
		driver, err := s.driverFromClaims(ctx, claims)
		if err != nil {
			return err
		}
		if driverID != "" && driverID != driver.ID {
			return errors.New("driver cannot update another driver")
		}
		driverID = driver.ID
	}
	if driverID == "" {
		return errors.New("driver_id is required")
	}
	if _, err := uuid.Parse(driverID); err != nil {
		return errors.New("driver_id must be a UUID")
	}
	storeID := ""
	if hasClaims && !claims.IsAdmin() {
		if shipmentID == "" {
			return errors.New("shipment_id is required")
		}
		var shipment domain.Shipment
		if err := s.db.WithContext(ctx).Where("id = ?", shipmentID).Take(&shipment).Error; err != nil {
			return err
		}
		switch claims.Role {
		case "DRIVER":
			if shipment.DriverID != driverID {
				return errors.New("driver cannot update another shipment")
			}
		case identity.RoleStoreMgr:
			if !claims.CanAccessStore(shipment.DestinationStoreID) {
				return gorm.ErrRecordNotFound
			}
		case identity.RoleWarehouseMgr:
			if !claims.CanAccessWarehouse(shipment.WarehouseID) {
				return gorm.ErrRecordNotFound
			}
		}
		storeID = firstNonEmpty(shipment.DestinationStoreID, shipment.StoreID)
	} else if shipmentID != "" {
		var shipment domain.Shipment
		if err := s.db.WithContext(ctx).Where("id = ?", shipmentID).Take(&shipment).Error; err == nil {
			storeID = firstNonEmpty(shipment.DestinationStoreID, shipment.StoreID)
		}
	}
	if err := s.storeDriverLocation(ctx, driverID, lat, long); err != nil {
		return err
	}
	_ = s.db.WithContext(ctx).Model(&domain.Vehicle{}).
		Where("id = (SELECT vehicle_id FROM drivers WHERE id = ?)", driverID).
		Updates(map[string]interface{}{"current_latitude": lat, "current_longitude": long}).Error
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
	return s.ConfirmDelivery(ctx, shipmentID)
}

func (s *Service) AssignShipment(ctx context.Context, shipmentID string, driverID string, vehicleID string) (*domain.Shipment, error) {
	var shipment domain.Shipment
	now := time.Now()
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("id = ?", shipmentID).Take(&shipment).Error; err != nil {
			return err
		}
		if claims, ok := identity.FromContext(ctx); ok && claims.Role == identity.RoleWarehouseMgr && !claims.CanAccessWarehouse(shipment.WarehouseID) {
			return gorm.ErrRecordNotFound
		}
		driver, vehicle, err := s.resolveDriverVehicle(ctx, tx, driverID, vehicleID)
		if err != nil {
			return err
		}
		if err := assignShipment(tx, &shipment, driver.ID, vehicle.ID, &now); err != nil {
			return err
		}
		if err := markDriverBusy(tx, driver.ID, vehicle.ID, shipment.ID); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	if err := s.publishAssigned(ctx, shipment); err != nil {
		return nil, err
	}
	return &shipment, nil
}

func (s *Service) DepartShipment(ctx context.Context, shipmentID string) (*domain.Shipment, error) {
	return s.transitionDriverAction(ctx, shipmentID, "depart")
}

func (s *Service) ArriveShipment(ctx context.Context, shipmentID string) (*domain.Shipment, error) {
	return s.transitionDriverAction(ctx, shipmentID, "arrive")
}

func (s *Service) ConfirmLoad(ctx context.Context, shipmentID string) (*domain.Shipment, error) {
	return s.transitionDriverAction(ctx, shipmentID, "confirm-load")
}

func (s *Service) ConfirmDelivery(ctx context.Context, shipmentID string) error {
	_, err := s.transitionDriverAction(ctx, shipmentID, "confirm-delivery")
	return err
}

func (s *Service) ReturnShipment(ctx context.Context, shipmentID string) (*domain.Shipment, error) {
	return s.transitionDriverAction(ctx, shipmentID, "return")
}

func (s *Service) transitionDriverAction(ctx context.Context, shipmentID string, action string) (*domain.Shipment, error) {
	var shipment domain.Shipment
	var topic string
	now := time.Now()
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("id = ?", shipmentID).Take(&shipment).Error; err != nil {
			return err
		}
		if err := s.authorizeShipmentAction(ctx, tx, shipment); err != nil {
			return err
		}

		updates, eventTopic, err := transitionUpdates(shipment, action, now)
		if err != nil {
			return err
		}
		topic = eventTopic
		if err := tx.Model(&shipment).Updates(updates).Error; err != nil {
			return err
		}
		if status, ok := updates["status"].(domain.ShipmentStatus); ok {
			shipment.Status = status
		}
		if leg, ok := updates["current_leg"].(domain.ShipmentLeg); ok {
			shipment.CurrentLeg = leg
		}
		if action == "arrive" && (shipment.Status == domain.ShipmentStatusArrivedWarehouse || shipment.Status == domain.ShipmentStatusReturnedToBase) {
			if err := tx.Model(&domain.Driver{}).Where("id = ?", shipment.DriverID).Updates(map[string]interface{}{
				"status":              domain.DriverStatusIdle,
				"is_available":        true,
				"current_shipment_id": "",
			}).Error; err != nil {
				return err
			}
			if shipment.VehicleID != "" {
				if err := tx.Model(&domain.Vehicle{}).Where("id = ?", shipment.VehicleID).Update("status", domain.DriverStatusIdle).Error; err != nil {
					return err
				}
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	if topic != "" {
		if err := s.publishStatusChanged(ctx, topic, shipment); err != nil {
			return nil, err
		}
	}
	return &shipment, nil
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
			assigned := shipment
			var err error
			if shipment.Status != domain.ShipmentStatusAssigned {
				assigned, err = s.assignNearestDriver(ctx, tx, shipment)
				if err != nil {
					return err
				}
			}
			if assigned != nil {
				correlationID := assigned.OrderID
				traceID := ""
				if cloudEvent, err := events.ParseCloudEvent(msg.Value); err == nil {
					correlationID = firstNonEmpty(events.ExtensionString(cloudEvent, "correlationid"), assigned.OrderID, assigned.ID)
					traceID = events.ExtensionString(cloudEvent, "traceid")
				}
				if err := s.publishAssignedWithMeta(ctx, *assigned, correlationID, messageID, traceID); err != nil {
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
	case events.TopicWarehousePickupRequested:
		event, err := events.DataAs[events.WarehousePickupRequested](cloudEvent)
		if err != nil {
			return nil, err
		}
		if event.Status != "DISPATCHED" {
			return nil, nil
		}
		var existing domain.Shipment
		if err := tx.Where("type = ? AND harvest_id = ?", domain.ShipmentTypeFarmPickup, event.HarvestID).Take(&existing).Error; err == nil {
			return &existing, nil
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		shipment := &domain.Shipment{
			ID:                    uuid.NewString(),
			Type:                  domain.ShipmentTypeFarmPickup,
			Status:                domain.ShipmentStatusPending,
			CurrentLeg:            domain.ShipmentLegOutbound,
			RouteID:               "route-warehouse-farm-demo",
			OriginLocationID:      event.WarehouseID,
			DestinationLocationID: event.FarmID,
			FarmID:                event.FarmID,
			HarvestID:             event.HarvestID,
			WarehouseID:           event.WarehouseID,
			OriginWarehouseID:     event.WarehouseID,
			DestinationAddress:    event.FarmID,
		}
		if err := tx.Create(shipment).Error; err != nil {
			return nil, err
		}
		return s.assignNearestDriver(ctx, tx, shipment)
	case events.TopicLogisticsDeliveryAssigned:
		event, err := events.DataAs[events.LogisticsShipmentAssigned](cloudEvent)
		if err != nil {
			return nil, err
		}
		now := time.Now()
		shipment := &domain.Shipment{
			ID:                    uuid.NewString(),
			Type:                  domain.ShipmentTypeRetailDelivery,
			Status:                domain.ShipmentStatusAssigned, // Already assigned by Warehouse
			CurrentLeg:            domain.ShipmentLegOutbound,
			RouteID:               "route-warehouse-store-demo",
			OrderID:               event.OrderID,
			DriverID:              event.DriverID,
			VehicleID:             event.VehicleID,
			WarehouseID:           firstNonEmpty(events.ExtensionString(cloudEvent, "warehouseid"), s.defaultWarehouseID),
			OriginWarehouseID:     firstNonEmpty(events.ExtensionString(cloudEvent, "warehouseid"), s.defaultWarehouseID),
			OriginLocationID:      firstNonEmpty(events.ExtensionString(cloudEvent, "warehouseid"), s.defaultWarehouseID),
			DestinationStoreID:    event.StoreID,
			StoreID:               event.StoreID,
			DestinationLocationID: event.StoreID,
			DestinationAddress:    event.StoreID,
			AssignedAt:            &now,
		}
		if err := tx.Create(shipment).Error; err != nil {
			return nil, err
		}
		if err := markDriverBusy(tx, event.DriverID, event.VehicleID, shipment.ID); err != nil {
			return nil, err
		}
		return shipment, nil
	case events.TopicWarehouseInventoryUpdated:
		event, err := events.DataAs[events.WarehouseStockUpdated](cloudEvent)
		if err != nil {
			return nil, err
		}
		shipment := &domain.Shipment{
			ID:                    uuid.NewString(),
			Type:                  domain.ShipmentTypeRetailDelivery,
			Status:                domain.ShipmentStatusPending,
			CurrentLeg:            domain.ShipmentLegOutbound,
			RouteID:               "route-warehouse-store-demo",
			BatchID:               event.BatchID,
			WarehouseID:           s.defaultWarehouseID,
			OriginWarehouseID:     s.defaultWarehouseID,
			OriginLocationID:      s.defaultWarehouseID,
			DestinationStoreID:    s.defaultDestinationStoreID,
			StoreID:               s.defaultDestinationStoreID,
			DestinationLocationID: s.defaultDestinationStoreID,
			DestinationAddress:    s.defaultDestinationStoreID,
		}
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
	var driver domain.Driver
	if err := tx.Where("id = ? AND is_available = ?", driverID, true).Take(&driver).Error; err != nil {
		return nil, err
	}

	// 3. Atomic Update: Driver Busy + Shipment Assigned
	now := time.Now()
	if err := markDriverBusy(tx, driver.ID, driver.VehicleID, shipment.ID); err != nil {
		return nil, err
	}

	if err := assignShipment(tx, shipment, driver.ID, driver.VehicleID, &now); err != nil {
		return nil, err
	}
	return shipment, nil
}

func (s *Service) resolveDriverVehicle(ctx context.Context, tx *gorm.DB, driverID string, vehicleID string) (*domain.Driver, *domain.Vehicle, error) {
	var driver domain.Driver
	query := tx.WithContext(ctx).Where("is_available = ?", true)
	if driverID != "" {
		query = query.Where("id = ?", driverID)
	}
	if err := query.Order("created_at").Take(&driver).Error; err != nil {
		return nil, nil, err
	}
	if vehicleID == "" {
		vehicleID = driver.VehicleID
	}
	var vehicle domain.Vehicle
	if vehicleID != "" {
		if err := tx.WithContext(ctx).Where("id = ?", vehicleID).Take(&vehicle).Error; err != nil {
			return nil, nil, err
		}
	} else {
		if err := tx.WithContext(ctx).Where("status = ?", domain.DriverStatusIdle).Order("created_at").Take(&vehicle).Error; err != nil {
			return nil, nil, err
		}
	}
	return &driver, &vehicle, nil
}

func assignShipment(tx *gorm.DB, shipment *domain.Shipment, driverID string, vehicleID string, now *time.Time) error {
	if err := tx.Model(shipment).Updates(map[string]interface{}{
		"driver_id":   driverID,
		"vehicle_id":  vehicleID,
		"status":      domain.ShipmentStatusAssigned,
		"assigned_at": now,
	}).Error; err != nil {
		return err
	}
	shipment.DriverID = driverID
	shipment.VehicleID = vehicleID
	shipment.Status = domain.ShipmentStatusAssigned
	shipment.AssignedAt = now
	return nil
}

func markDriverBusy(tx *gorm.DB, driverID string, vehicleID string, shipmentID string) error {
	if err := tx.Model(&domain.Driver{}).Where("id = ?", driverID).Updates(map[string]interface{}{
		"status":              domain.DriverStatusBusy,
		"is_available":        false,
		"vehicle_id":          vehicleID,
		"current_shipment_id": shipmentID,
	}).Error; err != nil {
		return err
	}
	if vehicleID == "" {
		return nil
	}
	return tx.Model(&domain.Vehicle{}).Where("id = ?", vehicleID).Update("status", domain.DriverStatusBusy).Error
}

func (s *Service) driverFromClaims(ctx context.Context, claims identity.Claims) (*domain.Driver, error) {
	return s.driverFromClaimsDB(ctx, s.db, claims)
}

func (s *Service) driverFromClaimsDB(ctx context.Context, db *gorm.DB, claims identity.Claims) (*domain.Driver, error) {
	var driver domain.Driver
	if err := db.WithContext(ctx).
		Where("user_id IN ? OR id = ?", []string{claims.Subject, claims.Email}, claims.Subject).
		Take(&driver).Error; err != nil {
		return nil, err
	}
	return &driver, nil
}

func (s *Service) authorizeShipmentAction(ctx context.Context, db *gorm.DB, shipment domain.Shipment) error {
	claims, ok := identity.FromContext(ctx)
	if !ok || claims.IsAdmin() {
		return nil
	}
	switch claims.Role {
	case "DRIVER":
		driver, err := s.driverFromClaimsDB(ctx, db, claims)
		if err != nil {
			return err
		}
		if driver.ID != shipment.DriverID {
			return errors.New("driver cannot update another shipment")
		}
	case identity.RoleWarehouseMgr:
		if !claims.CanAccessWarehouse(shipment.WarehouseID) {
			return gorm.ErrRecordNotFound
		}
	case identity.RoleStoreMgr:
		if !claims.CanAccessStore(shipment.DestinationStoreID) {
			return gorm.ErrRecordNotFound
		}
	default:
		return errors.New("role cannot update shipment")
	}
	return nil
}

func transitionUpdates(shipment domain.Shipment, action string, now time.Time) (map[string]interface{}, string, error) {
	updates := map[string]interface{}{}
	switch shipment.Type {
	case domain.ShipmentTypeFarmPickup:
		switch action {
		case "depart":
			if shipment.Status != domain.ShipmentStatusAssigned {
				return nil, "", fmt.Errorf("cannot depart farm pickup from %s", shipment.Status)
			}
			updates["status"] = domain.ShipmentStatusInTransitToFarm
			updates["current_leg"] = domain.ShipmentLegOutbound
			updates["departed_at"] = &now
			return updates, events.TopicLogisticsPickupDeparted, nil
		case "arrive":
			if shipment.Status == domain.ShipmentStatusInTransitToFarm {
				updates["status"] = domain.ShipmentStatusArrivedAtFarm
				updates["arrived_at_origin_at"] = &now
				return updates, events.TopicLogisticsPickupArrivedAtFarm, nil
			}
			if shipment.Status == domain.ShipmentStatusReturningToWarehouse {
				updates["status"] = domain.ShipmentStatusArrivedWarehouse
				updates["arrived_destination_at"] = &now
				updates["returned_at"] = &now
				return updates, events.TopicLogisticsPickupArrivedAtWarehouse, nil
			}
			return nil, "", fmt.Errorf("cannot arrive farm pickup from %s", shipment.Status)
		case "confirm-load":
			if shipment.Status != domain.ShipmentStatusArrivedAtFarm {
				return nil, "", fmt.Errorf("cannot confirm load from %s", shipment.Status)
			}
			updates["status"] = domain.ShipmentStatusPickedUp
			updates["loaded_at"] = &now
			return updates, events.TopicLogisticsPickupLoadingConfirmed, nil
		case "return":
			if shipment.Status != domain.ShipmentStatusPickedUp {
				return nil, "", fmt.Errorf("cannot start return from %s", shipment.Status)
			}
			updates["status"] = domain.ShipmentStatusReturningToWarehouse
			updates["current_leg"] = domain.ShipmentLegReturn
			updates["return_started_at"] = &now
			return updates, events.TopicLogisticsPickupReturnStarted, nil
		default:
			return nil, "", fmt.Errorf("action %s is invalid for farm pickup", action)
		}
	case domain.ShipmentTypeRetailDelivery:
		switch action {
		case "depart":
			if shipment.Status != domain.ShipmentStatusAssigned {
				return nil, "", fmt.Errorf("cannot depart retail delivery from %s", shipment.Status)
			}
			updates["status"] = domain.ShipmentStatusInTransitToStore
			updates["current_leg"] = domain.ShipmentLegOutbound
			updates["departed_at"] = &now
			return updates, events.TopicLogisticsDeliveryDeparted, nil
		case "arrive":
			if shipment.Status == domain.ShipmentStatusInTransitToStore {
				updates["status"] = domain.ShipmentStatusArrivedAtStore
				updates["arrived_destination_at"] = &now
				return updates, events.TopicLogisticsDeliveryArrivedAtStore, nil
			}
			if shipment.Status == domain.ShipmentStatusReturningToBase {
				updates["status"] = domain.ShipmentStatusReturnedToBase
				updates["returned_at"] = &now
				return updates, events.TopicLogisticsDriverReturnedToBase, nil
			}
			return nil, "", fmt.Errorf("cannot arrive retail delivery from %s", shipment.Status)
		case "confirm-delivery":
			if shipment.Status != domain.ShipmentStatusArrivedAtStore {
				return nil, "", fmt.Errorf("cannot confirm delivery from %s", shipment.Status)
			}
			updates["status"] = domain.ShipmentStatusDelivered
			updates["delivered_at"] = &now
			return updates, events.TopicLogisticsDeliveryDriverConfirmed, nil
		case "return":
			if shipment.Status != domain.ShipmentStatusDelivered {
				return nil, "", fmt.Errorf("cannot start return from %s", shipment.Status)
			}
			updates["status"] = domain.ShipmentStatusReturningToBase
			updates["current_leg"] = domain.ShipmentLegReturn
			updates["return_started_at"] = &now
			return updates, events.TopicLogisticsDriverReturnStarted, nil
		default:
			return nil, "", fmt.Errorf("action %s is invalid for retail delivery", action)
		}
	default:
		return nil, "", fmt.Errorf("unsupported shipment type %s", shipment.Type)
	}
}

func (s *Service) publishAssigned(ctx context.Context, shipment domain.Shipment) error {
	return s.publishAssignedWithMeta(ctx, shipment, firstNonEmpty(shipment.OrderID, shipment.HarvestID, shipment.ID), "", "")
}

func (s *Service) publishAssignedWithMeta(ctx context.Context, shipment domain.Shipment, correlationID string, causationID string, traceID string) error {
	if shipment.Type == domain.ShipmentTypeFarmPickup {
		event := events.LogisticsPickupStatusChanged{
			EventID:     uuid.NewString(),
			ShipmentID:  shipment.ID,
			HarvestID:   shipment.HarvestID,
			FarmID:      shipment.FarmID,
			WarehouseID: shipment.WarehouseID,
			DriverID:    shipment.DriverID,
			VehicleID:   shipment.VehicleID,
			Status:      string(domain.ShipmentStatusAssigned),
			OccurredAt:  time.Now(),
		}
		return s.publishLogisticsEvent(ctx, s.pickupAssignedTopic, shipment.ID, event, correlationID, causationID, traceID)
	}
	event := events.LogisticsShipmentAssigned{EventID: uuid.NewString(), ShipmentID: shipment.ID, OrderID: shipment.OrderID, StoreID: shipment.DestinationStoreID, DriverID: shipment.DriverID, OccurredAt: time.Now()}
	return s.publishLogisticsEvent(ctx, s.shipmentAssignedTopic, firstNonEmpty(shipment.OrderID, shipment.ID), event, correlationID, causationID, traceID)
}

func (s *Service) publishStatusChanged(ctx context.Context, topic string, shipment domain.Shipment) error {
	if shipment.Type == domain.ShipmentTypeFarmPickup {
		event := events.LogisticsPickupStatusChanged{
			EventID:     uuid.NewString(),
			ShipmentID:  shipment.ID,
			HarvestID:   shipment.HarvestID,
			FarmID:      shipment.FarmID,
			WarehouseID: shipment.WarehouseID,
			DriverID:    shipment.DriverID,
			VehicleID:   shipment.VehicleID,
			Status:      string(shipment.Status),
			OccurredAt:  time.Now(),
		}
		return s.publishLogisticsEvent(ctx, topic, shipment.ID, event, firstNonEmpty(shipment.HarvestID, shipment.ID), "", "")
	}
	event := events.LogisticsDeliveryStatusChanged{
		EventID:     uuid.NewString(),
		ShipmentID:  shipment.ID,
		OrderID:     shipment.OrderID,
		StoreID:     shipment.DestinationStoreID,
		WarehouseID: shipment.WarehouseID,
		DriverID:    shipment.DriverID,
		VehicleID:   shipment.VehicleID,
		Status:      string(shipment.Status),
		OccurredAt:  time.Now(),
	}
	return s.publishLogisticsEvent(ctx, topic, shipment.ID, event, firstNonEmpty(shipment.OrderID, shipment.ID), "", "")
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
	case events.LogisticsPickupStatusChanged:
		return events.Metadata{EventID: event.EventID, OccurredAt: event.OccurredAt, ShipmentID: event.ShipmentID, HarvestID: event.HarvestID, FarmID: event.FarmID, WarehouseID: event.WarehouseID, DriverID: event.DriverID, VehicleID: event.VehicleID}
	case events.LogisticsDeliveryStatusChanged:
		return events.Metadata{EventID: event.EventID, OccurredAt: event.OccurredAt, ShipmentID: event.ShipmentID, OrderID: event.OrderID, StoreID: event.StoreID, WarehouseID: event.WarehouseID, DriverID: event.DriverID, VehicleID: event.VehicleID}
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
