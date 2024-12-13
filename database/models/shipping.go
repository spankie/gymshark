package models

type ShippingPack struct {
	ID        int    `json:"id"`
	Quantity  int    `json:"quantity"`
	CreatedAt string `json:"created_at"`
	UpdateAt  string `json:"updated_at"`
}

type OrderShipping struct {
	ID                   int    `json:"id"`
	OrderID              int    `json:"order_id"`
	PackSize             int    `json:"pack_size"`
	ShippingPackQuantity int    `json:"shipping_pack_quantity"`
	CreatedAt            string `json:"created_at"`
	UpdateAt             string `json:"updated_at"`
}
