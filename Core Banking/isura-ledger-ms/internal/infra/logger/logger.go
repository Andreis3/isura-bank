package logger

import (
	"context"
	"log/slog"
	"os"

	"github.com/lmittmann/tint"
	"go.opentelemetry.io/otel/trace"
)

// LevelCritical Níveis customizados
const LevelCritical = slog.LevelError + 1

const (
	timeFormat = "01-02-2006 15:04:05.000"
)

// Logger encapsula os handlers para diferentes ambientes.
type Logger struct {
	json *slog.Logger
	text *slog.Logger
}

// NewLogger constrói os handlers e define o padrão global.
func NewLogger() *Logger {
	// Opções para JSON (Produção)
	jsonOpts := &slog.HandlerOptions{
		Level:       slog.LevelDebug,
		ReplaceAttr: replaceAttrJSON,
	}

	// Opções para Texto (Desenvolvimento)
	textOpts := &tint.Options{
		Level:      slog.LevelDebug,
		TimeFormat: timeFormat,
		NoColor:    false,
	}

	l := &Logger{
		json: slog.New(slog.NewJSONHandler(os.Stdout, jsonOpts)),
		text: slog.New(tint.NewHandler(os.Stderr, textOpts)),
	}

	// Define o default do sistema baseado no ambiente
	if os.Getenv("ENV") == "development" {
		slog.SetDefault(l.text)
	} else {
		slog.SetDefault(l.json)
	}

	return l
}

// ── Métodos JSON (Produção / Loki) ───────────────────────────────────────────

func (l *Logger) DebugJSON(msg string, args ...any) { l.json.Debug(msg, args...) }
func (l *Logger) InfoJSON(msg string, args ...any)  { l.json.Info(msg, args...) }
func (l *Logger) WarnJSON(msg string, args ...any)  { l.json.Warn(msg, args...) }
func (l *Logger) ErrorJSON(msg string, args ...any) { l.json.Error(msg, args...) }
func (l *Logger) CriticalJSON(msg string, args ...any) {
	l.json.Log(context.Background(), LevelCritical, msg, args...)
}

// ── Métodos Text (Desenvolvimento / Terminal) ────────────────────────────────

func (l *Logger) DebugText(msg string, args ...any) { l.text.Debug(msg, args...) }
func (l *Logger) InfoText(msg string, args ...any)  { l.text.Info(msg, args...) }
func (l *Logger) WarnText(msg string, args ...any)  { l.text.Warn(msg, args...) }
func (l *Logger) ErrorText(msg string, args ...any) { l.text.Error(msg, args...) }
func (l *Logger) CriticalText(msg string, args ...any) {
	l.text.Log(context.Background(), LevelCritical, msg, args...)
}

// ── Observabilidade ───────────────────────────────────────────────────────────

// WithTrace retorna o logger JSON enriquecido com dados do OpenTelemetry.
func (l *Logger) WithTrace(ctx context.Context) *slog.Logger {
	span := trace.SpanContextFromContext(ctx)
	if !span.IsValid() || !span.HasTraceID() {
		return l.json
	}
	return l.json.With(
		slog.String("trace_id", span.TraceID().String()),
		slog.String("span_id", span.SpanID().String()),
	)
}

// ── Helpers e Atributos ───────────────────────────────────────────────────────

func replaceAttrJSON(_ []string, a slog.Attr) slog.Attr {
	switch a.Key {
	case slog.TimeKey:
		// Otimização: usa o tempo que o slog já fornece
		a.Value = slog.StringValue(a.Value.Time().Format(timeFormat))
	case slog.LevelKey:
		a.Value = slog.StringValue(levelLabel(a))
	}
	return a
}

func levelLabel(a slog.Attr) string {
	level, ok := a.Value.Any().(slog.Level)
	if !ok {
		return a.Value.String()
	}
	switch level {
	case LevelCritical:
		return "CRITICAL"
	case slog.LevelDebug:
		return "DEBUG"
	case slog.LevelInfo:
		return "INFO"
	case slog.LevelWarn:
		return "WARN"
	case slog.LevelError:
		return "ERROR"
	default:
		return level.String()
	}
}

// SlogJSON e SlogText permitem acesso direto caso necessário
func (l *Logger) SlogJSON() *slog.Logger { return l.json }
func (l *Logger) SlogText() *slog.Logger { return l.text }
