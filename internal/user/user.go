package user

type User struct {
	UserID   int    `json:"userId"`
	Email    string `json:"email"`
	Password string `json:"password"`
	IsAdmin  int    `json:"isAdmin"`
}
