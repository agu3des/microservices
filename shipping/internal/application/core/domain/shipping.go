package domain

type ShippingItem struct {
	ProductCode string
	Quantity    int32
}

func CalculateDays(items []ShippingItem) int32 {
	totalQty := int32(0)
	for _, item := range items {
		totalQty += item.Quantity
	}
	// Lógica: 1 dia base + (total / 5)
	// Em Go, divisão de inteiros já arredonda para baixo (floor), que é o comportamento desejado aqui para somar dias completos
	return 1 + (totalQty / 5)
}