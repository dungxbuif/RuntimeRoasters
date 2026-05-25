package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"

	"RuntimeRoasters/pkg/events"
	"RuntimeRoasters/pkg/kafka"
	"github.com/google/uuid"
)

type HarvestCreatedEvent struct {
	HarvestID  string  `json:"harvest_id"`
	CoffeeType string  `json:"coffee_type"`
	OriginCode string  `json:"origin_code"`
	Quantity   float64 `json:"quantity"`
}

func main() {
	rand.Seed(time.Now().UnixNano())

	brokers := []string{"localhost:9094"}
	producer := kafka.NewProducer(brokers)
	defer producer.Close()

	topic := events.TopicFarmHarvestCreated
	coffeeTypes := []string{"ARABICA", "ROBUSTA"}
	origins := []string{"CD", "BMT", "GL"}

	fmt.Println("🚀 Starting Harvest Simulation (Demo for RR-19.1)...")

	for i := 1; i <= 5; i++ {
		harvestID := uuid.New().String()
		coffee := coffeeTypes[rand.Intn(len(coffeeTypes))]
		origin := origins[rand.Intn(len(origins))]
		quantity := 100.0 + rand.Float64()*400.0 // 100kg to 500kg

		event := HarvestCreatedEvent{
			HarvestID:  harvestID,
			CoffeeType: coffee,
			OriginCode: origin,
			Quantity:   quantity,
		}
		cloudEvent, err := events.NewCloudEvent(context.Background(), topic, events.SourceFarmService, fmt.Sprintf("harvests/%s", harvestID), event, events.Metadata{
			EventID:       uuid.NewString(),
			CorrelationID: harvestID,
			HarvestID:     harvestID,
			OccurredAt:    time.Now(),
		})
		if err != nil {
			log.Printf("Failed to build CloudEvent: %v", err)
			continue
		}

		fmt.Printf("[%d] Sending Harvest: %s | %s | %.2fkg\n", i, coffee, origin, quantity)

		err = producer.Publish(context.Background(), topic, harvestID, cloudEvent)
		if err != nil {
			log.Printf("Failed to publish: %v", err)
		}

		time.Sleep(2 * time.Second)
	}

	fmt.Println("✅ Simulation finished. Check Warehouse Service logs for processing details.")
}
