package ports

import "github.com/agu3des/microservices/order/internal/application/core/domain"

type PaymentPort interface {
	Charge(order *domain.Order) error
}