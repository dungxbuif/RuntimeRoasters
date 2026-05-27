package topology

import (
	"encoding/json"
	"strings"
	"time"

	"RuntimeRoasters/pkg/events"
	cloudevents "github.com/cloudevents/sdk-go/v2"
)

const (
	VisibilityPublic  = "public"
	VisibilityPrivate = "private"
)

type Node struct {
	ID       string `json:"id"`
	Label    string `json:"label"`
	Subtitle string `json:"subtitle"`
	Icon     string `json:"icon"`
	Tone     string `json:"tone"`
}

type Edge struct {
	ID      string `json:"id"`
	Source  string `json:"source"`
	Target  string `json:"target"`
	Label   string `json:"label,omitempty"`
	Pattern string `json:"pattern,omitempty"`
}

type Flow struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	NodeIDs     []string `json:"node_ids"`
	EdgeIDs     []string `json:"edge_ids"`
	LastSeenAt  *string  `json:"last_seen_at,omitempty"`
}

type Config struct {
	Nodes []Node `json:"nodes"`
	Edges []Edge `json:"edges"`
	Flows []Flow `json:"flows"`
}

type Projection struct {
	FlowID         string                 `json:"flow_id"`
	NodeID         string                 `json:"node_id"`
	EdgeID         string                 `json:"edge_id"`
	Pattern        string                 `json:"pattern"`
	SourceService  string                 `json:"source_service"`
	Visibility     string                 `json:"visibility"`
	TraceID        string                 `json:"trace_id,omitempty"`
	DisplayPayload map[string]interface{} `json:"display_payload"`
}

type HistoryEntry struct {
	EventID        string                 `json:"event_id"`
	Topic          string                 `json:"topic"`
	Status         string                 `json:"status"`
	FlowID         string                 `json:"flow_id"`
	NodeID         string                 `json:"node_id"`
	EdgeID         string                 `json:"edge_id"`
	Pattern        string                 `json:"pattern"`
	SourceService  string                 `json:"source_service"`
	Visibility     string                 `json:"visibility"`
	TraceID        string                 `json:"trace_id,omitempty"`
	OccurredAt     time.Time              `json:"occurred_at"`
	DisplayPayload map[string]interface{} `json:"display_payload,omitempty"`
}

type BroadcastRequested struct {
	EventID       string                 `json:"event_id"`
	Channel       string                 `json:"channel"`
	FlowID        string                 `json:"flow_id"`
	NodeID        string                 `json:"node_id"`
	EdgeID        string                 `json:"edge_id"`
	Visibility    string                 `json:"visibility"`
	Status        string                 `json:"status"`
	SourceService string                 `json:"source_service"`
	Payload       map[string]interface{} `json:"payload"`
	OccurredAt    time.Time              `json:"occurred_at"`
}

type topicDef struct {
	FlowID     string
	NodeID     string
	EdgeID     string
	Pattern    string
	Visibility string
}

var topicMap = map[string]topicDef{
	events.TopicFarmHarvestCreated:                {"flow.farm.harvest-to-pickup", "service.farm", "edge.farm.kafka", "outbox-event", VisibilityPublic},
	events.TopicWarehousePickupRequested:          {"flow.farm.harvest-to-pickup", "service.warehouse", "edge.kafka.warehouse", "consumer-projection", VisibilityPublic},
	events.TopicLogisticsPickupAssigned:           {"flow.farm.harvest-to-pickup", "service.logistics", "edge.kafka.logistics", "worker-assignment", VisibilityPublic},
	events.TopicLogisticsPickupDeparted:           {"flow.farm.harvest-to-pickup", "service.logistics", "edge.logistics.valkey", "driver-state", VisibilityPublic},
	events.TopicLogisticsPickupArrivedAtFarm:      {"flow.farm.harvest-to-pickup", "service.logistics", "edge.logistics.valkey", "driver-state", VisibilityPublic},
	events.TopicLogisticsPickupLoadingConfirmed:   {"flow.farm.harvest-to-pickup", "service.logistics", "edge.logistics.kafka", "status-event", VisibilityPublic},
	events.TopicLogisticsPickupReturnStarted:      {"flow.farm.harvest-to-pickup", "service.logistics", "edge.logistics.valkey", "driver-state", VisibilityPublic},
	events.TopicLogisticsPickupArrivedAtWarehouse: {"flow.farm.harvest-to-pickup", "service.logistics", "edge.logistics.kafka", "return-event", VisibilityPublic},
	events.TopicWarehousePickupReceived:           {"flow.farm.harvest-to-pickup", "service.warehouse", "edge.kafka.warehouse", "inventory-intake", VisibilityPublic},
	events.TopicWarehouseIntakeCreated:            {"flow.warehouse.intake-processing-stock", "service.warehouse", "edge.warehouse.postgres", "db-transaction", VisibilityPublic},
	events.TopicWarehouseInventoryUpdated:         {"flow.warehouse.intake-processing-stock", "service.warehouse", "edge.warehouse.kafka", "inventory-event", VisibilityPublic},
	events.TopicRetailOrderCreated:                {"flow.retail.paid-order-fulfillment", "service.retail", "edge.retail.kafka", "outbox-event", VisibilityPublic},
	events.TopicPaymentIntentCreated:              {"flow.retail.paid-order-fulfillment", "service.payment", "edge.kafka.payment", "consumer-projection", VisibilityPublic},
	events.TopicPaymentCompleted:                  {"flow.retail.paid-order-fulfillment", "service.payment", "edge.payment.kafka", "webhook-event", VisibilityPublic},
	events.TopicPaymentSimulatedCompleted:         {"flow.retail.paid-order-fulfillment", "service.payment", "edge.payment.kafka", "legacy-simulated-event", VisibilityPublic},
	events.TopicPaymentFailed:                     {"flow.retail.paid-order-fulfillment", "service.payment", "edge.payment.kafka", "webhook-event", VisibilityPublic},
	events.TopicWarehouseStockReserved:            {"flow.retail.paid-order-fulfillment", "service.warehouse", "edge.kafka.warehouse", "reservation-event", VisibilityPublic},
	events.TopicWarehouseStockReservationFailed:   {"flow.retail.paid-order-fulfillment", "service.warehouse", "edge.warehouse.kafka", "reservation-event", VisibilityPublic},
	events.TopicWarehouseDispatchRequested:        {"flow.retail.paid-order-fulfillment", "service.warehouse", "edge.warehouse.kafka", "dispatch-event", VisibilityPublic},
	events.TopicLogisticsDeliveryAssigned:         {"flow.logistics.delivery-return-to-base", "service.logistics", "edge.kafka.logistics", "worker-assignment", VisibilityPublic},
	events.TopicLogisticsDeliveryDeparted:         {"flow.logistics.delivery-return-to-base", "service.logistics", "edge.logistics.valkey", "driver-state", VisibilityPublic},
	events.TopicLogisticsDeliveryArrivedAtStore:   {"flow.logistics.delivery-return-to-base", "service.logistics", "edge.logistics.valkey", "driver-state", VisibilityPublic},
	events.TopicLogisticsDeliveryDriverConfirmed:  {"flow.logistics.delivery-return-to-base", "service.logistics", "edge.logistics.kafka", "driver-confirmation", VisibilityPublic},
	events.TopicLogisticsDriverReturnStarted:      {"flow.logistics.delivery-return-to-base", "service.logistics", "edge.logistics.valkey", "return-state", VisibilityPublic},
	events.TopicLogisticsDriverReturnedToBase:     {"flow.logistics.delivery-return-to-base", "service.logistics", "edge.logistics.kafka", "return-event", VisibilityPublic},
	events.TopicLogisticsGPSUpdated:               {"flow.logistics.delivery-return-to-base", "service.logistics", "edge.logistics.valkey", "gps-coalesced", VisibilityPublic},
	events.TopicSocketBroadcastRequested:          {"flow.realtime.socket-push-pull", "service.socket", "edge.socket.kafka", "websocket-broadcast", VisibilityPublic},
}

func CanonicalConfig() Config {
	return Config{
		Nodes: []Node{
			{"client.web", "Client App", "Next.js control plane", "devices", "touchpoint"},
			{"gateway.krakend", "KrakenD", "REST gateway", "router", "gateway"},
			{"identity.kratos", "Ory Kratos", "Identity", "fingerprint", "identity"},
			{"identity.hydra", "Ory Hydra", "OIDC", "lock", "identity"},
			{"service.auth", "Auth Service", "Casbin policy owner", "passkey", "service"},
			{"service.farm", "Farm Service", "Harvest operations", "agriculture", "service"},
			{"service.retail", "Retail Service", "Paid orders", "storefront", "service"},
			{"service.payment", "Payment Service", "Stripe webhook demo", "payments", "service"},
			{"service.warehouse", "Warehouse Service", "Inventory & reservation", "warehouse", "service"},
			{"service.logistics", "Logistics Service", "Driver state", "local_shipping", "service"},
			{"service.trace", "Trace Service", "History/read models", "timeline", "service"},
			{"service.audit", "Audit Service", "Immutable logs", "policy", "service"},
			{"service.socket", "Socket Service", "WebSocket fanout", "sensors", "service"},
			{"infra.kafka", "Kafka", "Event backbone", "hub", "infra"},
			{"infra.postgres", "PostgreSQL", "Service data", "database", "db"},
			{"infra.elasticsearch", "Elasticsearch", "Trace search", "manage_search", "db"},
			{"infra.cassandra", "Cassandra", "Audit log", "storage", "db"},
			{"infra.valkey", "Valkey", "Realtime session/GPS", "bolt", "infra"},
			{"infra.otel", "OpenTelemetry", "Trace correlation", "monitoring", "infra"},
		},
		Edges: []Edge{
			{"edge.client.gateway", "client.web", "gateway.krakend", "REST", "gateway"},
			{"edge.gateway.auth", "gateway.krakend", "service.auth", "HTTP/gRPC", "authz"},
			{"edge.gateway.farm", "gateway.krakend", "service.farm", "HTTP/gRPC", "api"},
			{"edge.gateway.retail", "gateway.krakend", "service.retail", "HTTP/gRPC", "api"},
			{"edge.gateway.payment", "gateway.krakend", "service.payment", "HTTP", "webhook"},
			{"edge.gateway.trace", "gateway.krakend", "service.trace", "HTTP", "query"},
			{"edge.client.socket", "client.web", "service.socket", "WebSocket", "realtime"},
			{"edge.auth.postgres", "service.auth", "infra.postgres", "SQL", "db"},
			{"edge.farm.postgres", "service.farm", "infra.postgres", "SQL", "db"},
			{"edge.retail.postgres", "service.retail", "infra.postgres", "SQL", "db"},
			{"edge.payment.postgres", "service.payment", "infra.postgres", "SQL", "db"},
			{"edge.warehouse.postgres", "service.warehouse", "infra.postgres", "SQL", "db"},
			{"edge.trace.postgres", "service.trace", "infra.postgres", "SQL", "read-model"},
			{"edge.trace.elasticsearch", "service.trace", "infra.elasticsearch", "Index", "search"},
			{"edge.audit.cassandra", "service.audit", "infra.cassandra", "Write", "audit"},
			{"edge.logistics.valkey", "service.logistics", "infra.valkey", "GEO/TTL", "realtime-state"},
			{"edge.socket.valkey", "service.socket", "infra.valkey", "Session TTL", "sticky-session"},
			{"edge.farm.kafka", "service.farm", "infra.kafka", "Outbox", "event"},
			{"edge.retail.kafka", "service.retail", "infra.kafka", "Outbox", "event"},
			{"edge.payment.kafka", "service.payment", "infra.kafka", "Payment event", "event"},
			{"edge.warehouse.kafka", "service.warehouse", "infra.kafka", "Inventory event", "event"},
			{"edge.logistics.kafka", "service.logistics", "infra.kafka", "Shipment event", "event"},
			{"edge.socket.kafka", "service.socket", "infra.kafka", "Broadcast event", "event"},
			{"edge.kafka.payment", "infra.kafka", "service.payment", "Consume", "consumer"},
			{"edge.kafka.warehouse", "infra.kafka", "service.warehouse", "Consume", "consumer"},
			{"edge.kafka.logistics", "infra.kafka", "service.logistics", "Consume", "consumer"},
			{"edge.kafka.trace", "infra.kafka", "service.trace", "Consume", "projection"},
			{"edge.kafka.socket", "infra.kafka", "service.socket", "Consume", "fanout"},
			{"edge.services.otel", "service.trace", "infra.otel", "Trace IDs", "correlation"},
		},
		Flows: []Flow{
			{"flow.auth.oidc-login", "OIDC Login", "Browser login through KrakenD, Hydra, Kratos, Auth, and policy data.", []string{"client.web", "gateway.krakend", "identity.hydra", "identity.kratos", "service.auth", "infra.postgres"}, []string{"edge.client.gateway", "edge.gateway.auth", "edge.auth.postgres"}, nil},
			{"flow.farm.harvest-to-pickup", "Harvest To Pickup", "Farm harvest events create warehouse pickup and logistics return flow.", []string{"service.farm", "infra.kafka", "service.warehouse", "service.logistics", "infra.valkey"}, []string{"edge.farm.kafka", "edge.kafka.warehouse", "edge.kafka.logistics", "edge.logistics.valkey"}, nil},
			{"flow.warehouse.intake-processing-stock", "Warehouse Processing", "Warehouse intake, processing, and stock update projection.", []string{"service.warehouse", "infra.postgres", "infra.kafka", "service.trace"}, []string{"edge.warehouse.postgres", "edge.warehouse.kafka", "edge.kafka.trace"}, nil},
			{"flow.retail.paid-order-fulfillment", "Paid Order Fulfillment", "Retail order waits for Stripe webhook, then warehouse reserves stock.", []string{"service.retail", "infra.kafka", "service.payment", "service.warehouse", "service.logistics"}, []string{"edge.retail.kafka", "edge.kafka.payment", "edge.payment.kafka", "edge.kafka.warehouse", "edge.kafka.logistics"}, nil},
			{"flow.logistics.delivery-return-to-base", "Delivery Return To Base", "Driver delivery confirmation does not complete the order until returned to base.", []string{"service.logistics", "infra.valkey", "infra.kafka", "service.retail"}, []string{"edge.logistics.valkey", "edge.logistics.kafka", "edge.kafka.trace"}, nil},
			{"flow.trace.public-qr", "Public QR Trace", "Public sanitized query of trace documents and history.", []string{"client.web", "gateway.krakend", "service.trace", "infra.elasticsearch", "infra.postgres"}, []string{"edge.client.gateway", "edge.gateway.trace", "edge.trace.elasticsearch", "edge.trace.postgres"}, nil},
			{"flow.realtime.socket-push-pull", "Socket Push/Pull", "Socket-service streams live events while trace-service owns pull history.", []string{"client.web", "service.socket", "infra.valkey", "infra.kafka", "service.trace"}, []string{"edge.client.socket", "edge.socket.valkey", "edge.kafka.socket", "edge.kafka.trace"}, nil},
			{"flow.authz.policy-sync", "AuthZ Policy Sync", "Auth owns policy writes; services refresh local Casbin readers.", []string{"service.auth", "infra.postgres", "infra.kafka", "gateway.krakend"}, []string{"edge.auth.postgres", "edge.kafka.trace"}, nil},
		},
	}
}

func ProjectEvent(event cloudevents.Event, rawPayload []byte) Projection {
	def, ok := topicMap[event.Type()]
	if !ok {
		def = topicDef{FlowID: "flow.trace.public-qr", NodeID: "service.trace", EdgeID: "edge.kafka.trace", Pattern: "observed-event", Visibility: VisibilityPublic}
	}
	ids := eventIDs(event)
	payload := sanitizePayload(event.Type(), ids)
	if status, ok := eventStatus(event.Type()); ok {
		payload["status"] = status
	}
	if len(rawPayload) > 0 {
		var envelope map[string]interface{}
		if err := json.Unmarshal(rawPayload, &envelope); err == nil {
			if data, ok := envelope["data"].(map[string]interface{}); ok {
				for _, key := range []string{"status", "reason"} {
					if value, ok := data[key]; ok {
						payload[key] = value
					}
				}
			}
		}
	}
	return Projection{
		FlowID:         def.FlowID,
		NodeID:         def.NodeID,
		EdgeID:         def.EdgeID,
		Pattern:        def.Pattern,
		SourceService:  sourceService(event.Source()),
		Visibility:     def.Visibility,
		TraceID:        ids["trace_id"],
		DisplayPayload: payload,
	}
}

func eventIDs(event cloudevents.Event) map[string]string {
	ids := map[string]string{}
	for ext, key := range map[string]string{
		"traceid": "trace_id", "orderid": "order_id", "paymentid": "payment_id",
		"shipmentid": "shipment_id", "storeid": "store_id", "harvestid": "harvest_id",
		"farmid": "farm_id", "warehouseid": "warehouse_id", "driverid": "driver_id",
		"vehicleid": "vehicle_id", "batchid": "batch_id",
	} {
		if value := extensionString(event, ext); value != "" {
			ids[key] = value
		}
	}
	return ids
}

func sanitizePayload(topic string, ids map[string]string) map[string]interface{} {
	payload := map[string]interface{}{"topic": topic}
	for _, key := range []string{"order_id", "payment_id", "shipment_id", "store_id", "harvest_id", "farm_id", "warehouse_id", "driver_id", "vehicle_id", "batch_id"} {
		if value := ids[key]; value != "" {
			payload[key] = shortID(value)
		}
	}
	return payload
}

func eventStatus(topic string) (string, bool) {
	switch topic {
	case events.TopicPaymentFailed, events.TopicWarehouseStockReservationFailed:
		return "failed", true
	case events.TopicPaymentCompleted, events.TopicPaymentSimulatedCompleted, events.TopicLogisticsDriverReturnedToBase:
		return "completed", true
	case events.TopicLogisticsGPSUpdated:
		return "moving", true
	default:
		return "observed", true
	}
}

func extensionString(event cloudevents.Event, key string) string {
	value, ok := event.Extensions()[key]
	if !ok || value == nil {
		return ""
	}
	switch typed := value.(type) {
	case string:
		return typed
	default:
		return strings.TrimSpace(strings.Trim(strings.ReplaceAll(strings.ReplaceAll(strings.TrimSpace(jsonString(typed)), "\"", ""), "\\", ""), " "))
	}
}

func sourceService(source string) string {
	source = strings.TrimPrefix(source, "/services/")
	if source == "" {
		return "unknown"
	}
	return source
}

func jsonString(value interface{}) string {
	body, err := json.Marshal(value)
	if err != nil {
		return ""
	}
	return string(body)
}

func shortID(value string) string {
	if len(value) <= 12 {
		return value
	}
	return value[:8]
}
