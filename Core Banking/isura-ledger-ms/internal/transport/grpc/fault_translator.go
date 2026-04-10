package grpc

import (
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/andreis3/isura-ledger-ms/internal/domain/fault"
)

// ProtocolError mapeia um Code de domínio para status de protocolos externos.
// Adicione novos protocolos aqui (gRPC, GraphQL, etc.) sem tocar no domínio.
type ProtocolError struct {
	GRPCCode codes.Code
}

// translator é o mapa interno de conversão Code → protocolos.
var translator = map[fault.Code]ProtocolError{
	fault.CodeBadRequest:          {GRPCCode: codes.InvalidArgument},
	fault.CodeUnauthorized:        {GRPCCode: codes.Unauthenticated},
	fault.CodeForbidden:           {GRPCCode: codes.PermissionDenied},
	fault.CodeNotFound:            {GRPCCode: codes.NotFound},
	fault.CodeConflict:            {GRPCCode: codes.AlreadyExists},
	fault.CodeUnprocessableEntity: {GRPCCode: codes.InvalidArgument},
	fault.CodeInternal:            {GRPCCode: codes.Internal},
}

// GRPCStatus retorna o status GRPC correspondente ao erro.
// Se o erro não for um DomainError, retorna 500.
// Se o Code não estiver mapeado, retorna 500.
func GRPCStatus(err error) codes.Code {
	var de *fault.DomainError
	if !errors.As(err, &de) {
		return codes.Internal
	}

	if p, ok := translator[de.Code]; ok {
		return p.GRPCCode
	}

	return codes.Internal
}

// Response é a estrutura que vai no body da resposta de erro HTTP.
// Nunca exponha DomainError.Error() aqui — contém informação técnica.
type Response struct {
	Code    fault.Code     `json:"code"`
	Message string         `json:"message"`
	Fields  map[string]any `json:"fields,omitempty"`
}

// ToGRPCError converte um DomainError para um erro GRPC.
func ToGRPCError(err error) error {
	if err == nil {
		return nil
	}

	code := GRPCStatus(err)
	response := ToResponse(err)
	return status.Error(code, response.Message)
}

// ToResponse converte um DomainError para a resposta segura ao client.
func ToResponse(err error) Response {
	var de *fault.DomainError
	if !errors.As(err, &de) {
		return Response{
			Code:    fault.CodeInternal,
			Message: "Internal server error",
		}
	}

	return Response{
		Code:    de.Code,
		Message: de.FriendlyMessage,
		Fields:  de.Fields,
	}
}
