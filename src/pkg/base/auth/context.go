package auth

import (
	"context"

	"github.com/dungxbuif/RuntimeRoasters/pkg/base/identity"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)
func SetIdentityInContext(ctx context.Context, idClaims identity.Claims) context.Context {
	return identity.InjectContext(ctx, idClaims)
}

func RecordTracingData(ctx context.Context, sub string) {
	span := trace.SpanFromContext(ctx)
	if span.IsRecording() && sub != "" {
		span.SetAttributes(attribute.String(TraceAttrUserID, sub))
	}
}