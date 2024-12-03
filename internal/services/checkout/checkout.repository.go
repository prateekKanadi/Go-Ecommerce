package checkout

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/ecommerce/internal/services/user"
)

type CheckoutRepository struct {
	db *sql.DB
}

const (
	TABLE_NAME = "address"
)

func NewCheckoutRepository(db *sql.DB) *CheckoutRepository {
	return &CheckoutRepository{db: db}
}

func (repo *CheckoutRepository) getAddressDetailsOfUser(userId int) (*user.Address, error) {

	address := &user.Address{}
	// whereClause := fmt.Sprintf("%s = ?", userId)
	// query := utils.BuildSelectQuery(TABLE_NAME, address, whereClause)
    query := "SELECT houseNo, landmark, city, state, pincode, phoneNumber FROM address WHERE userId=?"
	fmt.Println("Query formed is ,",query)

	row := repo.db.QueryRow(query, userId)
	err := row.Scan(
		&address.HouseNo,
		&address.Landmark,
		&address.State,
		&address.City,
		&address.Pincode,
		&address.PhoneNumber)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		log.Println(err)
		return nil, err
	}
	return address, nil
}
