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
func (repo *CartRepository) getAllCartItems(cartID int) ([]CartItem, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	query := `
        SELECT 
            ci.ID AS item_id,
            ci.cart_id,
            p.productId,
            ci.quantity,
            ci.created_at AS item_created_at,
            ci.updated_at AS item_updated_at,
            p.productName,
            p.productBrand,
            p.description,
            p.pricePerUnit
        FROM 
            cart_items ci
        LEFT JOIN 
            products p ON ci.product_id = p.productId
        WHERE 
            ci.cart_id = ?`

	rows, err := repo.db.QueryContext(ctx, query, cartID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []CartItem
	var cartTotal float64
	for rows.Next() {
		var item CartItem
		var productName, productBrand, description string
		var pricePerUnit float64
		err := rows.Scan(
			&item.ID,
			&item.CartID,
			&item.ProductID,
			&item.Quantity,
			&item.CreatedAt,
			&item.UpdatedAt,
			&productName,
			&productBrand,
			&description,
			&pricePerUnit,
		)
		if err != nil {
			return nil, err
		}

		// Attach product details to CartItem struct directly without modifying CartItem struct
		// You can choose to use these values for other operations, but don't store them in CartItem database.
		// item.ProductID =
		item.ProductName = productName
		item.ProductBrand = productBrand
		item.Description = description
		item.PricePerUnit = pricePerUnit
		totalPrice := ((float64)(item.Quantity) * pricePerUnit)
		item.TotalPrice = float64(totalPrice)
		cartTotal += totalPrice

		// Append the item to the items slice
		items = append(items, item)
	}

	if len(items) == 0 {
		return nil, nil // No items found for this cart
	}

	log.Println("Cart items with product details fetched from database")
	return items, nil
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
