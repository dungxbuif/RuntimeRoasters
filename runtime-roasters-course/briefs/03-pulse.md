# Module 3: The Logistics Pulse

### Teaching Arc
- **Metaphor:** The "System's Nervous System". Just like your brain sends electrical pulses to your legs to move, our "Events" (Kafka) send digital pulses to drivers to move physical coffee.
- **Opening hook:** A harvest is sitting at the farm. How does the Warehouse know it's time to send a truck? They listen for a "Pulse".
- **Key insight:** "Event-Driven" means the system reacts to things that *happen*, rather than asking "is anything new?" every second.
- **"Why should I care?":** Understanding events helps you design real-time apps. You'll know when to ask AI for a "Webhook" or a "Topic" instead of a "Polling loop".

### Code Snippets (pre-extracted)

File: src/pkg/events/topics.go (lines 5-15)
```go
const (
	TopicFarmHarvestCreated    = "farm.harvest.created"
	TopicWarehousePickupReady  = "warehouse.pickup.ready"
	TopicLogisticsGpsUpdated   = "logistics.driver.gps_updated"
	TopicRetailSaleCompleted   = "retail.sale.completed"
)
```

File: src/apps/logistics-service/internal/domain/shipment.go (lines 20-30)
```go
type DriverLocationUpdate struct {
	DriverID   string    `json:"driver_id"`
	ShipmentID string    `json:"shipment_id"`
	Latitude   float64   `json:"latitude"`
	Longitude  float64   `json:"longitude"`
	OccurredAt time.Time `json:"occurred_at"`
}
```

### Interactive Elements

- [x] **Message Flow Animation** — actors: [Farm Service, Kafka, Warehouse Service]. Steps: [1. Harvest created event, 2. Flies through Kafka "Pipe", 3. Lands in Warehouse listener, 4. Dashboard updates].
- [x] **Code↔English translation** — `DriverLocationUpdate` and Event Topics.
- [x] **Live Pulse visualization** — A vibrating "Event" circle that changes text based on which topic it represents.
- [x] **Quiz** — 3 questions, style: debugging. "The Warehouse Manager says they don't see the new harvest. Which 'wire' (topic) might be broken?"

### Reference Files to Read

- `references/interactive-elements.md` → "Message Flow Animation", "Code ↔ English Translations"
- `references/content-philosophy.md`
- `references/gotchas.md`

### Connections

- **Previous module:** The Cast of Services (Who they are)
- **Next module:** The Retail Dance (Complex multi-step actions)
- **Tone/style notes:** Use "Electric" or "Pulse" visual metaphors. High energy.
