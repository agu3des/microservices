package domain

import (
	"errors"
	"time"
)

type Order struct {
	ID         int64
	CustomerID int64
	Status     string
	OrderItems []OrderItem
	TotalPrice float32 
	DeliveryDays int32
	CreatedAt  int64
}

type OrderItem struct {
	ProductCode string
	UnitPrice   float32
	Quantity    int32
}

func NewOrder(customerId int64, orderItems []OrderItem) (Order, error) {
	if len(orderItems) == 0 {
		return Order{}, errors.New("o pedido deve ter pelo menos um item")
	}

	var totalPrice float32 = 0.0

	totalQuantity := 0
	for _, item := range orderItems {
		if item.Quantity <= 0 {
			return Order{}, errors.New("a quantidade do item deve ser maior que zero")
		}
		if item.UnitPrice < 0 {
			return Order{}, errors.New("o preço do item não pode ser negativo")
		}
		
		totalQuantity += int(item.Quantity)
		
		totalPrice += item.UnitPrice * float32(item.Quantity)
	}

	if totalQuantity > 50 {
		return Order{}, errors.New("pedidos não podem ter mais de 50 itens no total")
	}

	return Order{
		CustomerID: customerId,
		Status:     "Pendente", 
		OrderItems: orderItems,
		TotalPrice: totalPrice, 
		CreatedAt:  time.Now().Unix(),
	}, nil
}