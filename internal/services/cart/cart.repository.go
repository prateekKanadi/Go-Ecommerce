package cart

import (
	"database/sql"
	"fmt"
	"sync"
)

const (
	CART_ID    = "id"
	TABLE_NAME = "carts"
)

type CartRepository struct {
	db *sql.DB
	mu sync.Mutex // Mutex to prevent concurrent access
}

func NewCartRepository(db *sql.DB) *CartRepository {
	return &CartRepository{db: db}
}

func (r *CartRepository) AddOrUpdateCartItem(cartID, productID, quantity int) error {
	// Upsert query: Insert if the product doesn't exist in the cart, or update the quantity if it does
	query := `
		INSERT INTO cart_items (cart_id, product_id, quantity)
		VALUES (?, ?, ?)
		ON DUPLICATE KEY UPDATE quantity = quantity + VALUES(quantity), updated_at = CURRENT_TIMESTAMP
	`
	_, err := r.db.Exec(query, cartID, productID, quantity)
	if err != nil {
		return fmt.Errorf("failed to add/update cart item: %v", err)
	}
	return nil
}

// func (r *CartRepository) CreateCartForUser(userID int) (int, error) {
// 	// Create a new cart record for the user
// 	query := `INSERT INTO carts (user_id) VALUES (?)`
// 	result, err := r.db.Exec(query, userID)
// 	if err != nil {
// 		return 0, fmt.Errorf("failed to create cart for user %d: %v", userID, err)
// 	}

// 	cartID, err := result.LastInsertId()
// 	if err != nil {
// 		log.Println(err.Error())
// 		return 0, err
// 	}

// 	log.Printf("Cart created for user %d successfully", userID)
// 	return int(cartID), nil
// }
