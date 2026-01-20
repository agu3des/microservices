package payment

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/agu3des/microservices/order/internal/application/core/domain"
	
	"github.com/agu3des/microservices-proto/golang/payment"
	
	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

type Adapter struct {
	paymentServiceUrl string
}

func NewAdapter(paymentServiceUrl string) (*Adapter, error) {
	return &Adapter{
		paymentServiceUrl: paymentServiceUrl,
	}, nil
}

func (a *Adapter) Charge(order *domain.Order) error {
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(grpc_retry.UnaryClientInterceptor(
			grpc_retry.WithCodes(codes.Unavailable, codes.ResourceExhausted),
			grpc_retry.WithMax(5),
			grpc_retry.WithBackoff(grpc_retry.BackoffLinear(1 * time.Second)),
		)),
	}

	conn, err := grpc.NewClient(a.paymentServiceUrl, opts...)
	if err != nil {
		return fmt.Errorf("falha ao conectar no serviço de pagamento: %v", err)
	}
	defer conn.Close()

	client := payment.NewPaymentClient(conn)

	req := &payment.CreatePaymentRequest{
		UserId:     order.CustomerID,
		OrderId:    order.ID,
		TotalPrice: order.TotalPrice(),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	_, err = client.Create(ctx, req)
	if err != nil {
		if status.Code(err) == codes.DeadlineExceeded {
			log.Println("LOG: Timeout exceeded calling Payment Service (DeadlineExceeded)")
		}
		return fmt.Errorf("erro no pagamento: %v", err)
	}

	fmt.Printf("Pagamento autorizado externamente para o pedido %d\n", order.ID)
	return nil
}