package fault

import (
	"errors"
	"fmt"
	"log/slog"
	"runtime"
	"strings"
)

// Code representa a classificação semântica do erro no domínio.
// Usar string permite serialização direta em logs e respostas JSON.
type Code string

const (
	CodeBadRequest          Code = "ERR_400"
	CodeUnauthorized        Code = "ERR_401"
	CodeForbidden           Code = "ERR_403"
	CodeNotFound            Code = "ERR_404"
	CodeConflict            Code = "ERR_409"
	CodeInternal            Code = "ERR_500"
	CodeUnprocessableEntity Code = "ERR_422"
)

// DomainError é o erro rico do domínio.
// Implementa a interface error e é compatível com errors.Is / errors.As / errors.Unwrap.
type DomainError struct {
	Code            Code           // Classificação semântica
	FriendlyMessage string         // Mensagem segura para expor ao client
	Fields          map[string]any // Erros por campo (ex: validação de formulário)
	Origin          string         // Função de origem, preenchida automaticamente
	Cause           error          // erro original preservado (permite wrapping)
}

// Error implementa a interface error.
// Retorna informação técnica completa — use apenas em logs internos.
func (e *DomainError) Error() string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("[%s]", e.Code))

	if e.Origin != "" {
		sb.WriteString(fmt.Sprintf(" %s", e.Origin))
	}

	if e.Cause != nil {
		sb.WriteString(fmt.Sprintf(" %s", e.Cause.Error()))
	}

	if len(e.Fields) > 0 {
		//sb.WriteString(fmt.Sprintf(" fields=%v", e.Fields))
		sb.WriteString(" [")
		for k, v := range e.Fields {
			sb.WriteString(fmt.Sprintf("%s=%v ", k, v))
		}
		sb.WriteString("]")
	}

	return sb.String()
}

// Is permite comparar DomainError por Code usando errors.Is.
//
// Exemplo:
//
//	var ErrNotFound = &DomainError{Code: CodeNotFound}
//	errors.Is(err, ErrNotFound) // true se ambos tiverem CodeNotFound
func (e *DomainError) Is(target error) bool {
	if t, ok := errors.AsType[*DomainError](target); ok {
		return e.Code == t.Code
	}
	return false
}

// Unwrap implementa a interface errors.Unwrap
func (e *DomainError) Unwrap() error {
	return e.Cause
}

// New cria um DomainError.
// origin é preenchido automaticamente com o nome da função chamadora.
func New(code Code, friendly string, cause error) *DomainError {
	return &DomainError{
		Code:            code,
		FriendlyMessage: friendly,
		Cause:           cause,
		Origin:          callerName(2),
	}
}

// NewWithFields cria um DomainError com mapa de erros por campo.
// Útil para erros de validação onde cada campo tem sua mensagem.
func NewWithFields(code Code, friendly string, fields map[string]any) *DomainError {
	return &DomainError{
		Code:            code,
		FriendlyMessage: friendly,
		Fields:          fields,
		Origin:          callerName(2),
	}
}

// Wrap envolve um erro existente num DomainError, preservando a causa original.
// Equivalente semântico ao fmt.Errorf("op: %w", err) mas com metadados ricos.
func Wrap(code Code, friendly string, cause error) *DomainError {
	return &DomainError{
		Code:            code,
		FriendlyMessage: friendly,
		Cause:           cause,
		Origin:          callerName(2),
	}
}

func callerName(skip int) string {
	pc, _, _, ok := runtime.Caller(skip)
	if !ok {
		return "unknown"
	}

	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return "unknown"
	}

	// runtime retorna "pkg/path.TypeName.MethodName" — pega só "TypeName.MethodName"
	full := fn.Name()

	if idx := strings.LastIndex(full, "/"); idx >= 0 {
		full = full[idx+1:]

	}
	return full
}

// Attrs retorna os atributos slog do DomainError para uso em logs.
// Se err não for DomainError, retorna só o erro como string.
func Attrs(err error) []any {
	if de, ok := errors.AsType[*DomainError](err); ok {
		return []any{
			slog.String("error_code", string(de.Code)),
			slog.String("error_origin", de.Origin),
			slog.String("error_message", de.FriendlyMessage),
		}
	}
	return []any{
		slog.String("error", err.Error()),
	}
}
