package order

import (
	"log"
	"net/http"
)

type OrderService struct {
	Repo *OrderRepository
}

func NewOrderService(repo *OrderRepository) *OrderService {
	return &OrderService{
		Repo: repo,
	}
}
func (s *OrderService) getOrderService(orderId int) (Order, int, error) {
	order, err := s.Repo.getOrderDetailsWithId(orderId)
	if err != nil {
		log.Printf("Error fetching order: %v", err)
		return *order, http.StatusInternalServerError, err
	}

	return *order, http.StatusOK, nil
}
