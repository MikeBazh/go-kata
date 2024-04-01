package services

import (
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-chi/jwtauth"
	"github.com/go-redis/redis"
	"github.com/prometheus/client_golang/prometheus"
	"go-kata/2.server/5.server_http_api/geoservice_metrics/Dadata"
	"go-kata/2.server/5.server_http_api/geoservice_metrics/storage"
	"golang.org/x/crypto/bcrypt"
	"time"
)

var (
	SearchByQueryExternalDuration = prometheus.NewHistogram(prometheus.HistogramOpts{
		Name:    "myapp_SearchByQueryExternal_duration_seconds",
		Help:    "Request duration in seconds",
		Buckets: prometheus.LinearBuckets(0.1, 0.1, 10),
	})

	SearchByQueryDatabaseDuration = prometheus.NewHistogram(prometheus.HistogramOpts{
		Name:    "myapp_SearchByQueryDatabase_duration_seconds",
		Help:    "Request duration in seconds",
		Buckets: prometheus.LinearBuckets(0.1, 0.1, 10),
	})

	SearchByQueryCacheDuration = prometheus.NewHistogram(prometheus.HistogramOpts{
		Name:    "myapp_SearchByQueryCache_duration_seconds",
		Help:    "Request duration in seconds",
		Buckets: prometheus.LinearBuckets(0.1, 0.1, 10),
	})

	SearchByGeoExternalDuration = prometheus.NewHistogram(prometheus.HistogramOpts{
		Name:    "myapp_SearchByGeoExternal_duration_seconds",
		Help:    "Request duration in seconds",
		Buckets: prometheus.LinearBuckets(0.1, 0.1, 10),
	})

	SearchByGeoCacheDuration = prometheus.NewHistogram(prometheus.HistogramOpts{
		Name:    "myapp_SearchByGeoCache_duration_seconds",
		Help:    "Request duration in seconds",
		Buckets: prometheus.LinearBuckets(0.1, 0.1, 10),
	})
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
	startTime := time.Now()
	duration := time.Since(startTime).Seconds()
	defer SearchByGeoExternalDuration.Observe(duration)

	newSearchResponse, err := Dadata.AskByGeo(Lat, Lng)
	if err != nil {
		return Dadata.GeocodeResponse{}, err
	}
	return newSearchResponse, nil
}

func (s *Service) SearchByQuery(Query string) (Dadata.SearchResponse, error) {
	//Проверяем запрос в истории поиска
	if addressList, err := s.UserStorage.CheckHistory(Query); err == nil {
		fmt.Println("service: запрос найден в БД, ответ отправлен из БД")
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
		startTime := time.Now()
		duration := time.Since(startTime).Seconds()

		newSearchResponse, err := Dadata.AskByQuery(Query)
		SearchByQueryExternalDuration.Observe(duration)

		if err != nil {
			return Dadata.SearchResponse{}, err
		}
		fmt.Println("service: не найдено в БД, отправлен запрос во внешний сервис")
		// Записываем запрос и результаты запроса в базу данных
		startTime2 := time.Now()
		duration2 := time.Since(startTime2).Seconds()
		defer SearchByQueryDatabaseDuration.Observe(duration2)

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
		fmt.Println("service: запрос и результаты записаны в БД")
		return newSearchResponse, err
	}
}
