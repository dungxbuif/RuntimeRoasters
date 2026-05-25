package seed

import (
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"

	"RuntimeRoasters/apps/logistics-service/internal/domain"
)

type Data struct {
	Vehicles  []domain.Vehicle  `json:"vehicles"`
	Drivers   []domain.Driver   `json:"drivers"`
	Locations []domain.Location `json:"locations"`
}

//go:embed logistics.json
var logisticsSeed []byte

func Load() (Data, error) {
	var data Data
	if err := json.Unmarshal(logisticsSeed, &data); err != nil {
		return Data{}, fmt.Errorf("load logistics seed: %w", err)
	}
	if len(data.Vehicles) == 0 {
		return Data{}, errors.New("logistics vehicles seed is empty")
	}
	if len(data.Drivers) == 0 {
		return Data{}, errors.New("logistics drivers seed is empty")
	}
	if len(data.Locations) == 0 {
		return Data{}, errors.New("logistics locations seed is empty")
	}
	for _, vehicle := range data.Vehicles {
		if vehicle.ID == "" || vehicle.PlateNumber == "" || vehicle.Status == "" {
			return Data{}, fmt.Errorf("invalid logistics vehicle seed: id=%q plate=%q status=%q", vehicle.ID, vehicle.PlateNumber, vehicle.Status)
		}
	}
	for _, driver := range data.Drivers {
		if driver.ID == "" || driver.UserID == "" || driver.Name == "" || driver.Status == "" {
			return Data{}, fmt.Errorf("invalid logistics driver seed: id=%q user_id=%q name=%q status=%q", driver.ID, driver.UserID, driver.Name, driver.Status)
		}
	}
	for _, location := range data.Locations {
		if location.ID == "" || location.Name == "" || location.Type == "" {
			return Data{}, fmt.Errorf("invalid logistics location seed: id=%q name=%q type=%q", location.ID, location.Name, location.Type)
		}
	}
	return data, nil
}
