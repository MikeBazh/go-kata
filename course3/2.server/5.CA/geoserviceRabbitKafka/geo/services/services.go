package services

import (
	"encoding/json"
	"fmt"
	//_ "github.com/go-chi/jwtauth"
	"github.com/go-redis/redis"
	"go-kata/2.server/5.CA/geoserviceRabbitKafka/geo/Dadata"
	"go-kata/2.server/5.CA/geoserviceRabbitKafka/geo/storage"
)

const (
	authAddress = "localhost:50051"
)

type Servicer interface {
	//RegisterUser(login, password string) error
	//LoginUser(login, password string) (string, error)
	SearchByQuery(Query string) (Dadata.SearchResponse, error, bool)
	SearchByGeo(Lat string, Lng string) (Dadata.GeocodeResponse, error, bool)
	SendRateLimitExceededMessage(email string) error
}

type Service struct {
	UserStorage storage.UserRepository
}

func NewService(UserStorage storage.UserRepository) *Service {
	return &Service{
		UserStorage: UserStorage,
		//UserCache:   UserCache,
	}
}

type ServiceProxy struct {
	Service Servicer
	cache   redis.Client
}

// NewServiceProxy - конструктор хранилища пользователей
func NewServiceProxy(Service Service, cache redis.Client) *ServiceProxy {
	return &ServiceProxy{
		&Service,
		cache,
	}
}

func (s *Service) SearchByGeo(Lat, Lng string) (Dadata.GeocodeResponse, error, bool) {

	newSearchResponse, err, limited := Dadata.RateLimitedAskByGeo(Lat, Lng)
	if limited {
		if err != nil {
			fmt.Println(err)
		}
		return Dadata.GeocodeResponse{}, err, limited
	}
	if err != nil {
		return Dadata.GeocodeResponse{}, err, false
	}
	return newSearchResponse, nil, false
}

func (s *Service) SearchByQuery(Query string) (Dadata.SearchResponse, error, bool) {
	//Проверяем запрос в истории поиска
	//if addressList, err := s.UserStorage.CheckHistory(Query); err == nil {
	//	fmt.Println("service: запрос найден в БД, ответ отправлен из БД")
	//	var responseFromCache Dadata.SearchResponse
	//	var address Dadata.Address
	//	for _, SingleAddress := range addressList {
	//		err = json.Unmarshal(SingleAddress, &address)
	//		if err != nil {
	//			fmt.Println(err)
	//			continue
	//		}
	//		responseFromCache.Addresses = append(responseFromCache.Addresses, &address)
	//	}
	//	return responseFromCache, nil
	//
	//} else {
	// Отправляем апрос во внешний сервис
	newSearchResponse, err, limited := Dadata.RateLimitedAskByQuery(Query)
	if limited {
		if err != nil {
			fmt.Println(err)
		}
		return Dadata.SearchResponse{}, err, limited
	}
	if err != nil {
		return Dadata.SearchResponse{}, err, limited
	}
	//fmt.Println("service: не найдено в БД, отправлен запрос во внешний сервис")
	// Записываем запрос и результаты запроса в базу данных
	SearchHistoryID, err := s.UserStorage.CreateSearchHistory(Query)
	if err != nil {
		fmt.Println(err)
	}
	//Обработка "пустого ответа"
	if newSearchResponse.Addresses == nil {
		jsonAddress, err := json.Marshal(Dadata.Address{})
		RespondHistoryID, err := s.UserStorage.CreateRespondHistory(jsonAddress, " ")
		if err != nil {
			fmt.Println(err)
		}
		err = s.UserStorage.CreateHistorySearch(SearchHistoryID, RespondHistoryID)
		if err != nil {
			fmt.Println(err)
		}
	} else {
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
	}
	//fmt.Println("service: запрос и результаты записаны в БД")
	return newSearchResponse, err, limited
}

//}
