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
    shippingServiceUrl string
}

func NewAdapter(shippingServiceUrl string) (*Adapter, error) {
    return &Adapter{shippingServiceUrl: shippingServiceUrl}, nil
}

func (a *Adapter) Ship(order *domain.Order) error {
    var opts []grpc.DialOption
    opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

    conn, err := grpc.NewClient(a.shippingServiceUrl, opts...)
    if err != nil {
        return err
    }
    defer conn.Close()

    client := shipping.NewShippingClient(conn)

    var items []*shipping.ShippingItem
    for _, item := range order.OrderItems {
        items = append(items, &shipping.ShippingItem{
            ProductCode: item.ProductCode,
            Quantity:    item.Quantity,
        })
    }

    resp, err := client.Create(context.Background(), &shipping.CreateShippingRequest{
        OrderId: order.ID,
        Items:   items,
    })

    if err != nil {
        return fmt.Errorf("erro no serviço de shipping: %v", err)
    }

    order.DeliveryDays = resp.DeliveryDays
    
    return nil
}