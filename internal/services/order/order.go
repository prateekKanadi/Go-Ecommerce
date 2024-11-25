package order

import "time"

type Order struct {
	ID          int         `json:"id"`
	UserID      int         `json:"user_id"`
	TotalAmount float64     `json:"total_amount"`
	Status      string      `json:"status"` // 'pending', 'completed', 'cancelled'
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
	Items       []OrderItem // One-to-many relationship
}

type OrderItem struct {
	ID           int       `json:"id"`
	OrderID      int       `json:"order_id"`
	ProductID    int       `json:"product_id"`
	Quantity     int       `json:"quantity"`
	PricePerUnit float64   `json:"price_per_unit"`
	TotalPrice   float64   `json:"total_price"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
