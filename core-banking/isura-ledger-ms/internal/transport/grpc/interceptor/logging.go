package interceptor

import (
	"context"
	"log/slog"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// LoggingInterceptor é um interceptor unário que loga cada request gRPC.
// Deve ser registrado no grpc.NewServer via grpc.UnaryInterceptor.
func LoggingInterceptor(log *slog.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		start := time.Now()

		// chama o handler real

		resp, err := handler(ctx, req)

		duration := time.Since(start)
		code := status.Code(err)

		attrs := []any{
			slog.String("method", info.FullMethod),
			slog.String("status", code.String()),
			slog.String("duration", duration.String()),
		}

		// codes que são erros de cliente — loga como WARN
		// codes que são erros de servidor — loga como ERROR
		switch code {
		case codes.Internal, codes.Unavailable, codes.DataLoss, codes.Unknown:
			log.ErrorContext(ctx, "gRPC request failed", attrs...)
		case codes.OK:
			log.InfoContext(ctx, "gRPC request completed", attrs...)
		default:
			log.WarnContext(ctx, "gRPC request rejected", attrs...)
		}

		return resp, err
	}
}
