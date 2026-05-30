package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/gocql/gocql"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

// Operational DB struct mappings
type Farm struct {
	ID         int64
	Name       string
	Location   string
	CoffeeType string
	OwnerID    string
}

type DriverInfo struct {
	ID        string
	VehicleID string
}

type TimelineEntry struct {
	Topic       string                 `json:"topic"`
	Status      string                 `json:"status"`
	Title       string                 `json:"title"`
	Description string                 `json:"description"`
	Timestamp   string                 `json:"timestamp"`
	Payload     map[string]interface{} `json:"payload,omitempty"`
}

type TraceDocument struct {
	EntityID   string                 `json:"entity_id"`
	EntityType string                 `json:"entity_type"`
	StoreID    string                 `json:"store_id,omitempty"`
	Status     string                 `json:"status"`
	Order      map[string]interface{} `json:"order,omitempty"`
	Payment    map[string]interface{} `json:"payment,omitempty"`
	Shipment   map[string]interface{} `json:"shipment,omitempty"`
	TraceIDs   []string               `json:"trace_ids,omitempty"`
	Timeline   []TimelineEntry        `json:"timeline"`
	Events     []interface{}          `json:"events"`
	UpdatedAt  string                 `json:"updated_at"`
}

func main() {
	rand.Seed(time.Now().UnixNano())
	log.Println("🚀 Initializing Ultra-Rich Direct-DB Demo Seeder Simulator...")

	// 1. Establish DB Connections
	farmDB := openDB("postgres://user:password@localhost:54321/farm_db?sslmode=disable")
	warehouseDB := openDB("postgres://user:password@localhost:54321/warehouse_db?sslmode=disable")
	retailDB := openDB("postgres://user:password@localhost:54321/retail_db?sslmode=disable")
	paymentDB := openDB("postgres://user:password@localhost:54321/payment_db?sslmode=disable")
	logisticsDB := openDB("postgres://user:password@localhost:54321/logistics_db?sslmode=disable")
	traceDB := openDB("postgres://user:password@localhost:54321/trace_db?sslmode=disable")
	auditDB := openDB("postgres://user:password@localhost:54321/audit_db?sslmode=disable")

	defer farmDB.Close()
	defer warehouseDB.Close()
	defer retailDB.Close()
	defer paymentDB.Close()
	defer logisticsDB.Close()
	defer traceDB.Close()
	defer auditDB.Close()

	// Connect to Cassandra
	log.Println("⚡ Connecting to Cassandra Audit Database...")
	cassandraCluster := gocql.NewCluster("127.0.0.1:9042")
	cassandraCluster.Keyspace = "runtime_roasters_audit"
	cassandraCluster.Consistency = gocql.Quorum
	cassandraCluster.Timeout = 10 * time.Second
	cassandraCluster.ConnectTimeout = 10 * time.Second
	cassandraSession, err := cassandraCluster.CreateSession()
	if err != nil {
		log.Printf("⚠️ Cassandra direct connection failed: %v. Audit logs will only be seeded into Postgres fallback.", err)
	} else {
		defer cassandraSession.Close()
		log.Println("✅ Cassandra Audit Store connected successfully.")
	}

	// 2. Clear out old transactional data (keeping master seeds)
	log.Println("🧹 Clearing existing operational transactional data...")
	clearTables(farmDB, "harvests")
	clearTables(warehouseDB, "roast_runs", "intakes", "production_batches", "pick_up_requests")
	clearTables(retailDB, "orders")
	clearTables(paymentDB, "payments", "webhook_events")
	clearTables(logisticsDB, "shipments")
	clearTables(traceDB, "trace_events", "trace_documents")
	clearTables(auditDB, "audit_logs")

	// Clean Elasticsearch
	log.Println("🧹 Resetting Elasticsearch 'coffee_traceability' Index...")
	resetElasticsearch("http://localhost:9200", "coffee_traceability")

	// 3. Dynamic Variables & Master Data Retrieval
	now := time.Now()
	drivers := getLogisticsDrivers(logisticsDB)

	// Fetch seeded Farms from DB
	rows, err := farmDB.Query("SELECT id, name, location, coffee_type, owner_id FROM farms")
	if err != nil {
		log.Fatalf("Failed to fetch seeded farms: %v", err)
	}
	defer rows.Close()

	var farmList []Farm
	for rows.Next() {
		var f Farm
		if err := rows.Scan(&f.ID, &f.Name, &f.Location, &f.CoffeeType, &f.OwnerID); err != nil {
			log.Fatalf("Scan farm error: %v", err)
		}
		farmList = append(farmList, f)
	}

	if len(farmList) == 0 {
		log.Fatal("❌ No seeded farms found in farm_db. Please run migrations first.")
	}

	// Fetch seeded Warehouses from DB
	wRows, err := warehouseDB.Query("SELECT id FROM warehouses")
	if err != nil {
		log.Fatalf("Failed to fetch seeded warehouses: %v", err)
	}
	defer wRows.Close()

	var warehouses []string
	for wRows.Next() {
		var wID string
		if err := wRows.Scan(&wID); err != nil {
			log.Fatalf("Scan warehouse error: %v", err)
		}
		warehouses = append(warehouses, wID)
	}
	if len(warehouses) == 0 {
		warehouses = []string{"WAREHOUSE-HN-001", "WAREHOUSE-HCM-001", "WAREHOUSE-DN-001"}
	}

	// Fetch seeded Stores from DB
	sRows, err := retailDB.Query("SELECT id FROM stores")
	if err != nil {
		log.Fatalf("Failed to fetch seeded stores: %v", err)
	}
	defer sRows.Close()

	var stores []string
	for sRows.Next() {
		var sID string
		if err := sRows.Scan(&sID); err != nil {
			log.Fatalf("Scan store error: %v", err)
		}
		stores = append(stores, sID)
	}
	if len(stores) == 0 {
		stores = []string{
			"11111111-1111-1111-1111-111111111101",
			"11111111-1111-1111-1111-111111111102",
			"11111111-1111-1111-1111-111111111103",
			"11111111-1111-1111-1111-111111111104",
			"11111111-1111-1111-1111-111111111105",
		}
	}

	log.Println("🌾 [Cycle 1] Seeding 20+ rich simulation harvests, processing, batches and transportation runs...")

	// 4. Generate 22 Harvest Supply Chain Loops
	for i := 0; i < 22; i++ {
		daysAgo := rand.Intn(30)
		t := now.AddDate(0, 0, -daysAgo).Add(-time.Duration(rand.Intn(24)) * time.Hour)

		farm := farmList[rand.Intn(len(farmList))]
		qty := float64(400 + rand.Intn(1200))

		status := "COMPLETED"
		if daysAgo < 2 {
			status = "NEW"
		} else if daysAgo == 3 {
			status = "PROCESSING"
		}

		// Insert Harvest
		var harvestID int64
		err := farmDB.QueryRow(
			"INSERT INTO harvests (farm_id, owner_id, coffee_type, quantity, harvest_date, status, notes, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id",
			farm.ID, farm.OwnerID, farm.CoffeeType, qty, t, status, "Simulated harvest run", t, t,
		).Scan(&harvestID)
		if err != nil {
			log.Fatalf("Insert harvest failed: %v", err)
		}

		harvestIDStr := fmt.Sprintf("%d", harvestID)
		traceID := uuid.NewString()

		// Generate trace & audit logs for harvest creation
		harvestPayload := fmt.Sprintf(`{"harvest_id":"%d","quantity":%f,"coffee_type":"%s","farm_id":"%d"}`, harvestID, qty, farm.CoffeeType, farm.ID)
		insertTrace(traceDB, "farm-service", "farm.harvest.created", harvestIDStr, traceID, "flow.farm.harvest-to-pickup", "service.farm", "edge.farm.kafka", "outbox-event", harvestPayload, t)
		insertAudit(auditDB, cassandraSession, "harvest", "farm.harvest.created", harvestPayload, t)

		if status == "COMPLETED" {
			whID := warehouses[rand.Intn(len(warehouses))]
			pickupID := uuid.NewString()
			shipmentID := uuid.NewString()

			// Insert Pickup Request
			_, err = warehouseDB.Exec(
				"INSERT INTO pick_up_requests (id, harvest_id, farm_id, warehouse_id, quantity, coffee_type, origin_code, status, shipment_id, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)",
				pickupID, harvestIDStr, fmt.Sprintf("%d", farm.ID), whID, qty, farm.CoffeeType, "VN-LD", "RECEIVED", shipmentID, t.Add(time.Hour), t.Add(12*time.Hour),
			)
			if err != nil {
				log.Fatalf("Insert pickup request failed: %v", err)
			}

			// Generate trace for pickup requested
			pickupRequestedPayload := fmt.Sprintf(`{"pickup_id":"%s","harvest_id":"%s","farm_id":"%d","warehouse_id":"%s","quantity":%f}`, pickupID, harvestIDStr, farm.ID, whID, qty)
			insertTrace(traceDB, "warehouse-service", "warehouse.pickup.requested", harvestIDStr, traceID, "flow.farm.harvest-to-pickup", "service.warehouse", "edge.kafka.warehouse", "consumer-projection", pickupRequestedPayload, t.Add(time.Hour))

			// Insert Shipment
			driver := drivers[rand.Intn(len(drivers))]
			_, err = logisticsDB.Exec(
				"INSERT INTO shipments (id, type, harvest_id, driver_id, vehicle_id, origin_location_id, destination_location_id, status, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)",
				shipmentID, "HARVEST_PICKUP", harvestIDStr, driver.ID, driver.VehicleID, fmt.Sprintf("FARM-CAUDAT-%03d", farm.ID), whID, "DELIVERED", t.Add(2*time.Hour), t.Add(8*time.Hour),
			)
			if err != nil {
				log.Fatalf("Insert pickup shipment failed: %v", err)
			}

			// Trace logs for Logistics progress
			logisticsPayload := fmt.Sprintf(`{"shipment_id":"%s","harvest_id":"%s","driver_id":"%s","vehicle_id":"%s","status":"DELIVERED"}`, shipmentID, harvestIDStr, driver.ID, driver.VehicleID)
			insertTrace(traceDB, "logistics-service", "logistics.pickup.completed", harvestIDStr, traceID, "flow.farm.harvest-to-pickup", "service.logistics", "edge.logistics.kafka", "return-event", logisticsPayload, t.Add(8*time.Hour))

			// Insert Intake
			intakeID := uuid.NewString()
			batchUUID := uuid.NewString()
			batchIDStr := fmt.Sprintf("BATCH-%s-%02d", t.Format("060102"), i)
			_, err = warehouseDB.Exec(
				"INSERT INTO intakes (id, harvest_id, pickup_id, warehouse_id, coffee_type, origin_code, quantity, status, batch_id, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)",
				intakeID, harvestIDStr, pickupID, whID, farm.CoffeeType, "VN-LD", qty, "BATCHED", batchUUID, t.Add(13*time.Hour), t.Add(14*time.Hour),
			)
			if err != nil {
				log.Fatalf("Insert intake failed: %v", err)
			}

			// Trace and Audit for intake created
			intakePayload := fmt.Sprintf(`{"intake_id":"%s","harvest_id":"%s","warehouse_id":"%s","quantity":%f}`, intakeID, harvestIDStr, whID, qty)
			insertTrace(traceDB, "warehouse-service", "warehouse.intake.created", harvestIDStr, traceID, "flow.warehouse.intake-processing-stock", "service.warehouse", "edge.warehouse.postgres", "db-transaction", intakePayload, t.Add(13*time.Hour))
			insertAudit(auditDB, cassandraSession, "intake", "warehouse.intake.created", intakePayload, t.Add(13*time.Hour))

			// Insert Production Batch
			loss := 12.0 + rand.Float64()*10.0
			outQty := qty * (100.0 - loss) / 100.0
			_, err = warehouseDB.Exec(
				"INSERT INTO production_batches (id, batch_id, warehouse_id, status, total_input_weight, total_output_weight, weight_loss_percent, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)",
				batchUUID, batchIDStr, whID, "COMPLETED", qty, outQty, loss, t.Add(24*time.Hour), t.Add(28*time.Hour),
			)
			if err != nil {
				log.Fatalf("Insert production batch failed: %v", err)
			}

			// Insert Roast Runs
			_, err = warehouseDB.Exec(
				"INSERT INTO roast_runs (id, batch_id, run_number, input_weight, output_weight, status, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7)",
				uuid.NewString(), batchUUID, 1, qty, outQty, "COMPLETED", t.Add(25*time.Hour),
			)
			if err != nil {
				log.Fatalf("Insert roast run failed: %v", err)
			}

			// Trace for inventory updated
			inventoryPayload := fmt.Sprintf(`{"batch_id":"%s","warehouse_id":"%s","quantity":%f}`, batchIDStr, whID, outQty)
			insertTrace(traceDB, "warehouse-service", "warehouse.inventory.updated", batchIDStr, traceID, "flow.warehouse.intake-processing-stock", "service.warehouse", "edge.warehouse.kafka", "inventory-event", inventoryPayload, t.Add(28*time.Hour))
		}
	}

	log.Println("🛍️ [Cycle 2] Seeding 45+ rich retail orders, payment intents, and logistics runs...")

	// 5. Generate 46 Orders and their complete SAGA journeys
	for i := 0; i < 46; i++ {
		daysAgo := rand.Intn(30)
		t := now.AddDate(0, 0, -daysAgo).Add(-time.Duration(rand.Intn(24)) * time.Hour)

		storeID := stores[rand.Intn(len(stores))]
		qty := float64(5 + rand.Intn(95))
		amount := qty * 250000.0 // VND

		status := "COMPLETED"
		payStatus := "COMPLETED"

		randVal := rand.Intn(100)
		if randVal < 5 {
			status = "CANCELLED"
			payStatus = "FAILED"
		} else if daysAgo == 0 {
			if randVal < 40 {
				status = "PENDING"
				payStatus = "PENDING"
			} else {
				status = "PROCESSING"
				payStatus = "COMPLETED"
			}
		} else if randVal > 95 {
			status = "REFUNDED"
			payStatus = "REFUNDED"
		}

		orderID := uuid.NewString()
		itemJSON := fmt.Sprintf(`[{"sku":"SKU-AR-VN-LD-001","quantity":%f,"price":250000.0}]`, qty)

		// Insert Order
		_, err := retailDB.Exec(
			"INSERT INTO orders (id, store_id, items, total_amount, status, idempotency_key, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)",
			orderID, storeID, itemJSON, amount, status, uuid.NewString(), t, t,
		)
		if err != nil {
			log.Fatalf("Insert order failed: %v", err)
		}

		traceID := uuid.NewString()

		// Trace & Audit for order created
		orderPayload := fmt.Sprintf(`{"order_id":"%s","store_id":"%s","total_amount":%f,"items":[{"sku":"SKU-AR-VN-LD-001","quantity":%f}]}`, orderID, storeID, amount, qty)
		insertTrace(traceDB, "retail-service", "retail.order.created", orderID, traceID, "flow.retail.paid-order-fulfillment", "service.retail", "edge.retail.kafka", "outbox-event", orderPayload, t)
		insertAudit(auditDB, cassandraSession, "order", "retail.order.created", orderPayload, t)

		// Insert Payment
		paymentID := uuid.NewString()
		_, err = paymentDB.Exec(
			"INSERT INTO payments (id, order_id, store_id, provider, provider_ref, amount, status, simulated, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)",
			paymentID, orderID, storeID, "STRIPE", fmt.Sprintf("ch_%s", uuid.NewString()[:12]), amount, payStatus, true, t.Add(2*time.Minute), t.Add(3*time.Minute),
		)
		if err != nil {
			log.Fatalf("Insert payment failed: %v", err)
		}

		// Payment Intent Trace
		intentPayload := fmt.Sprintf(`{"payment_id":"%s","order_id":"%s","amount":%f,"currency":"VND","provider":"STRIPE"}`, paymentID, orderID, amount)
		insertTrace(traceDB, "payment-service", "payment.intent.created", orderID, traceID, "flow.retail.paid-order-fulfillment", "service.payment", "edge.kafka.payment", "consumer-projection", intentPayload, t.Add(1*time.Minute))

		if payStatus == "COMPLETED" {
			// Payment completed Trace & Audit
			completedPayload := fmt.Sprintf(`{"payment_id":"%s","order_id":"%s","amount":%f,"currency":"VND","status":"SUCCESS"}`, paymentID, orderID, amount)
			insertTrace(traceDB, "payment-service", "payment.completed", orderID, traceID, "flow.retail.paid-order-fulfillment", "service.payment", "edge.payment.kafka", "webhook-event", completedPayload, t.Add(3*time.Minute))
			insertAudit(auditDB, cassandraSession, "payment", "payment.completed", completedPayload, t.Add(3*time.Minute))

			// Warehouse stock reservation Trace
			reservationPayload := fmt.Sprintf(`{"order_id":"%s","store_id":"%s","sku":"SKU-AR-VN-LD-001","quantity":%f}`, orderID, storeID, qty)
			insertTrace(traceDB, "warehouse-service", "warehouse.stock.reserved", orderID, traceID, "flow.retail.paid-order-fulfillment", "service.warehouse", "edge.kafka.warehouse", "reservation-event", reservationPayload, t.Add(5*time.Minute))

			// Shipment Seeding
			if status == "COMPLETED" || status == "PROCESSING" {
				shipmentID := uuid.NewString()
				shipStatus := "DELIVERED"
				if status == "PROCESSING" {
					shipStatus = "DEPARTED"
				}

				driver := drivers[rand.Intn(len(drivers))]
				whID := warehouses[rand.Intn(len(warehouses))]

				_, err = logisticsDB.Exec(
					"INSERT INTO shipments (id, type, order_id, driver_id, vehicle_id, origin_location_id, destination_location_id, status, destination_address, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)",
					shipmentID, "RETAIL_DELIVERY", orderID, driver.ID, driver.VehicleID, whID, storeID, shipStatus, "Simulated Retail Delivery Address", t.Add(15*time.Minute), t.Add(2*time.Hour),
				)
				if err != nil {
					log.Fatalf("Insert delivery shipment failed: %v", err)
				}

				// Shipment assigned trace
				assignPayload := fmt.Sprintf(`{"shipment_id":"%s","order_id":"%s","driver_id":"%s","vehicle_id":"%s","status":"ASSIGNED"}`, shipmentID, orderID, driver.ID, driver.VehicleID)
				insertTrace(traceDB, "logistics-service", "logistics.delivery.assigned", orderID, traceID, "flow.logistics.delivery-return-to-base", "service.logistics", "edge.kafka.logistics", "worker-assignment", assignPayload, t.Add(15*time.Minute))

				if shipStatus == "DELIVERED" {
					// Delivered trace & audit
					deliveredPayload := fmt.Sprintf(`{"shipment_id":"%s","order_id":"%s","status":"DELIVERED"}`, shipmentID, orderID)
					insertTrace(traceDB, "logistics-service", "logistics.delivery.completed", orderID, traceID, "flow.logistics.delivery-return-to-base", "service.logistics", "edge.logistics.kafka", "driver-confirmation", deliveredPayload, t.Add(2*time.Hour))
					insertAudit(auditDB, cassandraSession, "shipment", "logistics.delivery.completed", deliveredPayload, t.Add(2*time.Hour))
				}
			}
		} else if payStatus == "FAILED" {
			failedPayload := fmt.Sprintf(`{"payment_id":"%s","order_id":"%s","reason":"Stripe card declined"}`, paymentID, orderID)
			insertTrace(traceDB, "payment-service", "payment.failed", orderID, traceID, "flow.retail.paid-order-fulfillment", "service.payment", "edge.payment.kafka", "webhook-event", failedPayload, t.Add(3*time.Minute))
		}
	}

	log.Println("📊 [Cycle 3] Generating & Aggregating rich Trace Documents (Relational ERD) for fast UI retrieval...")
	generateAndSyncTraceDocuments(traceDB)

	log.Println("✨ Dynamic Rich Demo Data Seeding Completed Successfully! 100% Consistent across PG, Cassandra and ES.")
}

func openDB(url string) *sql.DB {
	db, err := sql.Open("postgres", url)
	if err != nil {
		log.Fatalf("Failed to connect to database %s: %v", url, err)
	}
	if err := db.Ping(); err != nil {
		log.Fatalf("Ping failed to database %s: %v", url, err)
	}
	return db
}

func clearTables(db *sql.DB, tables ...string) {
	for _, t := range tables {
		_, err := db.Exec(fmt.Sprintf("TRUNCATE TABLE %s CASCADE", t))
		if err != nil {
			log.Fatalf("Truncate table %s failed: %v", t, err)
		}
	}
}

func getLogisticsDrivers(db *sql.DB) []DriverInfo {
	rows, err := db.Query("SELECT id, vehicle_id FROM drivers WHERE vehicle_id IS NOT NULL")
	if err != nil {
		log.Fatalf("Failed to fetch drivers: %v", err)
	}
	defer rows.Close()

	var list []DriverInfo
	for rows.Next() {
		var d DriverInfo
		if err := rows.Scan(&d.ID, &d.VehicleID); err != nil {
			log.Fatalf("Scan driver failed: %v", err)
		}
		list = append(list, d)
	}

	if len(list) == 0 {
		list = []DriverInfo{
			{"22222222-2222-2222-2222-222222222201", "VEHICLE-DEMO-001"},
			{"22222222-2222-2222-2222-222222222202", "VEHICLE-DEMO-002"},
			{"22222222-2222-2222-2222-222222222203", "VEHICLE-DEMO-003"},
		}
	}
	return list
}

func insertTrace(db *sql.DB, service, topic, entityID, traceID, flowID, nodeID, edgeID, pattern, payload string, t time.Time) {
	// Parse fields depending on type
	harvestID, batchID, orderID := "", "", ""
	if strings.Contains(topic, "harvest") || strings.Contains(topic, "pickup") {
		harvestID = entityID
	} else if strings.Contains(topic, "inventory") || strings.Contains(topic, "batch") {
		batchID = entityID
	} else {
		orderID = entityID
	}

	_, err := db.Exec(
		`INSERT INTO trace_events (
			id, message_id, topic, batch_id, order_id, shipment_id, store_id, harvest_id, farm_id, warehouse_id, 
			driver_id, vehicle_id, trace_id, flow_id, node_id, edge_id, pattern, source_service, visibility, 
			display_payload, payload, occurred_at, created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23)`,
		uuid.NewString(), uuid.NewString(), topic, batchID, orderID, "", "", harvestID, "", "",
		"", "", traceID, flowID, nodeID, edgeID, pattern, service, "public",
		"{}", payload, t, t,
	)
	if err != nil {
		log.Printf("⚠️ Insert trace failed: %v", err)
	}
}

func insertAudit(db *sql.DB, cassandraSession *gocql.Session, partitionKey, topic, payload string, t time.Time) {
	messageID := uuid.NewString()
	id := uuid.NewString()
	prevHash := "0000000000000000000000000000000000000000000000000000000000000000"
	currHash := uuid.NewString() // Simulated cryptographic SHA-256 equivalent hash

	// 1. Insert into Postgres fallback
	_, err := db.Exec(
		"INSERT INTO audit_logs (id, partition_key, message_id, topic, store_id, payload, previous_hash, current_hash, occurred_at, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)",
		id, partitionKey, messageID, topic, "", payload, prevHash, currHash, t, t,
	)
	if err != nil {
		log.Printf("⚠️ Insert audit log fallback failed: %v", err)
	}

	// 2. Insert into Cassandra
	if cassandraSession != nil {
		err = cassandraSession.Query(`
			INSERT INTO audit_logs (
				partition_key, occurred_at, id, message_id, topic, store_id, payload, previous_hash, current_hash, created_at
			) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			partitionKey, t, id, messageID, topic, "", payload, prevHash, currHash, t,
		).Exec()
		if err != nil {
			log.Printf("⚠️ Cassandra audit log insert failed: %v", err)
		}
	}
}

func resetElasticsearch(baseURL, index string) {
	client := &http.Client{Timeout: 5 * time.Second}
	req, _ := http.NewRequest(http.MethodDelete, baseURL+"/"+index, nil)
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("⚠️ Failed to delete ES index: %v", err)
		return
	}
	resp.Body.Close()

	// Recreate mapping
	mapping := `{
		"mappings": {
			"properties": {
				"entity_id": {"type": "keyword"},
				"entity_type": {"type": "keyword"},
				"status": {"type": "keyword"},
				"trace_ids": {"type": "keyword"},
				"updated_at": {"type": "date"},
				"order": {"type": "object"},
				"payment": {"type": "object"},
				"shipment": {"type": "object"},
				"timeline": {
					"type": "nested",
					"properties": {
						"topic": {"type": "keyword"},
						"status": {"type": "keyword"},
						"title": {"type": "text"},
						"description": {"type": "text"},
						"timestamp": {"type": "date"},
						"payload": {"type": "object", "enabled": false}
					}
				},
				"events": {
					"type": "nested"
				}
			}
		}
	}`
	req, _ = http.NewRequest(http.MethodPut, baseURL+"/"+index, bytes.NewReader([]byte(mapping)))
	req.Header.Set("Content-Type", "application/json")
	resp, err = client.Do(req)
	if err != nil {
		log.Printf("⚠️ Failed to create ES mapping: %v", err)
		return
	}
	resp.Body.Close()
	log.Println("✅ Elasticsearch mapping successfully created.")
}

func generateAndSyncTraceDocuments(db *sql.DB) {
	// Rebuild TraceDocuments in GORM/Postgres and push them synchronously to Elasticsearch
	rows, err := db.Query("SELECT id, message_id, topic, batch_id, order_id, shipment_id, store_id, harvest_id, trace_id, occurred_at, payload FROM trace_events ORDER BY occurred_at ASC")
	if err != nil {
		log.Printf("⚠️ Fetch trace events failed: %v", err)
		return
	}
	defer rows.Close()

	type TraceEvent struct {
		ID            string    `json:"id"`
		MessageID     string    `json:"message_id"`
		Topic         string    `json:"topic"`
		BatchID       string    `json:"batch_id"`
		OrderID       string    `json:"order_id"`
		ShipmentID    string    `json:"shipment_id"`
		StoreID       string    `json:"store_id"`
		HarvestID     string    `json:"harvest_id"`
		TraceID       string    `json:"trace_id"`
		OccurredAt    time.Time `json:"occurred_at"`
		PayloadString string    `json:"-"`
	}

	entityEventsMap := make(map[string][]TraceEvent)
	for rows.Next() {
		var e TraceEvent
		if err := rows.Scan(&e.ID, &e.MessageID, &e.Topic, &e.BatchID, &e.OrderID, &e.ShipmentID, &e.StoreID, &e.HarvestID, &e.TraceID, &e.OccurredAt, &e.PayloadString); err != nil {
			continue
		}

		// Group by EntityID
		entityID := ""
		if e.OrderID != "" {
			entityID = e.OrderID
		} else if e.HarvestID != "" {
			entityID = e.HarvestID
		} else if e.BatchID != "" {
			entityID = e.BatchID
		}

		if entityID != "" {
			entityEventsMap[entityID] = append(entityEventsMap[entityID], e)
		}
	}

	httpClient := &http.Client{Timeout: 5 * time.Second}

	for entityID, eventsList := range entityEventsMap {
		entityType := "ORDER"
		if len(eventsList) > 0 {
			if eventsList[0].HarvestID != "" {
				entityType = "HARVEST"
			} else if eventsList[0].BatchID != "" {
				entityType = "BATCH"
			}
		}

		// Build Read Model
		var timeline []TimelineEntry
		var traceIDs []string
		traceIDsMap := make(map[string]bool)

		status := "COMPLETED"
		var orderPayload, paymentPayload, shipmentPayload map[string]interface{}

		for _, ev := range eventsList {
			if ev.TraceID != "" && !traceIDsMap[ev.TraceID] {
				traceIDsMap[ev.TraceID] = true
				traceIDs = append(traceIDs, ev.TraceID)
			}

			var payload map[string]interface{}
			_ = json.Unmarshal([]byte(ev.PayloadString), &payload)

			title := ev.Topic
			desc := fmt.Sprintf("Simulated event of type %s", ev.Topic)

			// Fill real dynamic descriptions
			switch ev.Topic {
			case "farm.harvest.created":
				title = "Harvest Declared"
				desc = fmt.Sprintf("Simulated harvest created for %.2f kg", payload["quantity"])
			case "warehouse.pickup.requested":
				title = "Pickup Requested"
				desc = "Pickup requested for harvest transport"
			case "logistics.pickup.completed":
				title = "Pickup Delivered"
				desc = "Harvest safely arrived at Roastery"
			case "warehouse.intake.created":
				title = "Warehouse Inbound Intake"
				desc = fmt.Sprintf("Intake completed for %.2f kg", payload["quantity"])
			case "warehouse.inventory.updated":
				title = "Roasting & Inventory Finished"
				desc = fmt.Sprintf("Finished processing roasted batch %s", ev.BatchID)
			case "retail.order.created":
				title = "Order Created"
				desc = fmt.Sprintf("Order placed successfully for %.2f VND", payload["total_amount"])
				orderPayload = payload
			case "payment.intent.created":
				title = "Checkout Created"
				desc = "Stripe simulated checkout transaction initialized"
			case "payment.completed":
				title = "Payment Successful"
				desc = "Payment successfully confirmed by provider webhooks"
				paymentPayload = payload
			case "payment.failed":
				title = "Payment Declined"
				desc = "Payment checkout was declined by the provider"
				status = "FAILED"
			case "warehouse.stock.reserved":
				title = "Warehouse Stock Allocated"
				desc = fmt.Sprintf("Successfully reserved %v units of SKU", payload["quantity"])
			case "logistics.delivery.assigned":
				title = "Logistics Fleet Assigned"
				desc = "Assigned driver and vehicle logistics courier"
			case "logistics.delivery.completed":
				title = "Delivered to Store"
				desc = "Retail shipment safely arrived at store front"
				shipmentPayload = payload
			}

			timeline = append(timeline, TimelineEntry{
				Topic:       ev.Topic,
				Status:      status,
				Title:       title,
				Description: desc,
				Timestamp:   ev.OccurredAt.Format(time.RFC3339),
				Payload:     payload,
			})
		}

		traceIDsJSON, _ := json.Marshal(traceIDs)

		doc := TraceDocument{
			EntityID:   entityID,
			EntityType: entityType,
			Status:     status,
			Order:      orderPayload,
			Payment:    paymentPayload,
			Shipment:   shipmentPayload,
			TraceIDs:   traceIDs,
			Timeline:   timeline,
			UpdatedAt:  time.Now().Format(time.RFC3339),
		}

		docBytes, _ := json.Marshal(doc)

		// 1. Insert into trace_documents table in PostgreSQL
		_, err = db.Exec(
			"INSERT INTO trace_documents (id, entity_id, entity_type, store_id, trace_ids, document, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7) ON CONFLICT (entity_id) DO UPDATE SET document = EXCLUDED.document, trace_ids = EXCLUDED.trace_ids, updated_at = EXCLUDED.updated_at",
			uuid.NewString(), entityID, entityType, "", string(traceIDsJSON), string(docBytes), time.Now(),
		)
		if err != nil {
			log.Printf("⚠️ SQL Trace document write failed for %s: %v", entityID, err)
		}

		// 2. Synchronize directly to Elasticsearch
		esURL := fmt.Sprintf("http://localhost:9200/coffee_traceability/_doc/%s?refresh=true", entityID)
		req, _ := http.NewRequest(http.MethodPut, esURL, bytes.NewReader(docBytes))
		req.Header.Set("Content-Type", "application/json")
		resp, err := httpClient.Do(req)
		if err != nil {
			log.Printf("⚠️ ES write failed for %s: %v", entityID, err)
		} else {
			resp.Body.Close()
		}
	}

	log.Println("✅ All Trace Documents successfully synchronized to ES.")
}

