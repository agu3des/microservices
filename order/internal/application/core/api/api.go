package api

import (
	"github.com/agu3des/microservices/order/internal/application/core/domain"
	"github.com/agu3des/microservices/order/internal/ports"
)

type Application struct {
	db      ports.DBPort
	payment ports.PaymentPort
}

func NewApplication(db ports.DBPort, payment ports.PaymentPort) *Application {
	return &Application{
		db:      db,
		payment: payment,
	}
}

func (a *Application) PlaceOrder(customerId int64, orderItems []domain.OrderItem) (domain.Order, error) {
	order, err := domain.NewOrder(customerId, orderItems)
	if err != nil {
		return domain.Order{}, err
	}

	err = a.db.Save(&order)
	if err != nil {
		return domain.Order{}, err
	}

	paymentErr := a.payment.Charge(&order)
	if paymentErr != nil {
		return domain.Order{}, paymentErr
	}

	return order, nil
}