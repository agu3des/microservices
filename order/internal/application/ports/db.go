package ports

import "github.com/agu3des/microservices/order/internal/application/core/domain"

type DBPort interface {
	Get (id string) (domain.Order, error)
	Save (*domain.Order) error
}