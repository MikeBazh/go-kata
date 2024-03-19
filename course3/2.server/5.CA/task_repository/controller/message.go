package controller

type RegisterRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type RegisterResponse struct {
	Success   bool `json:"success"`
	ErrorCode int  `json:"error_code,omitempty"`
	Data      Data `json:"data"`
}

type Data struct {
	Message string `json:"message"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SearchRequest struct {
	Query string `json:"query"`
}

type GeocodeRequest struct {
	Lat string `json:"lat"`
	Lng string `json:"lng"`
}

//type User struct {
//	ID    int    `json:"id"`
//	Name  string `json:"name"`
//	Email string `json:"email"`
//	//Verified      bool   `json:"verified"`
//}

//type RequestUser struct {
//	Name  string `json:"name"`
//	Email string `json:"email"`
//}
