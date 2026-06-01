package usecase

import (
	"context"
	"fmt"
	"time"

	"RuntimeRoasters/apps/warehouse-service/internal/domain"
	"RuntimeRoasters/pkg/base/identity"
	"RuntimeRoasters/pkg/events"
	"RuntimeRoasters/pkg/kafka"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DispatchUseCase struct {
	db          *gorm.DB
	producer    kafka.Producer
	pickupTopic string // using same logic for dispatch event
}

func NewDispatchUseCase(db *gorm.DB, producer kafka.Producer, pickupTopic string) *DispatchUseCase {
	return &DispatchUseCase{db: db, producer: producer, pickupTopic: pickupTopic}
}

func (uc *DispatchUseCase) ListDispatchRequests(ctx context.Context) ([]domain.DispatchRequest, error) {
	var requests []domain.DispatchRequest
	query := uc.db.WithContext(ctx).Order("created_at DESC")
	_, warehouseIDs, allWarehouses := identity.WarehouseScopeFromContext(ctx)
	if !allWarehouses {
		if len(warehouseIDs) == 0 {
			return []domain.DispatchRequest{}, nil
		}
		query = query.Where("warehouse_id IN ?", warehouseIDs)
	}
	err := query.Find(&requests).Error
	return requests, err
}

func (uc *DispatchUseCase) DispatchDelivery(ctx context.Context, id string, driverID, vehicleID string) (*domain.DispatchRequest, error) {
	var req domain.DispatchRequest
	now := time.Now()
	err := uc.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("id = ?", id).Take(&req).Error; err != nil {
			return err
		}
		if claims, _, allWarehouses := identity.WarehouseScopeFromContext(ctx); !allWarehouses && !claims.CanAccessWarehouse(req.WarehouseID) {
			return fmt.Errorf("warehouse access denied")
		}
		if req.Status == domain.DispatchRequestStatusCompleted || req.Status == domain.DispatchRequestStatusCancelled {
			return fmt.Errorf("dispatch request %s cannot be dispatched from status %s", id, req.Status)
		}

		if req.Status != domain.DispatchRequestStatusDispatched {
			if err := tx.Model(&req).Updates(map[string]interface{}{
				"status":     domain.DispatchRequestStatusDispatched,
				"updated_at": now,
			}).Error; err != nil {
				return err
			}
			req.Status = domain.DispatchRequestStatusDispatched
		}

		// Publish Event to Logistics
		payload := events.LogisticsShipmentAssigned{
			EventID:    uuid.NewString(),
			OrderID:    req.OrderID,
			StoreID:    req.StoreID,
			DriverID:   driverID,
			VehicleID:  vehicleID,
			OccurredAt: now,
		}
		// Using correlationID = OrderID
		cloudEvent, err := events.NewCloudEvent(ctx, events.TopicLogisticsDeliveryAssigned, events.SourceWarehouseService, fmt.Sprintf("orders/%s", req.OrderID), payload, events.Metadata{
			EventID:       payload.EventID,
			CorrelationID: req.OrderID,
			OrderID:       req.OrderID,
			StoreID:       req.StoreID,
			WarehouseID:   req.WarehouseID,
			DriverID:      driverID,
			VehicleID:     vehicleID,
		})
		if err != nil {
			return err
		}
		return uc.producer.Publish(ctx, events.TopicLogisticsDeliveryAssigned, req.OrderID, cloudEvent)
	})
	if err != nil {
		return nil, err
	}
	return &req, nil
}
