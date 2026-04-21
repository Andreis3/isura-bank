package interceptor

import (
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"

	"github.com/andreis3/isura-ledger-ms/internal/application"
)

func MetricsInterceptor(metrics application.Metrics) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		start := time.Now()

		resp, err := handler(ctx, req)

		duration := float64(time.Since(start).Milliseconds())
		code := int(status.Code(err))

		metrics.RecordRequestTotal(info.FullMethod, "grpc", code)
		metrics.RecordRequestDuration(info.FullMethod, "grpc", code, duration)

		return resp, err
	}
}
