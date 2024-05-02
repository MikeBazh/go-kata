package services

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-chi/jwtauth"
	"go-kata/2.server/5.server_http_api/geoservise_jsonrpc/dto"
	"golang.org/x/crypto/bcrypt"
)

var users = make(map[string]string)

var TokenAuth *jwtauth.JWTAuth = jwtauth.New("HS256", []byte("secret"), nil)

type Servicer interface {
	RegisterUser(login, password string) error
	LoginUser(login, password string) (string, error)
	SearchByQuery(Query string) (dto.SearchResponse, error)
	SearchByGeo(Lat, Lng string) (dto.GeocodeResponse, error)
}

type Service struct {
	provider GeoProvider
}

func NewService(protocol string) Servicer {
	return &Service{provider: NewGeoProvider(protocol)}
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

func (s *Service) SearchByQuery(Query string) (dto.SearchResponse, error) {
	//g := NewGeoProvider()
	ProviderReply, err := s.provider.AddressSearch(Query)
	var reply dto.SearchResponse
	reply.Addresses = ProviderReply
	if err != nil {
		fmt.Println("Ошибка при вызове удаленного метода:", err)
		return reply, err
	}
	return reply, nil
}

func (s *Service) SearchByGeo(Lat, Lng string) (dto.GeocodeResponse, error) {
	//g := NewGeoProvider()
	ProviderReply, err := s.provider.GeoCode(Lat, Lng)
	var reply dto.GeocodeResponse
	reply.Addresses = ProviderReply
	if err != nil {
		fmt.Println("Ошибка при вызове удаленного метода:", err)
		return reply, err
	}
	return reply, nil
}
