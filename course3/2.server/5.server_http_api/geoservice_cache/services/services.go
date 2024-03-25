package services

import (
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-chi/jwtauth"
	"go-kata/2.server/5.server_http_api/geoservice_cache/Dadata"
	"go-kata/2.server/5.server_http_api/geoservice_cache/storage"
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
	UserStorage storage.UserRepository
}

func NewService(UserStorage storage.UserRepository) *Service {
	return &Service{
		UserStorage: UserStorage}
}

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

func (s *Service) SearchByGeo(Lat, Lng string) (Dadata.GeocodeResponse, error) {

	newSearchResponse, err := Dadata.AskByGeo(Lat, Lng)
	if err != nil {
		return Dadata.GeocodeResponse{}, err
	}
	return newSearchResponse, nil
}

func (s *Service) SearchByQuery(Query string) (Dadata.SearchResponse, error) {
	//Проверяем запрос в истории поиска
	if addressList, err := s.UserStorage.CheckHistory(Query); err == nil {
		fmt.Println("service: запрос найден в истории, ответ отправлен из кэша")
		var responseFromCache Dadata.SearchResponse
		var address Dadata.Address
		for _, SingleAddress := range addressList {
			err = json.Unmarshal(SingleAddress, &address)
			if err != nil {
				fmt.Println(err)
				continue
			}
			responseFromCache.Addresses = append(responseFromCache.Addresses, &address)
		}
		return responseFromCache, nil

	} else {
		// Отправляем апрос во внешний сервис
		newSearchResponse, err := Dadata.AskByQuery(Query)
		if err != nil {
			return Dadata.SearchResponse{}, err
		}
		fmt.Println("service: не найдено в истории, отправлен запрос во внешний сервис")
		// Записываем запрос и результаты запроса в базу данных
		SearchHistoryID, err := s.UserStorage.CreateSearchHistory(Query)
		if err != nil {
			fmt.Println(err)
		}
		for _, address := range newSearchResponse.Addresses {
			jsonAddress, err := json.Marshal(address)
			RespondHistoryID, err := s.UserStorage.CreateRespondHistory(jsonAddress, address.UnrestrictedValue)
			if err != nil {
				fmt.Println(err)
			}
			err = s.UserStorage.CreateHistorySearch(SearchHistoryID, RespondHistoryID)
			if err != nil {
				fmt.Println(err)
			}
		}
		fmt.Println("service: запрос и результаты записаны в базу данных")
		return newSearchResponse, err
	}
}
