package usecase

import (
	"context"

	"github.com/dungxbuif/RuntimeRoasters/apps/logistics-service/internal/domain"
)

type ShipmentRepository interface {
	Create(ctx context.Context, shipment *domain.Shipment) error
	Update(ctx context.Context, shipment *domain.Shipment) error
	GetByID(ctx context.Context, id string) (*domain.Shipment, error)
	List(ctx context.Context) ([]domain.Shipment, error)
}

type DriverRepository interface {
	Create(ctx context.Context, driver *domain.Driver) error
	Update(ctx context.Context, driver *domain.Driver) error
	GetByID(ctx context.Context, id string) (*domain.Driver, error)
	FindNearestAvailable(ctx context.Context, lat, lng float64, radiusKm float64) (*domain.Driver, error)
	SetAvailability(ctx context.Context, driverID string, available bool) error
}

type LocationRepository interface {
	Create(ctx context.Context, location *domain.Location) error
	ListByType(ctx context.Context, locType domain.LocationType) ([]domain.Location, error)
}

type IdempotencyRepository interface {
	IsMessageProcessed(ctx context.Context, msgKey string) (bool, error)
	MarkMessageProcessed(ctx context.Context, msgKey string, topic string) error
}
