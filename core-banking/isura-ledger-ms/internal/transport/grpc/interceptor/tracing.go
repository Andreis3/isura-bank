package interceptor

import (
	"context"

	"google.golang.org/grpc"

	"github.com/andreis3/isura-ledger-ms/internal/application"
)

func TracingInterceptor(tracer application.Tracer) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		ctx, span := tracer.Start(ctx, info.FullMethod)
		defer span.End()

		resp, err := handler(ctx, req)
		if err != nil {
			span.RecordError(err)
		}

		return resp, err
	}
}
