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
func (repo *CartRepository) addOrUpdateCartItem(cartID int, variantID string, quantity int, isFormQuantityNotNull bool) error {
	var query string
	if isFormQuantityNotNull {
		query = `INSERT INTO cart_items (cart_id, variantId, quantity)
	VALUES (?, ?, ?)
	ON DUPLICATE KEY UPDATE quantity = VALUES(quantity), updated_at = CURRENT_TIMESTAMP
`
	} else {
		query = `
		INSERT INTO cart_items (cart_id, variantId, quantity)
		VALUES (?, ?, ?)
		ON DUPLICATE KEY UPDATE quantity = quantity + VALUES(quantity), updated_at = CURRENT_TIMESTAMP
	`
	}
	_, err := repo.db.Exec(query, cartID, variantID, quantity)
	if err != nil {
		return fmt.Errorf("failed to add/update cart item: %v", err)
	}
	return nil
}

// --------------------REMOVE-CART-ITEM--------------------
func (repo *CartRepository) removeCartItem(cartID, cartItemID int) error {

	var query string
	query = `DELETE FROM cart_items WHERE cart_id=? AND id=?`
	_, err := repo.db.Exec(query, cartID, cartItemID)
	if err != nil {
		return fmt.Errorf("failed to delete cart item: %v", err)
	}
	return nil
}

// get all products from cart_items JOIN products table
func (repo *CartRepository) GetAllCartItems(cartID int, currency string) (*Cart, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	query := `
        SELECT 
            ci.ID AS item_id,
            ci.cart_id,
            ci.variantId,
            ci.quantity,
            ci.created_at AS item_created_at,
            ci.updated_at AS item_updated_at,
            v.variantName,
            p.productBrand,
            v.description,
			v.stockQuantity,
			v.imageURL,
			v.color,
            pp.price
        FROM 
            cart_items ci
        LEFT JOIN 
            variantProducts v ON ci.variantId = v.variantId
		LEFT JOIN 
            products p ON v.productId = p.productId
        LEFT JOIN 
            productPrices pp ON v.productId = pp.productId AND pp.currencyCode = ?
        WHERE 
            ci.cart_id = ?`

	// Pass currencyCode and cartID to the query
	rows, err := repo.db.QueryContext(ctx, query, currency, cartID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []CartItem
	var cartTotal float64
	var cart Cart

	for rows.Next() {
		var item CartItem
		var price float64

		err := rows.Scan(
			&item.ID,
			&item.CartID,
			&item.VariantID,
			&item.Quantity,
			&item.CreatedAt,
			&item.UpdatedAt,
			&item.ProductName,
			&item.ProductBrand,
			&item.Description,
			&item.StockQuantity,
			&item.ImageURL,
			&item.Color,
			&price,
		)
		if err != nil {
			return nil, err
		}

		item.PricePerUnit = price

		// Calculate the total price for the item based on quantity
		totalPrice := float64(item.Quantity) * price
		item.TotalPrice = totalPrice
		cartTotal += totalPrice

		items = append(items, item)
	}

	cart.Items = items
	cart.CartTotal = cartTotal

	if len(items) == 0 {
		return &cart, nil // No items found for this cart
	}

	log.Println("Cart items with product details fetched from database")
	return &cart, nil
}
