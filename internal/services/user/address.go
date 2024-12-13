package user

import "time"

type Address struct {
	ID          int       `json:"id"`
	UserID      int       `json:"userId"`
	HouseNo     string    `json:"houseNo"`
	Landmark    string    `json:"landmark"`
	City        string    `json:"city"`
	State       string    `json:"state"`
	Pincode     string    `json:"pincode"`
	PhoneNumber string    `json:"phoneNumber"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}
