package api

import (
	"fmt"

	"github.com/agu3des/microservices/order/internal/application/core/domain"
	"github.com/agu3des/microservices/order/internal/ports"
)

type Application struct {
	db       ports.DBPort
	payment  ports.PaymentPort
	shipping ports.ShippingPort
}

func NewApplication(db ports.DBPort, payment ports.PaymentPort, shipping ports.ShippingPort) *Application {
	return &Application{
		db:       db,
		payment:  payment,
		shipping: shipping,
	}
}

func (a Application) PlaceOrder(customerId int64, orderItems []domain.OrderItem) (domain.Order, error) {
	
	// 1. Validação de Estoque 
	for _, item := range orderItems {
		exists, err := a.db.ProductExists(item.ProductCode)
		if err != nil {
			return domain.Order{}, err
		}
		if !exists {
			return domain.Order{}, fmt.Errorf("produto não encontrado no estoque: %s", item.ProductCode)
		}
	}

	// 2. Criação da Entidade de Domínio
	newOrder, err := domain.NewOrder(customerId, orderItems)
	if err != nil {
		return domain.Order{}, err
	}

	// 3. Salva no Banco 
	err = a.db.Save(&newOrder)
	if err != nil {
		return domain.Order{}, err
	}

	// 4. Processa Pagamento 
	err = a.payment.Charge(&newOrder)
	if err != nil {
		return domain.Order{}, err
	}

	// 5. Processa Envio 
	err = a.shipping.Ship(newOrder)
	if err != nil {
		return domain.Order{}, fmt.Errorf("pagamento realizado, mas falha ao solicitar envio: %v", err)
	}

	return newOrder, nil
}