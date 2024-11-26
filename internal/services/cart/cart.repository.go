package cart

import (
	"database/sql"
	"fmt"
)

const (
	CART_ID    = "id"
	TABLE_NAME = "carts"
)

type CartRepository struct {
	db *sql.DB
}

func NewCartRepository(db *sql.DB) *CartRepository {
	return &CartRepository{db: db}
}

func (repo *CartRepository) AddOrUpdateCartItem(cartID, productID, quantity int) error {
	// Upsert query: Insert if the product doesn't exist in the cart, or update the quantity if it does
	query := `
		INSERT INTO cart_items (cart_id, product_id, quantity)
		VALUES (?, ?, ?)
		ON DUPLICATE KEY UPDATE quantity = quantity + VALUES(quantity), updated_at = CURRENT_TIMESTAMP
	`
	_, err := repo.db.Exec(query, cartID, productID, quantity)
	if err != nil {
		return fmt.Errorf("failed to add/update cart item: %v", err)
	}
	return nil
}
