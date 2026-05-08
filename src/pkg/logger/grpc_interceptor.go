package logger

import (
	"context"
	"time"

	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

// GRPCLoggerInterceptor returns a grpc.UnaryServerInterceptor that logs requests using zap.
func GRPCLoggerInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		start := time.Now()

		resp, err := handler(ctx, req)

		latency := time.Since(start)
		st, _ := status.FromError(err)
		
		log := FromContext(ctx)
		
		fields := []zap.Field{
			zap.String("grpc.method", info.FullMethod),
			zap.String("grpc.code", st.Code().String()),
			zap.Duration("latency", latency),
		}

		spanContext := trace.SpanContextFromContext(ctx)
		if spanContext.HasTraceID() {
			fields = append(fields, zap.String("trace_id", spanContext.TraceID().String()))
		}

		if err != nil {
			log.Error("gRPC request failed", append(fields, zap.Error(err))...)
		} else {
			log.Info("gRPC request success", fields...)
		}

		return resp, err
	}
}
