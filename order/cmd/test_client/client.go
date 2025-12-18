package main

import (
	"context"
	"fmt"
	"log"
	"time"

	// Importação baseada no seu script 'run' (MODULE_PATH)
	"github.com/agu3des/microservices-proto/golang/order"
	
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

func main() {
	// 1. Conecta ao Microsserviço de ORDER
	// Ajuste a porta ":3000" se o seu config estiver usando outra
	conn, err := grpc.NewClient("localhost:3000", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := order.NewOrderClient(conn)

	fmt.Println("--- INICIANDO TESTES DE INTEGRAÇÃO ---")

	// ---------------------------------------------------------
	// CENÁRIO 1: Pedido Válido (Deve passar no Order e tentar chamar Payment)
	// Nota: Para este teste dar 100% certo, o microsserviço PAYMENT deve estar rodando também.
	// ---------------------------------------------------------
	fmt.Println("\n1️⃣  Tentando criar pedido VÁLIDO (10 itens)...")
	criarPedido(client, 10, 10.0)

	// ---------------------------------------------------------
	// CENÁRIO 2: Regra de Negócio (Limite de 50 itens)
	// Este teste deve falhar APENAS no Order, sem nem chamar o Payment.
	// ---------------------------------------------------------
	fmt.Println("\n2️⃣  Tentando criar pedido INVÁLIDO (>50 itens)...")
	criarPedido(client, 60, 5.0)
}

func criarPedido(client order.OrderClient, quantidade int32, preco float32) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	// Simulando um request
	req := &order.CreateOrderRequest{
		UserId: 1001, // ID do usuário
		OrderItems: []*order.OrderItem{
			{
				ProductCode: "PROD-TEST-01",
				UnitPrice:   preco,
				Quantity:    quantidade,
			},
		},
	}

	res, err := client.Create(ctx, req)

	if err != nil {
		// Analisa o erro retornado pelo gRPC
		st, ok := status.FromError(err)
		if ok {
			switch st.Code() {
			case codes.InvalidArgument:
				fmt.Printf("✅ SUCESSO NO TESTE DE BLOQUEIO! O servidor retornou InvalidArgument: %s\n", st.Message())
			case codes.Internal:
				fmt.Printf("⚠️  Erro Interno (pode ser falha ao conectar no Payment): %s\n", st.Message())
			default:
				fmt.Printf("❌ Erro inesperado: [%s] %s\n", st.Code(), st.Message())
			}
		} else {
			log.Fatalf("Erro não-gRPC: %v", err)
		}
	} else {
		fmt.Printf("✅ Pedido criado com sucesso! Order ID: %d\n", res.OrderId)
	}
}