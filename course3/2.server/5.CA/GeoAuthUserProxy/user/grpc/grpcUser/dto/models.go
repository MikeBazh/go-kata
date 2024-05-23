package dto

type User struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RequestUser struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
