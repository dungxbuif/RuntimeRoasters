package usecase

import (
	"context"
	"fmt"

	"RuntimeRoasters/apps/farm-service/internal/domain"
	"RuntimeRoasters/pkg/logger"
	"go.uber.org/zap"
)

type SystemUsecase interface {
	GetStatus(ctx context.Context) (bool, map[string]int64, error)
	SeedData(ctx context.Context, usersMap map[string]string) (int64, error)
}

type systemUsecase struct {
	farmRepo FarmRepository
}

func NewSystemUsecase(farmRepo FarmRepository) SystemUsecase {
	return &systemUsecase{
		farmRepo: farmRepo,
	}
}

func (u *systemUsecase) GetStatus(ctx context.Context) (bool, map[string]int64, error) {
	count, err := u.farmRepo.Count(ctx)
	if err != nil {
		return false, nil, err
	}

	counts := map[string]int64{
		"farms": count,
	}

	// Consider seeded if we have at least 6 farms (as per master data)
	return count >= 6, counts, nil
}

func (u *systemUsecase) SeedData(ctx context.Context, usersMap map[string]string) (int64, error) {
	log := logger.GetLogger()
	log.Info("Starting farm-service seeding", zap.Int("users_map_size", len(usersMap)))

	// Check if already seeded
	count, _ := u.farmRepo.Count(ctx)
	if count >= 6 {
		log.Info("System already seeded, skipping", zap.Int64("count", count))
		return 0, nil
	}

	farmsToSeed := []struct {
		Name       string
		Location   domain.Location
		Area       float64
		CoffeeType domain.CoffeeType
		Email      string
	}{
		{"K'Ho Coffee Farm", domain.LocationCauDat, 15.5, domain.CoffeeTypeArabica, "mgr.farm.kho@runtimeroasters.com"},
		{"Cau Dat Arabica", domain.LocationCauDat, 45.0, domain.CoffeeTypeArabica, "mgr.farm.caudat@runtimeroasters.com"},
		{"Son Pacamara Farm", domain.LocationCauDat, 12.0, domain.CoffeeTypeArabica, "mgr.farm.sonpacamara@runtimeroasters.com"},
		{"Aeroco Coffee", domain.LocationBuonMaThuot, 20.0, domain.CoffeeTypeRobusta, "mgr.farm.aeroco@runtimeroasters.com"},
		{"Trung Nguyen Village", domain.LocationBuonMaThuot, 5.0, domain.CoffeeTypeRobusta, "mgr.farm.trungnguyen@runtimeroasters.com"},
		{"Chu Se Estate", domain.LocationPleiku, 30.0, domain.CoffeeTypeRobusta, "mgr.farm.chuse@runtimeroasters.com"},
	}

	var created int64
	for _, f := range farmsToSeed {
		ownerID, ok := usersMap[f.Email]
		if !ok {
			log.Warn("Manager email not found in users_map, skipping farm", zap.String("email", f.Email))
			continue
		}

		farm := &domain.Farm{
			Name:       f.Name,
			Location:   f.Location,
			Area:       f.Area,
			CoffeeType: f.CoffeeType,
			OwnerID:    ownerID,
		}

		if err := u.farmRepo.Create(ctx, farm); err != nil {
			return created, fmt.Errorf("failed to create farm %s: %w", f.Name, err)
		}
		created++
	}

	log.Info("Farm-service seeding completed", zap.Int64("created", created))
	return created, nil
}
