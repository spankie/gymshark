package models

type ShippingPack struct {
	ID        int
	Quantity  int
	CreatedAt string
	UpdateAt  string
}

type OrderShipping struct {
	ID                   int
	OrderID              int
	PackSize             int
	ShippingPackQuantity int
	CreatedAt            string
	UpdatedAt            string
}
