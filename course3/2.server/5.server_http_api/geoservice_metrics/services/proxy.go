package services

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	"go-kata/2.server/5.server_http_api/geoservice_metrics/Dadata"
	"time"
)

func (s *ServiceProxy) SearchByQuery(Query string) (Dadata.SearchResponse, error) {
	startTime := time.Now()
	duration := time.Since(startTime).Seconds()
	defer SearchByQueryCacheDuration.Observe(duration)

	var SearchResponse Dadata.SearchResponse
	// Проверка наличия данных в кэше
	jsonResponse, err := s.cache.Get(Query).Result()
	if err == redis.Nil {
		// Если данных нет в кэше, выполняем запрос к базе данных
		fmt.Println("service-proxy: данных нет в кэше. выполняем запрос к БД...")

		startTimeExclude := time.Now()
		newSearchResponse, err := s.Service.SearchByQuery(Query)

		durationExclude := time.Since(startTimeExclude).Seconds()
		duration = duration - durationExclude

		if err != nil {
			fmt.Println(err)
			return Dadata.SearchResponse{}, err
		}
		jsonResponse, err := json.Marshal(newSearchResponse)
		//fmt.Println("jsonResponse", jsonResponse)
		if err != nil {
			err := fmt.Errorf("service-proxy: ошибка json:", err)
			fmt.Println(err)
			return newSearchResponse, err
		}
		// Сохранение результата запроса в кэше (на 1 мин)
		err = s.cache.Set(Query, jsonResponse, 1*time.Minute).Err()
		if err != nil {
			fmt.Println("service-proxy: ошибка сохранения данных в кэше:", err)
			return Dadata.SearchResponse{}, err
		}
		fmt.Println("service-proxy: результат запроса сохранен в кэше на 1 мин")

		return newSearchResponse, err

	} else if err != nil {
		fmt.Println("service-proxy: ошибка получения данных из кэша:", err)
		newSearchResponse, err := s.SearchByQuery(Query)
		return newSearchResponse, err
	} else {
		// Если данные есть в кэше, выводим их
		err = json.Unmarshal([]byte(jsonResponse), &SearchResponse)
		if err != nil {
			err := fmt.Errorf("service-proxy: ошибка json:", err)
			fmt.Println(err)
			return Dadata.SearchResponse{}, err
		}
		fmt.Println("service-proxy: передан результат запроса из кэша")
	}
	return SearchResponse, nil
}

func (s *ServiceProxy) SearchByGeo(Lat, Lng string) (Dadata.GeocodeResponse, error) {
	startTime := time.Now()
	duration := time.Since(startTime).Seconds()
	defer SearchByGeoCacheDuration.Observe(duration)

	var SearchResponse Dadata.GeocodeResponse
	str := Lat + "-" + Lng
	// Проверка наличия данных в кэше
	jsonResponse, err := s.cache.Get(str).Result()
	if err == redis.Nil {
		// Если данных нет в кэше, выполняем запрос к базе данных
		fmt.Println("service-proxy: данных нет в кэше. выполняем запрос к внеш сервису...")
		startTimeExternal := time.Now()

		newSearchResponse, err := s.Service.SearchByGeo(Lat, Lng)

		durationExternal := time.Since(startTimeExternal).Seconds()
		duration = duration - durationExternal

		if err != nil {
			fmt.Println(err)
			return Dadata.GeocodeResponse{}, err
		}
		jsonResponse, err := json.Marshal(newSearchResponse)
		//fmt.Println("jsonResponse", jsonResponse)
		if err != nil {
			err := fmt.Errorf("service-proxy: ошибка json:", err)
			fmt.Println(err)
			return newSearchResponse, err
		}
		// Сохранение результата запроса в кэше (на 1 мин)
		err = s.cache.Set(str, jsonResponse, 1*time.Minute).Err()
		if err != nil {
			fmt.Println("service-proxy: ошибка сохранения данных в кэше:", err)
			return Dadata.GeocodeResponse{}, err
		}
		fmt.Println("service-proxy: результат запроса сохранен в кэше на 1 мин")

		return newSearchResponse, err

	} else if err != nil {
		fmt.Println("service-proxy: ошибка получения данных из кэша:", err)
		newSearchResponse, err := s.SearchByGeo(Lat, Lng)
		return newSearchResponse, err
	} else {
		// Если данные есть в кэше, выводим их
		err = json.Unmarshal([]byte(jsonResponse), &SearchResponse)
		if err != nil {
			err := fmt.Errorf("service-proxy: ошибка json:", err)
			fmt.Println(err)
			return Dadata.GeocodeResponse{}, err
		}
		fmt.Println("service-proxy: передан результат запроса из кэша")
	}
	return SearchResponse, nil
}