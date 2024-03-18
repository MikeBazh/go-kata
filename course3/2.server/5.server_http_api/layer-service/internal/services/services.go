package services

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-chi/jwtauth"
	"go-kata/2.server/5.server_http_api/layer-service/Dadata"
	"golang.org/x/crypto/bcrypt"
)

var users = make(map[string]string)

var TokenAuth *jwtauth.JWTAuth = jwtauth.New("HS256", []byte("secret"), nil)

type Servicer interface {
	RegisterUser(login, password string) error
	LoginUser(login, password string) (string, error)
	SearchByQuery(Query string) (Dadata.SearchResponse, error)
	SearchByGeo(Lat, Lng string) (Dadata.GeocodeResponse, error)
}

type Service struct {
}

func NewService() Servicer {
	return &Service{}
}

// RegisterUser регистрирует нового пользователя.
func (s *Service) RegisterUser(login, password string) error {

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	users[login] = string(hashedPassword)
	return nil
}

func (s *Service) LoginUser(login, password string) (string, error) {

	savedPassword, exists := users[login]
	err := bcrypt.CompareHashAndPassword([]byte(savedPassword), []byte(password))
	if !exists || err != nil {
		return "Ошибка: Пользователь не существует или пароль не совпадает", nil
	}
	// Генерация JWT токена
	_, tokenString, err := TokenAuth.Encode(jwt.MapClaims{"sub": login})
	if err != nil {
		return "", fmt.Errorf("ошибка генерации JWT токена")
	}
	// Отправка токена в ответе
	response := "Bearer " + tokenString
	return response, nil
}

func (s *Service) SearchByQuery(Query string) (Dadata.SearchResponse, error) {

	newSearchResponse, err := Dadata.AskByQuery(Query)
	if err != nil {
		return Dadata.SearchResponse{}, err
	}
	return newSearchResponse, nil
}

func (s *Service) SearchByGeo(Lat, Lng string) (Dadata.GeocodeResponse, error) {

	newSearchResponse, err := Dadata.AskByGeo(Lat, Lng)
	if err != nil {
		return Dadata.GeocodeResponse{}, err
	}
	return newSearchResponse, nil
}
