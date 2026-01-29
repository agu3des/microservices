package ports

import "github.com/agu3des/microservices/order/internal/application/core/domain"

type DBPort interface {
	Get (id string) (domain.Order, error)
	ProductExists(code string) (bool, error)
	Save (*domain.Order) error
}