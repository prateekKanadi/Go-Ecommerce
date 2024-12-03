package order

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/ecommerce/utils"
)

type OrderRepository struct {
	db *sql.DB
}

const (
	TABLE_NAME = "orders"
)

func NewOrderRepository(db *sql.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

func (repo *OrderRepository) getOrderDetailsWithId(orderId int) (*Order, error) {

	order := &Order{}
	whereClause := fmt.Sprintf("%s = ?", "orderId")
	query := utils.BuildSelectQuery(TABLE_NAME, order, whereClause)

	row := repo.db.QueryRow(query, orderId)
	err := row.Scan(
		&order.OrderID,
		&order.UserID,
		&order.CreatedAt,
		&order.UpdatedAt,
		&order.DeliveryMode,
		&order.PaymentMode,
		&order.OrderValue,
		&order.ShippingAddress,
		&order.OrderTotal)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		log.Println(err)
		return nil, err
	}
	return order, nil
}
