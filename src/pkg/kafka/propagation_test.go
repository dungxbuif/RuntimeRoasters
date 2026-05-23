package kafka

import (
	"context"
	"testing"

	"github.com/segmentio/kafka-go"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	oteltrace "go.opentelemetry.io/otel/trace"
)

func TestKafkaHeadersCarrierInjectExtractTraceparent(t *testing.T) {
	otel.SetTextMapPropagator(propagation.TraceContext{})
	tp := sdktrace.NewTracerProvider()
	defer func() {
		_ = tp.Shutdown(context.Background())
	}()
	otel.SetTracerProvider(tp)

	ctx, span := otel.Tracer("test").Start(context.Background(), "produce", oteltrace.WithSpanKind(oteltrace.SpanKindProducer))
	defer span.End()

	headers := kafkaHeadersCarrier{}
	otel.GetTextMapPropagator().Inject(ctx, &headers)

	traceparent := headers.Get("traceparent")
	if traceparent == "" {
		t.Fatal("expected traceparent header to be injected")
	}

	extracted := otel.GetTextMapPropagator().Extract(context.Background(), &headers)
	extractedSpanContext := oteltrace.SpanContextFromContext(extracted)
	if !extractedSpanContext.HasTraceID() {
		t.Fatal("expected extracted context to contain trace id")
	}
	if extractedSpanContext.TraceID() != span.SpanContext().TraceID() {
		t.Fatalf("expected extracted trace id %s, got %s", span.SpanContext().TraceID(), extractedSpanContext.TraceID())
	}
}

func TestKafkaHeadersCarrierSetOverwritesExistingHeader(t *testing.T) {
	headers := kafkaHeadersCarrier{
		{Key: "traceparent", Value: []byte("old")},
		{Key: "other", Value: []byte("value")},
	}

	headers.Set("traceparent", "new")

	if got := headers.Get("traceparent"); got != "new" {
		t.Fatalf("expected overwritten traceparent, got %q", got)
	}
	if len([]kafka.Header(headers)) != 2 {
		t.Fatalf("expected header count to remain 2, got %d", len(headers))
	}
}
