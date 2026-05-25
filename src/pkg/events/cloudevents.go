package events

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"go.opentelemetry.io/otel/trace"
)

const (
	SourceFarmService      = "/services/farm-service"
	SourceRetailService    = "/services/retail-service"
	SourcePaymentService   = "/services/payment-service"
	SourceWarehouseService = "/services/warehouse-service"
	SourceLogisticsService = "/services/logistics-service"
)

type Metadata struct {
	EventID       string
	CorrelationID string
	CausationID   string
	TraceID       string
	OccurredAt    time.Time
	OrgID         string
	OrderID       string
	PaymentID     string
	StoreID       string
	ShipmentID    string
	HarvestID     string
	BatchID       string
	FarmID        string
	WarehouseID   string
	DriverID      string
	VehicleID     string
}

func NewCloudEvent(ctx context.Context, eventType string, source string, subject string, data any, metadata Metadata) (cloudevents.Event, error) {
	event := cloudevents.NewEvent()
	event.SetSpecVersion(cloudevents.VersionV1)
	event.SetType(eventType)
	event.SetSource(source)
	event.SetSubject(subject)
	if metadata.EventID == "" {
		return event, errors.New("event id is required")
	}
	event.SetID(metadata.EventID)
	if metadata.OccurredAt.IsZero() {
		metadata.OccurredAt = time.Now()
	}
	event.SetTime(metadata.OccurredAt)
	if metadata.CorrelationID == "" {
		return event, errors.New("correlationid is required")
	}
	setExtension(&event, "correlationid", metadata.CorrelationID)
	setExtension(&event, "causationid", metadata.CausationID)
	if metadata.TraceID == "" {
		metadata.TraceID = TraceIDFromContext(ctx)
	}
	setExtension(&event, "traceid", metadata.TraceID)
	setExtension(&event, "orgid", metadata.OrgID)
	setExtension(&event, "orderid", metadata.OrderID)
	setExtension(&event, "paymentid", metadata.PaymentID)
	setExtension(&event, "storeid", metadata.StoreID)
	setExtension(&event, "shipmentid", metadata.ShipmentID)
	setExtension(&event, "harvestid", metadata.HarvestID)
	setExtension(&event, "batchid", metadata.BatchID)
	setExtension(&event, "farmid", metadata.FarmID)
	setExtension(&event, "warehouseid", metadata.WarehouseID)
	setExtension(&event, "driverid", metadata.DriverID)
	setExtension(&event, "vehicleid", metadata.VehicleID)
	if err := event.SetData(cloudevents.ApplicationJSON, data); err != nil {
		return event, err
	}
	if err := event.Validate(); err != nil {
		return event, err
	}
	return event, nil
}

func ParseCloudEvent(payload []byte) (cloudevents.Event, error) {
	var event cloudevents.Event
	if err := json.Unmarshal(payload, &event); err != nil {
		return event, fmt.Errorf("invalid CloudEvent JSON: %w", err)
	}
	if err := event.Validate(); err != nil {
		return event, fmt.Errorf("invalid CloudEvent: %w", err)
	}
	if event.DataContentType() != cloudevents.ApplicationJSON {
		return event, fmt.Errorf("invalid CloudEvent datacontenttype %q", event.DataContentType())
	}
	if ExtensionString(event, "correlationid") == "" {
		return event, errors.New("invalid CloudEvent: correlationid extension is required")
	}
	return event, nil
}

func DataAs[T any](event cloudevents.Event) (T, error) {
	var out T
	if err := event.DataAs(&out); err != nil {
		return out, err
	}
	return out, nil
}

func DataMap(event cloudevents.Event) map[string]interface{} {
	var out map[string]interface{}
	if err := event.DataAs(&out); err != nil {
		return map[string]interface{}{}
	}
	return out
}

func ExtensionString(event cloudevents.Event, key string) string {
	value, ok := event.Extensions()[key]
	if !ok || value == nil {
		return ""
	}
	switch typed := value.(type) {
	case string:
		return typed
	default:
		return fmt.Sprint(typed)
	}
}

func TraceIDFromContext(ctx context.Context) string {
	spanContext := trace.SpanContextFromContext(ctx)
	if spanContext.HasTraceID() {
		return spanContext.TraceID().String()
	}
	return ""
}

func setExtension(event *cloudevents.Event, key string, value string) {
	if value != "" {
		event.SetExtension(key, value)
	}
}
