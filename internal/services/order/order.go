package order

import "time"

type Order struct {
	OrderID         int         `json:"orderId"`
	UserID          int         `json:"userId"`
	CreatedAt       time.Time   `json:"createdAt"`
	UpdatedAt       time.Time   `json:"updatedAt"`
	DeliveryMode    string      `json:"deliveryMode"`
	PaymentMode     string      `json:"paymentMode"`
	OrderValue      float64     `json:"orderValue"`
	ShippingAddress string      `json:"shippingAddress"`
	OrderTotal      float64     `json:"orderTotal"`
	Items           []OrderItem // One-to-many relationship
}

type OrderItem struct {
	OrderItemID  int       `json:"orderItemId"`
	OrderID      int       `json:"orderId"`
	ProductID    int       `json:"productId"`
	Quantity     int       `json:"quantity"`
	PricePerUnit float64   `json:"priceperunit"`
	TotalPrice   float64   `json:"totalPrice"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}
