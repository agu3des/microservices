package main

import (
	"log"

	"github.com/agu3des/microservices/order/config"
	"github.com/agu3des/microservices/order/internal/adapters/db"
	"github.com/agu3des/microservices/order/internal/adapters/grpc"
	"github.com/agu3des/microservices/order/internal/application/core/api"
)

func main() {
	// 1. Inicializa o Adapter do Banco de Dados (Driven Adapter)
	dbAdapter, err := db.NewAdapter(config.GetDataSourceURL())
	if err != nil {
		log.Fatalf("Failed to connect to database. Error: %v", err)
	}

	// 2. Cria a Aplicação Core injetando a dependência do banco (Dependency Injection)
	application := api.NewApplication(dbAdapter)

	// 3. Inicializa o Adapter gRPC (Driver Adapter) injetando a aplicação
	grpcAdapter := grpc.NewAdapter(application, config.GetApplicationPort())

	// 4. Inicia o servidor
	grpcAdapter.Run()
}