package shipping

import (
	"context"
	"fmt"
	"github.com/agu3des/microservices-proto/golang/shipping"
	"github.com/agu3des/microservices/order/internal/application/core/domain"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Adapter struct {
	shippingURL string
}

func NewAdapter(url string) *Adapter {
	return &Adapter{shippingURL: url}
}

func (a *Adapter) Ship(order domain.Order) error {
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	conn, err := grpc.NewClient(a.shippingURL, opts...)
	if err != nil {
		return err
	}
	defer conn.Close()

	client := shipping.NewShippingClient(conn)

	var items []*shipping.ShippingItem
	for _, i := range order.OrderItems {
		items = append(items, &shipping.ShippingItem{
			ProductCode: i.ProductCode,
			Quantity:    i.Quantity,
		})
	}

	resp, err := client.Create(context.Background(), &shipping.CreateShippingRequest{
		OrderId: order.ID,
		Items:   items,
	})

	if err != nil {
		return fmt.Errorf("erro no envio: %v", err)
	}

	fmt.Printf("Envio processado! Prazo: %d dias para o pedido %d\n", resp.DeliveryDays, order.ID)
	return nil
}