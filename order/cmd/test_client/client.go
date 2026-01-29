package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/agu3des/microservices-proto/golang/order"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

const (
	ProdutoExistente = "prod1" 
	ProdutoInexistente = "PRODUTO-FANTASMA-XYZ"
)

// Cores
const (
	Reset  = "\033[0m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Blue   = "\033[34m"
	Cyan   = "\033[36m"
)

func main() {
	conn, err := grpc.NewClient("localhost:3000", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := order.NewOrderClient(conn)

	fmt.Println("=============================================================")
	fmt.Println("   🚀 BATERIA DE TESTES (CORRIGIDA E BLINDADA)")
	fmt.Println("=============================================================")

	// 1. Produto Existente
	executarTeste(client, "1️⃣  CASO FELIZ (Tudo Válido)", ProdutoExistente, 2, 50.00)

	// 2. Produto Inexistente
	executarTeste(client, "2️⃣  PRODUTO NÃO CADASTRADO", ProdutoInexistente, 1, 100.00)

	// 3. Regras de Validação
	executarTeste(client, "3️⃣  QUANTIDADE ZERO", ProdutoExistente, 0, 50.00)
	executarTeste(client, "4️⃣  QUANTIDADE NEGATIVA", ProdutoExistente, -5, 50.00)
	executarTeste(client, "5️⃣  PREÇO NEGATIVO", ProdutoExistente, 1, -20.00)

	// 4. Limite
	executarTeste(client, "6️⃣  PEDIDO GIGANTE", ProdutoExistente, 10000, 1.50)

	fmt.Println("\n🏁 Fim da execução.")
}

func executarTeste(client order.OrderClient, titulo string, produto string, quantidade int32, preco float32) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	fmt.Printf("\n%s%s%s\n", Cyan, titulo, Reset)
	fmt.Printf("   📝 Enviando: Item=%s | Qtd=%d | Preço=%.2f\n", produto, quantidade, preco)

	req := &order.CreateOrderRequest{
		CostumerId: 1001,
		OrderItems: []*order.OrderItem{
			{
				ProductCode: produto,
				UnitPrice:   preco,
				Quantity:    quantidade,
			},
		},
	}

	res, err := client.Create(ctx, req)

	if err != nil {
		st, _ := status.FromError(err)
		code := st.Code()
		msg := st.Message()

		if code != codes.OK { 
			fmt.Printf("   %s✅ SUCESSO: O sistema bloqueou o pedido.%s\n", Green, Reset)
			fmt.Printf("      Erro retornado: [%s] %s\n", code, msg)
		} else {
			fmt.Printf("   ❌ ERRO ESTRANHO: gRPC retornou erro nil mas code não OK?\n")
		}

	} else {
		failCondition := false
		if produto == ProdutoInexistente { failCondition = true }
		if quantidade <= 0 { failCondition = true }
		if preco < 0 { failCondition = true }

		if failCondition {
			fmt.Printf("   %s❌ FALHA GRAVE: O sistema aceitou dados inválidos!%s\n", Red, Reset)
			fmt.Printf("      Order ID Criado: %d\n", res.OrderId)
		} else {
			fmt.Printf("   %s✅ SUCESSO: Pedido processado.%s Order ID: %d\n", Green, Reset, res.OrderId)
		}
	}
}