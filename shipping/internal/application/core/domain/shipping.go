package domain

import "errors"

type Shipping struct {
	OrderID           int64
	Items             []ShippingItem
	EstimatedDelivery int32
}

type ShippingItem struct {
	ProductCode string
	Quantity    int32
}

func NewShipping(orderID int64, items []ShippingItem) (Shipping, error) {
	if len(items) == 0 {
		return Shipping{}, errors.New("a remessa deve conter itens")
	}

	days := CalculateDays(items)

	return Shipping{
		OrderID:           orderID,
		Items:             items,
		EstimatedDelivery: days,
	}, nil
}

func CalculateDays(items []ShippingItem) int32 {
	totalQuantity := 0
	for _, item := range items {
		if item.Quantity > 0 {
			totalQuantity += int(item.Quantity)
		}
	}
	days := 1
	extraDays := totalQuantity / 5
	
	return int32(days + extraDays)
}