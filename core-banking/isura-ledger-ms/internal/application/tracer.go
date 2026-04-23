package application

import "context"

type Tracer interface {
	Start(ctx context.Context, spanName string) (context.Context, Span)
}

type Span interface {
	End()
	SpanContext() SpanContext
	RecordError(err error)
}

type SpanContext interface {
	TraceID() string
}
