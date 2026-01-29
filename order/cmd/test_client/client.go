package main

import (
	"context"
	"fmt"
	"log"
	"time"
	"strings"

	"github.com/agu3des/microservices-proto/golang/order"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

// Constantes de Teste
const (
	ProdutoExistente   = "prod1"                // Deve passar (se saldo permitir)
	ProdutoInexistente = "PRODUTO-FANTASMA-XYZ" // Deve falhar no Payment
)

// Cores para o Terminal
const (
	Reset  = "\033[0m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Blue   = "\033[34m"
	Cyan   = "\033[36m"
)

func main() {
	// 1. Conexão gRPC
	conn, err := grpc.NewClient("localhost:3000", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := order.NewOrderClient(conn)

	fmt.Println("=============================================================")
	fmt.Println("    🚀 BATERIA DE TESTES DE INTEGRAÇÃO (ORDER SERVICE)      ")
	fmt.Println("=============================================================")

	// --- CAMINHOS FELIZES (Deve Funcionar) ---

	// Cenário 1: Pedido Pequeno
	// Lógica: 2 itens -> Total < 5. Base 1 dia.
	// Resultado Esperado: 1 dia
	executarTeste(client, "1️⃣  PEDIDO PEQUENO (2 itens)", ProdutoExistente, 2, 50.00)

	// Cenário 2: Pedido Médio
	// Lógica: 10 itens -> 10/5 = 2 dias extras + 1 base.
	// Resultado Esperado: 3 dias
	executarTeste(client, "2️⃣  PEDIDO MÉDIO (10 itens)", ProdutoExistente, 10, 50.00)
	
	// --- CENÁRIOS DE ERRO/VALIDAÇÃO (Deve Bloquear) ---
	
	fmt.Println("\n--- TESTES DE VALIDAÇÃO (DEVEM FALHAR) ---")

	// Cenário 3: Pedido Grande
	// Lógica: 45 itens -> 45/5 = 9 dias extras + 1 base.
	// Resultado Esperado: 10 dias
	executarTeste(client, "3️⃣  PEDIDO GRANDE (45 itens)", ProdutoExistente, 45, 50.00)

	// Cenário 4: Estouro de Limite
	// Lógica: > 50 itens. O Domain do Order deve bloquear.
	executarTeste(client, "4️⃣  ESTOURO DE LIMITE (>50)", ProdutoExistente, 55, 1.50)

	// Cenário 5: Produto Inexistente
	// Lógica: O Payment Service vai rejeitar.
	executarTeste(client, "5️⃣  PRODUTO INEXISTENTE", ProdutoInexistente, 1, 100.00)

	// Cenário 6: Preço Inválido
	// Lógica: O Domain do Order deve bloquear preço negativo.
	executarTeste(client, "6️⃣  PREÇO NEGATIVO", ProdutoExistente, 1, -50.00)

	fmt.Println("\n=============================================================")
	fmt.Println(" 🏁 Fim da execução.")
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
		msg := st.Message() // Mensagem técnica completa

		// --- LIMPEZA DA MENSAGEM ---
		if strings.Contains(msg, "1000") || strings.Contains(msg, "Payment over") {
			fmt.Printf("   %s⛔ Pagamentos acima de 1000 reais não são permitidos.%s\n", Yellow, Reset)
		} else if strings.Contains(msg, "50 itens") {
			fmt.Printf("   %s⛔ Limite de quantidade atingido (Max 50).%s\n", Yellow, Reset)
		} else if strings.Contains(msg, "não encontrado") || strings.Contains(msg, "not found") {
			fmt.Printf("   %s⛔ Produto não encontrado no catálogo.%s\n", Yellow, Reset)
		} else if strings.Contains(msg, "negativo") || strings.Contains(msg, "negative") {
			fmt.Printf("   %s⛔ Preço inválido (negativo).%s\n", Yellow, Reset)
		} else {
			fmt.Printf("   %s❌ %s%s\n", Red, msg, Reset)
		}

	} else {
		fmt.Printf("   %s✅ PEDIDO APROVADO!%s\n", Green, Reset)
		fmt.Printf("      🆔 Order ID: %d\n", res.OrderId)
		
		if res.DeliveryDays > 0 {
			fmt.Printf("      📦 Entrega estimada em: %s%d dias%s\n", Yellow, res.DeliveryDays, Reset)
		} else {
			fmt.Printf("      ⚠️  Entrega: %sNão informada%s\n", Red, Reset)
		}
	}
}