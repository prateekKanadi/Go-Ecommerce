package order

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/ecommerce/internal/services/cart"
	"github.com/ecommerce/internal/services/user"
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

func (repo *OrderRepository) createOrder(userID int, deliveryMode string, paymentMode string, orderValue float64, orderTotal float64) (int, error) {

	//userRepo := &user.UserRepository{db: repo.db}

	// Get the user's address by their userID
	//address, err := userRepo.GetAddressByUserId(userID)
	address, err := repo.GetAddressAsString(userID)
	//if err != nil {
	//	log.Printf("Error fetching address for user %d: %v", userID, err)
	//}
	//shippingAddress := checkout.CheckoutRepository.getAddressDetailsOfUser(userID)

	query := `
		INSERT INTO orders (userId, createdAt, updatedAt, deliveryMode, paymentMode, orderValue, shippingAddress, orderTotal)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`

	// Execute the query
	result, err := repo.db.Exec(query,
		userID,
		time.Now(),
		time.Now(),
		deliveryMode,
		paymentMode,
		orderValue,
		address,
		orderTotal)
	if err != nil {
		log.Println(err.Error())
		return 0, err
	}
	insertID, err := result.LastInsertId()
	if err != nil {
		log.Println(err.Error())
		return 0, err
	}

	return int(insertID), nil
}

func (repo *OrderRepository) createOrderItems(cartList *cart.Cart, orderId int, orderDetail Order) (Order, error) {
	//var orderItem []OrderItem
	for _, cartItem := range cartList.Items {
		orderItem := OrderItem{
			OrderID:      orderId,
			ProductID:    cartItem.ProductID,
			Quantity:     cartItem.Quantity,
			PricePerUnit: cartItem.PricePerUnit,
			TotalPrice:   cartItem.TotalPrice,
		}
		orderDetail.Items = append(orderDetail.Items, orderItem)

		// Insert the order item into the database
		query := `
		INSERT INTO order_Items (orderId, productId, quantity, priceperunit, totalPrice, createdAt, updatedAt)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`
		// Execute the query
		repo.db.Exec(query,
			orderId,
			orderItem.ProductID,
			orderItem.Quantity,
			orderItem.PricePerUnit,
			orderItem.TotalPrice,
			time.Now(),
			time.Now(),
		)
	}

	return orderDetail, nil
}

func (repo *OrderRepository) GetAddressAsString(userId int) (string, error) {
	// Prepare the query to get the address details by userId
	row := repo.db.QueryRow(`
		SELECT 
			userId, 
			houseNo, 
			landmark, 
			city, 
			state, 
			pincode, 
			phoneNumber 
		FROM address 
		WHERE userId = ?`, userId)

	// Create an Address struct to hold the result
	var address user.Address

	// Scan the row data into the address struct
	err := row.Scan(
		&address.UserID,
		&address.HouseNo,
		&address.Landmark,
		&address.City,
		&address.State,
		&address.Pincode,
		&address.PhoneNumber,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			// If no rows are found for the given userId
			return "", fmt.Errorf("no address found for userId %d", userId)
		}
		log.Println("Error scanning row:", err)
		return "", err
	}

	// Create a formatted string with all the address fields
	addressString := fmt.Sprintf("House No: %s\n, Landmark: %s\n, City: %s\n, State: %s\n, Pincode: %s\n, Phone Number: %s",
		address.HouseNo,
		address.Landmark,
		address.City,
		address.State,
		address.Pincode,
		address.PhoneNumber)

	return addressString, nil
}
