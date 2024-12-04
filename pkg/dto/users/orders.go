package users

type Order struct {
	Items      []ItemCount `json:"items"`
	ID         int         `json:"id"`
	PostomatID string      `json:"postomatID"`
	PaymentID  string      `json:"paymentID"`
	Status     string      `json:"status"`
}
