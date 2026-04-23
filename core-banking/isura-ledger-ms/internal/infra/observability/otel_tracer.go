package observability

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	sdkresource "go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.40.0"
	"go.opentelemetry.io/otel/trace"

	"github.com/andreis3/isura-ledger-ms/internal/application"
	"github.com/andreis3/isura-ledger-ms/internal/infra/configs"
)

type otelTracer struct {
	tracer trace.Tracer
}

func (o *otelTracer) Start(ctx context.Context, spanName string) (context.Context, application.Span) {
	ctx, s := o.tracer.Start(ctx, spanName)
	return ctx, &otelSpan{s}
}

type otelSpan struct {
	trace.Span
}

func (s *otelSpan) End() {
	s.Span.End()
}

func (s *otelSpan) RecordError(err error) {
	s.Span.RecordError(err)
}

func (s *otelSpan) SpanContext() application.SpanContext {
	return &otelSpanContext{s.Span.SpanContext()} // ✅ chama o SpanContext do campo do OTEL
}

type otelSpanContext struct {
	trace.SpanContext
}

func (sc *otelSpanContext) TraceID() string {
	return sc.SpanContext.TraceID().String()
}

// InitOtelTracer ✅ Esta função inicializa o OpenTelemetry por completo e retorna o adapter. Tracer
func InitOtelTracer(ctx context.Context, cfg *configs.Configs) (application.Tracer, func(context.Context) error, error) {
	exporter, err := otlptracehttp.New(
		ctx,
		otlptracehttp.WithEndpoint(cfg.OpenTelemetry.Host),
		otlptracehttp.WithInsecure(),
		otlptracehttp.WithCompression(otlptracehttp.GzipCompression),
	)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create OTLP exporter: %w", err)
	}

	resource, err := sdkresource.New(ctx,
		sdkresource.WithAttributes(
			semconv.ServiceName(cfg.ApplicationName),
			semconv.ServiceVersion(cfg.Version),
		),
		sdkresource.WithProcess(),
		sdkresource.WithOS(),
		sdkresource.WithHost(),
	)

	if err != nil {
		return nil, nil, fmt.Errorf("failed to create resource: %w", err)
	}

	sampler := sdktrace.AlwaysSample()
	if cfg.Env == "production" {
		sampler = sdktrace.TraceIDRatioBased(0.1)
	}

	provider := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(resource),
		sdktrace.WithSampler(sampler),
	)

	otel.SetTracerProvider(provider)

	// Retorna o adapter.Tracer e a função de shutdown
	return &otelTracer{tracer: provider.Tracer(cfg.ApplicationName)}, provider.Shutdown, nil
}
