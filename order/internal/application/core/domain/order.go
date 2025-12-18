package domain

import (
	"errors"
	"time"
)

type OrderItem struct {
	ProductCode string  `json:"product_code"`
	UnitPrice   float32 `json:"unit_price"`
	Quantity    int32   `json:"quantity"`
}

type Order struct {
	ID         int64       `json:"id"`
	CustomerID int64       `json:"customer_id"`
	Status     string      `json:"status"`
	OrderItems []OrderItem `json:"order_items"`
	CreatedAt  int64       `json:"created_at"`
}

// NewOrder agora retorna (Order, error) para validar a regra de negócio
func NewOrder(customerId int64, orderItems []OrderItem) (Order, error) {
	// Validação: Quantidade máxima de 50 itens
	var totalQuantity int32
	for _, item := range orderItems {
		totalQuantity += item.Quantity
	}

	if totalQuantity > 50 {
		return Order{}, errors.New("orders cannot have more than 50 items in total")
	}

	return Order{
		CreatedAt:  time.Now().Unix(),
		Status:     "Pending",
		CustomerID: customerId,
		OrderItems: orderItems,
	}, nil
}

func (o *Order) TotalPrice() float32 {
	var totalPrice float32
	for _, orderItem := range o.OrderItems {
		totalPrice += orderItem.UnitPrice * float32(orderItem.Quantity)
	}
	return totalPrice
}