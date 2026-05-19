package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type Location struct {
	ID   string
	Lat  float64
	Lng  float64
	Type string
}

type RouteData struct {
	ID          string      `json:"id"`
	Name        string      `json:"name"`
	From        string      `json:"from"`
	To          string      `json:"to"`
	Coordinates [][]float64 `json:"coordinates"` // [lat, lng]
}

func main() {
	locations := []Location{
		{"FARM_CAU_DAT", 11.895, 108.538, "FARM"},
		{"FARM_BMT", 12.666, 108.038, "FARM"},
		{"FARM_PLEIKU", 14.010, 108.040, "FARM"},
		{"WH_SONG_THAN", 10.880, 106.750, "ROASTERY"},
		{"WH_HOA_LAC", 21.010, 105.530, "ROASTERY"},
		{"WH_HOA_KHANH", 16.080, 108.150, "ROASTERY"},
		{"RET_HCM", 10.776, 106.700, "RETAILER"},
		{"RET_HN", 21.028, 105.852, "RETAILER"},
		{"RET_DN", 16.066, 108.216, "RETAILER"},
	}

	pairs := []struct{ from, to string }{
		{"FARM_CAU_DAT", "WH_SONG_THAN"},
		{"FARM_BMT", "WH_SONG_THAN"},
		{"FARM_PLEIKU", "WH_HOA_KHANH"},
		{"WH_SONG_THAN", "RET_HCM"},
		{"WH_HOA_KHANH", "RET_DN"},
		{"WH_HOA_LAC", "RET_HN"},
	}

	var allRoutes []RouteData

	for i, pair := range pairs {
		var start, end Location
		for _, loc := range locations {
			if loc.ID == pair.from {
				start = loc
			}
			if loc.ID == pair.to {
				end = loc
			}
		}

		fmt.Printf("Generating route %d: %s -> %s\n", i+1, start.ID, end.ID)
		coords, err := fetchOSRMRoute(start, end)
		if err != nil {
			fmt.Printf("Error fetching route %s -> %s: %v\n", start.ID, end.ID, err)
			continue
		}

		allRoutes = append(allRoutes, RouteData{
			ID:          fmt.Sprintf("route-%03d", i+1),
			Name:        fmt.Sprintf("%s to %s", start.ID, end.ID),
			From:        start.ID,
			To:          end.ID,
			Coordinates: coords,
		})
	}

	data, _ := json.MarshalIndent(allRoutes, "", "  ")

	// Save to testdata
	os.MkdirAll("src/apps/logistics-service/testdata", 0755)
	os.WriteFile("src/apps/logistics-service/testdata/routes.json", data, 0644)

	// Save to public data for FE
	os.MkdirAll("src/apps/client-app/public/data", 0755)
	os.WriteFile("src/apps/client-app/public/data/routes.json", data, 0644)

	fmt.Println("Routes generated successfully!")
}

type OSRMResponse struct {
	Routes []struct {
		Geometry struct {
			Coordinates [][]float64 `json:"coordinates"` // [lng, lat]
		} `json:"geometry"`
	} `json:"routes"`
}

func fetchOSRMRoute(start, end Location) ([][]float64, error) {
	url := fmt.Sprintf("http://router.project-osrm.org/route/v1/driving/%f,%f;%f,%f?overview=full&geometries=geojson",
		start.Lng, start.Lat, end.Lng, end.Lat)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var osrm OSRMResponse
	if err := json.Unmarshal(body, &osrm); err != nil {
		return nil, err
	}

	if len(osrm.Routes) == 0 {
		return nil, fmt.Errorf("no routes found")
	}

	var coords [][]float64
	for _, c := range osrm.Routes[0].Geometry.Coordinates {
		// OSRM returns [lng, lat], we want [lat, lng] for Leaflet
		coords = append(coords, []float64{c[1], c[0]})
	}

	return coords, nil
}
