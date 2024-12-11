package order

import (
	"log"
	"net/http"

	"github.com/ecommerce/internal/services/cart"
	"github.com/ecommerce/internal/services/user"
)

type OrderService struct {
	Repo        *OrderRepository
	UserService *user.UserService
	CartService *cart.CartService
}

func NewOrderService(userService *user.UserService, cartService *cart.CartService, repo *OrderRepository) *OrderService {
	return &OrderService{
		UserService: userService,
		CartService: cartService,
		Repo:        repo,
	}
}

func (s *OrderService) createOrderService(userID int, deliveryMode string, paymentMode string, orderValue float64, orderTotal float64, shippingAddress string) (int, int, error) {
	orderId, err := s.Repo.createOrder(userID, deliveryMode, paymentMode, orderValue, orderTotal,shippingAddress)
	if err != nil {
		log.Printf("Error fetching order: %v", err)
		return orderId, http.StatusInternalServerError, err
	}

	return orderId, http.StatusOK, nil
}
func (s *OrderService) createOrderItemsService(orderId int, userId int, orderDetail Order) (Order, int, error) {
	cartID, err := s.UserService.Repo.GetCartForUser(userId)

	cartList, err := s.CartService.Repo.GetAllCartItems(cartID)

	orderDetails, err := s.Repo.createOrderItems(cartList, orderId, orderDetail)
	if err != nil {
		log.Printf("Error fetching order: %v", err)
		return orderDetails, http.StatusInternalServerError, err
	}

	return orderDetails, http.StatusOK, nil
}

func (s *OrderService) getOrderService(orderId int) (Order, int, error) {
	order, err := s.Repo.getOrderDetailsWithId(orderId)
	if err != nil {
		log.Printf("Error fetching order: %v", err)
		return *order, http.StatusInternalServerError, err
	}

	return *order, http.StatusOK, nil
}

// func (s *OrderService) getOrderItemDetails(orderId int) ([]OrderItem, int, error){

// 	orderItems, err := s.Repo.getOrderItemDetails(orderId)
// 	if err != nil {
// 		log.Printf("Error fetching order: %v", err)
// 		return orderItems, http.StatusInternalServerError, err
// 	}

// 	return orderItems, http.StatusOK, nil
// }
