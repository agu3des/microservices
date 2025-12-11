package ports

import "github.com/agu3des/microservices/order/internal/application/core/domain"

type APIPort interface {
	PlaceOrder (order domain.Order) (domain.Order, error)
}