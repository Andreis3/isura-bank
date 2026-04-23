package grpc

import "google.golang.org/grpc"

// Module define o contrato para registrar handlers no servidor gRPC.
type Module interface {
	Register(server *grpc.Server)
}
