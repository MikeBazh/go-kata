package Handlers

import (
	"github.com/go-chi/jwtauth"
	"golang.org/x/crypto/bcrypt"
)

//type SearchRequest struct {
//	Query string `json:"query"`
//}
//
//type GeocodeRequest struct {
//	Lat string `json:"lat"`
//	Lng string `json:"lng"`
//}

type RegisterRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}


var users = make(map[string]string)
var TokenAuth *jwtauth.JWTAuth


type Servicer interface {
	RegisterUser(login, password string) error
}

type Service struct {
}

func NewService() Servicer {
	return &Service{}
}

// RegisterUser регистрирует нового пользователя.
func (s *Service) RegisterUser(login, password string) error {
	// Здесь происходит бизнес-логика регистрации пользователя
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	users[login] = string(hashedPassword)
	//fmt.Println("written!!")
	return nil
}


