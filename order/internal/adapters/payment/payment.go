package payment

import (
	"context"
	"fmt"

	"github.com/agu3des/microservices-proto/golang/payment"
	"github.com/agu3des/microservices/order/internal/application/core/domain"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Adapter struct {
	paymentServiceUrl string
}

func NewAdapter(paymentServiceUrl string) (*Adapter, error) {
	return &Adapter{paymentServiceUrl: paymentServiceUrl}, nil
}

func (a *Adapter) Charge(order *domain.Order) error {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

	conn, err := grpc.NewClient(a.paymentServiceUrl, opts...)
	if err != nil {
		return err
	}
	defer conn.Close()

	client := payment.NewPaymentClient(conn)

	_, err = client.Create(context.Background(), &payment.CreatePaymentRequest{
		UserId:     order.CustomerID,
		OrderId:    order.ID,
		TotalPrice: order.TotalPrice, 
	})

	if err != nil {
		return fmt.Errorf("erro no serviço de pagamento: %v", err)
	}

	return nil
}