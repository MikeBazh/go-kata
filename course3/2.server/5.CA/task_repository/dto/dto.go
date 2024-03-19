package dto

type User struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Email   string `json:"email"`
	Deleted bool   `json:"deleted"`
}

type RequestUser struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}
