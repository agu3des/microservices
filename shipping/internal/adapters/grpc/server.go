package grpc

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/agu3des/microservices-proto/golang/shipping"
	"github.com/agu3des/microservices/shipping/internal/application/core/domain"
	"github.com/agu3des/microservices/shipping/internal/ports"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Adapter struct {
	api  ports.APIPort
	port int
	shipping.UnimplementedShippingServer
}

func NewAdapter(api ports.APIPort, port int) *Adapter {
	return &Adapter{api: api, port: port}
}

func (a Adapter) Create(ctx context.Context, req *shipping.CreateShippingRequest) (*shipping.CreateShippingResponse, error) {
	var items []domain.ShippingItem
	for _, i := range req.Items {
		items = append(items, domain.ShippingItem{
			ProductCode: i.ProductCode,
			Quantity:    i.Quantity,
		})
	}

	days, err := a.api.ShipOrder(req.OrderId, items)
	if err != nil {
		return nil, err
	}

	return &shipping.CreateShippingResponse{
		OrderId:      req.OrderId,
		DeliveryDays: days,
	}, nil
}

func (a Adapter) Run() {
	listen, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		log.Fatalf("failed to listen on port %d: %v", a.port, err)
	}
	grpcServer := grpc.NewServer()
	shipping.RegisterShippingServer(grpcServer, a)
	reflection.Register(grpcServer)

	log.Printf("Starting Shipping Service on port %d...", a.port)
	if err := grpcServer.Serve(listen); err != nil {
		log.Fatalf("failed to serve grpc: %v", err)
	}
}