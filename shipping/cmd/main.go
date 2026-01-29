package main

import (
	"github.com/agu3des/microservices/shipping/internal/adapters/grpc"
	"github.com/agu3des/microservices/shipping/internal/application/core/api"
)

func main() {
	// Shipping roda na porta 50052 (pois Payment é 50051)
	application := api.NewApplication()
	grpcAdapter := grpc.NewAdapter(application, 50052)
	grpcAdapter.Run()
}