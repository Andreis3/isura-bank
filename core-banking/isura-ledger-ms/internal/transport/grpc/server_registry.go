package grpc

import "google.golang.org/grpc"

type ServerRegistry struct {
	server  *grpc.Server
	modules []Module
}

func NewServerRegistry(server *grpc.Server, modules ...Module) *ServerRegistry {
	return &ServerRegistry{
		server:  server,
		modules: modules,
	}
}

func (r *ServerRegistry) RegisterAll() {
	for _, module := range r.modules {
		module.Register(r.server)
	}
}
