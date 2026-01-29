package api

import (
	"github.com/agu3des/microservices/shipping/internal/application/core/domain"
)

type Application struct {
}

func NewApplication() *Application {
	return &Application{}
}

func (a Application) ShipOrder(orderId int64, items []domain.ShippingItem) (int32, error) {
	days := domain.CalculateDays(items)
	return days, nil
}