package ports

import "github.com/agu3des/microservices/shipping/internal/application/core/domain"

type APIPort interface {
	ShipOrder(orderId int64, items []domain.ShippingItem) (int32, error)
}