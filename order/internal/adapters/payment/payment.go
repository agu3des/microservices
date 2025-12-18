package payment

import (
	"context"

	"github.com/agu3des/microservices-proto/golang/payment"
	"github.com/agu3des/microservices/order/internal/application/core/domain"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Adapter struct {
	payment payment.PaymentClient // come from the generated code by the protobuf compiler
}

func NewAdapter(paymentServiceUrl string) (*Adapter, error) {
	var opts []grpc.DialOption
	
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	
	conn, err := grpc.Dial(paymentServiceUrl, opts...)
	if err != nil {
		return nil, err
	}
	
	// Nota: Em uma implementação real, você precisará gerenciar o fechamento 
	// desta conexão (conn.Close()) quando a aplicação encerrar.
	
	client := payment.NewPaymentClient(conn) // initialize the stub
	return &Adapter{payment: client}, nil
}

func (a *Adapter) Charge(order *domain.Order) error {
	_, err := a.payment.Create(context.Background(), &payment.CreatePaymentRequest{
		UserId:     order.CustomerID,
		OrderId:    order.ID,
		TotalPrice: order.TotalPrice(),
	})
	
	return err
}