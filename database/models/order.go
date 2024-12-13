package models

type Order struct {
	ID            int    `json:"id"`
	NumberOfItems int    `json:"number_of_items"`
	CreatedAt     string `json:"created_at"`
	UpdateAt      string `json:"updated_at"`
	Shipping      []OrderShipping
}
