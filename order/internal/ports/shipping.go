package ports

import "github.com/agu3des/microservices/order/internal/application/core/domain"

type ShippingPort interface {
	Ship(order *domain.Order) error
}