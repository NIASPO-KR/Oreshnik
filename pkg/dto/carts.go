package dto

type ItemCard struct {
	Item
	Count int `json:"count"`
}

type Cart struct {
	Items      []ItemCard `json:"items"`
	PriceTotal int        `json:"priceTotal"`
}
