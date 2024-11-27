package cart

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"
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

// ------------CART-ITEM RELATED------------
func (repo *CartRepository) addOrUpdateCartItem(cartID, productID, quantity int) error {
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

// get all products from cart_items JOIN products table
func (repo *CartRepository) getAllCartItem() error {
	return nil
}

// ------------CART RELATED------------
func (repo *CartRepository) getCartByID(cartID int) (*Cart, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	query := `
        SELECT
            c.ID AS cart_id,
            c.UserID,
            c.CreatedAt AS cart_created_at,
            c.UpdatedAt AS cart_updated_at,
            ci.ID AS item_id,
            ci.CartID,
            ci.ProductID,
            ci.Quantity,
            ci.CreatedAt AS item_created_at,
            ci.UpdatedAt AS item_updated_at
        FROM
            carts c
        LEFT JOIN
            cart_items ci ON c.ID = ci.CartID
        WHERE
            c.ID = ?`

	rows, err := repo.db.QueryContext(ctx, query, cartID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cart Cart
	var items []CartItem

	for rows.Next() {
		var cartID, userID, itemID, productID, quantity int
		var cartCreatedAt, cartUpdatedAt, itemCreatedAt, itemUpdatedAt time.Time

		err := rows.Scan(
			&cartID,
			&userID,
			&cartCreatedAt,
			&cartUpdatedAt,
			&itemID,
			&cartID,
			&productID,
			&quantity,
			&itemCreatedAt,
			&itemUpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		if cart.ID == 0 {
			cart.ID = cartID
			cart.UserID = userID
			cart.CreatedAt = cartCreatedAt
			cart.UpdatedAt = cartUpdatedAt
		}

		item := CartItem{
			ID:        itemID,
			CartID:    cartID,
			ProductID: productID,
			Quantity:  quantity,
			CreatedAt: itemCreatedAt,
			UpdatedAt: itemUpdatedAt,
		}
		items = append(items, item)
	}

	if len(items) > 0 {
		cart.Items = items
	} else {
		return nil, fmt.Errorf("no cart found with cartID %d: %v", cartID, err)
	}

	log.Println("Cart data coming from database")
	return &cart, nil
}

func (repo *CartRepository) getCartByUserID(userID int) (*Cart, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	query := `
        SELECT
            c.ID AS cart_id,
            c.UserID,
            c.CreatedAt AS cart_created_at,
            c.UpdatedAt AS cart_updated_at,
            ci.ID AS item_id,
            ci.CartID,
            ci.ProductID,
            ci.Quantity,
            ci.CreatedAt AS item_created_at,
            ci.UpdatedAt AS item_updated_at
        FROM
            carts c
        LEFT JOIN
            cart_items ci ON c.ID = ci.CartID
        WHERE
            c.UserID = ?`

	rows, err := repo.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cart Cart
	var items []CartItem

	for rows.Next() {
		var cartID, userID, itemID, productID, quantity int
		var cartCreatedAt, cartUpdatedAt, itemCreatedAt, itemUpdatedAt time.Time

		err := rows.Scan(
			&cartID,
			&userID,
			&cartCreatedAt,
			&cartUpdatedAt,
			&itemID,
			&cartID,
			&productID,
			&quantity,
			&itemCreatedAt,
			&itemUpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		if cart.ID == 0 {
			cart.ID = cartID
			cart.UserID = userID
			cart.CreatedAt = cartCreatedAt
			cart.UpdatedAt = cartUpdatedAt
		}

		item := CartItem{
			ID:        itemID,
			CartID:    cartID,
			ProductID: productID,
			Quantity:  quantity,
			CreatedAt: itemCreatedAt,
			UpdatedAt: itemUpdatedAt,
		}
		items = append(items, item)
	}
	if len(items) > 0 {
		cart.Items = items
	} else {
		return nil, fmt.Errorf("no cart found with userID %d: %v", userID, err) // No cart found or no items in the cart for the given user
	}

	log.Println("Cart data for user coming from database")
	return &cart, nil
}
