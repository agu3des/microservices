package ports

import "github.com/agu3des/microservices/order/internal/application/core/domain"

type APIPort interface {
	PlaceOrder(customerId int64, orderItems []domain.OrderItem) (domain.Order, error)
}