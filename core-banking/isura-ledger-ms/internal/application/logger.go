package application

import (
	"context"
	"log/slog"
)

type Logger interface {
	DebugJSON(msg string, info ...any)
	InfoJSON(msg string, info ...any)
	WarnJSON(msg string, info ...any)
	ErrorJSON(msg string, info ...any)
	CriticalJSON(msg string, info ...any)
	DebugText(msg string, info ...any)
	InfoText(msg string, info ...any)
	WarnText(msg string, info ...any)
	ErrorText(msg string, info ...any)
	CriticalText(msg string, info ...any)
	WithTrace(ctx context.Context) *slog.Logger
	SlogJSON() *slog.Logger
	SlogText() *slog.Logger
}
