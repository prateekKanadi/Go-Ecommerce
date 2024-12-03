package checkout

import (
	"log"
	"net/http"

	"github.com/ecommerce/internal/services/user"
)

type CheckoutService struct{
	Repo *CheckoutRepository
}


func NewCheckoutService(repo *CheckoutRepository) *CheckoutService {
	return &CheckoutService{
		Repo: repo,
	}
}
func (s *CheckoutService) getAddressDetailsOfUser(userId int) (*user.Address,int,error){
	addressDetails, err := s.Repo.getAddressDetailsOfUser(userId)
	if err != nil {
		log.Printf("Error fetching products: %v", err)
		return nil, http.StatusInternalServerError, err
	}

	return addressDetails, http.StatusOK, nil
}