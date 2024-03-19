package services

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-chi/jwtauth"
	"go-kata/2.server/5.CA/task_repository/Dadata"
	"go-kata/2.server/5.CA/task_repository/dto"
	"go-kata/2.server/5.CA/task_repository/storage"
	"golang.org/x/crypto/bcrypt"
)

//type User struct {
//	ID    int    `json:"id"`
//	Name  string `json:"name"`
//	Email string `json:"email"`
//	//Verified      bool   `json:"verified"`
//}

var users = make(map[string]string)

var TokenAuth *jwtauth.JWTAuth = jwtauth.New("HS256", []byte("secret"), nil)

type Servicer interface {
	RegisterUser(login, password string) error
	LoginUser(login, password string) (string, error)
	SearchByQuery(Query string) (Dadata.SearchResponse, error)
	SearchByGeo(Lat, Lng string) (Dadata.GeocodeResponse, error)

	DbAddUser(user dto.RequestUser) error
	DbGetUserByID(id int) (dto.User, error)
	UpdateUser(user dto.RequestUser) error
	DeleteByID(id int) error
	List(limit, offset int) ([]dto.User, error)
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

func (s *Service) DbAddUser(user dto.RequestUser) error {
	store := storage.NewUserStorage()
	err := store.Create(user)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) DbGetUserByID(id int) (dto.User, error) {
	store := storage.NewUserStorage()
	user, err := store.GetByID(id)
	if err != nil {
		fmt.Println(err)
	}
	return user, err
}

func (s *Service) UpdateUser(user dto.RequestUser) error {
	store := storage.NewUserStorage()
	err := store.Update(user)
	//if err != nil {
	//	return user, err
	//}
	return err
}

func (s *Service) DeleteByID(id int) error {
	store := storage.NewUserStorage()
	user, err := store.GetByID(id)
	//if err != nil {
	//	return err
	//}
	if err != nil || user.Deleted {
		err = errors.New("пользователь уже удален или не существует")
		fmt.Println(err)
		return err
	}
	err = store.Delete(id)
	return err
}

func (s *Service) List(limit, offset int) ([]dto.User, error) {
	store := storage.NewUserStorage()
	//users := []dto.User{}
	UsersReceived, err := store.List(limit, offset)
	//if err != nil {
	//	return user, err
	//}
	return UsersReceived, err
}
