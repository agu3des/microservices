package main

import (
	"log"
	"os"
	"strconv"

	"github.com/agu3des/microservices/order/internal/adapters/db"
	"github.com/agu3des/microservices/order/internal/adapters/grpc"
	"github.com/agu3des/microservices/order/internal/adapters/payment"
	"github.com/agu3des/microservices/order/internal/adapters/shipping"
	"github.com/agu3des/microservices/order/internal/application/core/api"
)

func main() {
	// 1. Configuração de Variáveis de Ambiente
	dataSourceUrl := os.Getenv("DATA_SOURCE_URL")
	if dataSourceUrl == "" {
		dataSourceUrl = "root:root@tcp(localhost:3306)/order"
	}

	paymentServiceUrl := os.Getenv("PAYMENT_SERVICE_URL")
	if paymentServiceUrl == "" {
		paymentServiceUrl = "localhost:50051"
	}

	shippingServiceUrl := os.Getenv("SHIPPING_SERVICE_URL")
	if shippingServiceUrl == "" {
		shippingServiceUrl = "localhost:50052"
	}

	applicationPortStr := os.Getenv("APPLICATION_PORT")
	if applicationPortStr == "" {
		applicationPortStr = "3000"
	}

	applicationPort, err := strconv.Atoi(applicationPortStr)
	if err != nil {
		log.Fatalf("Porta inválida: %v", err)
	}

	// 2. Inicializar Adapter de Banco de Dados
	dbAdapter, err := db.NewAdapter(dataSourceUrl)
	if err != nil {
		log.Fatalf("Falha ao conectar no banco de dados: %v", err)
	}

	// 3. Inicializar Adapter de Pagamento (gRPC Client)
	paymentAdapter, err := payment.NewAdapter(paymentServiceUrl)
	if err != nil {
		log.Fatalf("Falha ao conectar no serviço de pagamento: %v", err)
	}

	// 4. Inicializar Adapter de Envio/Shipping (gRPC Client)
	shippingAdapter := shipping.NewAdapter(shippingServiceUrl)

	// 5. Inicializar o Core da Aplicação (Injeção de Dependência)
	application := api.NewApplication(dbAdapter, paymentAdapter, shippingAdapter)

	// 6. Inicializar e Rodar o Servidor gRPC (Porta de Entrada)
	grpcAdapter := grpc.NewAdapter(application, applicationPort)
	grpcAdapter.Run()
}