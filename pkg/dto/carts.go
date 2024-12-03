package dto

type ItemCart struct {
	Item
	Count int `json:"count"`
}

type Cart struct {
	Items      []ItemCart `json:"items"`
	PriceTotal int        `json:"priceTotal"`
}
