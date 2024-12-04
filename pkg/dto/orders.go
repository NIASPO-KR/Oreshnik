package dto

type Order struct {
	Items      []ItemCard  `json:"items"`
	Postomat   PickupPoint `json:"postomat"`
	Payment    Payment     `json:"payment"`
	ID         int         `json:"id"`
	PriceTotal int         `json:"priceTotal"`
	Status     string      `json:"status"`
}
