package Handlers

import (
	"encoding/json"
	"fmt"
	"go-kata/2.server/5.server_http_api/task3.2.5.1_JWT/Dadata"
	"net/http"
)

type SearchRequest struct {
	Query string `json:"query"`
}

type GeocodeRequest struct {
	Lat string `json:"lat"`
	Lng string `json:"lng"`
}

func HandlerSearchByQuery(w http.ResponseWriter, r *http.Request) {
	//Проверка формата запроса
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Ошибка, запрос должен быть в формате JSON", http.StatusBadRequest)
		return
	}

	var newSearchRequest SearchRequest
	err := json.NewDecoder(r.Body).Decode(&newSearchRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	newSearchResponse, err := Dadata.AskByQuery(newSearchRequest.Query)
	fmt.Println(newSearchResponse)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Ошибка взаимодействия с внешним сервисом", http.StatusInternalServerError)
		return
	}
	respBody, err := json.Marshal(newSearchResponse)
	if err != nil {
		fmt.Println("Ошибка при кодировании JSON:", err)
		return
	}
	w.Write(respBody)

}

func HandlerSearchByGeo(w http.ResponseWriter, r *http.Request) {
	//Проверка формата запроса
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Ошибка, запрос должен быть в формате JSON", http.StatusBadRequest)
		return
	}

	var newGeocodeRequest GeocodeRequest
	err := json.NewDecoder(r.Body).Decode(&newGeocodeRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	addr, err := Dadata.AskByGeo(newGeocodeRequest.Lat, newGeocodeRequest.Lng)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Ошибка взаимодействия с внешним сервисом", http.StatusInternalServerError)
		return
	}
	fmt.Println(addr)
	respBody, err := json.Marshal(addr)
	if err != nil {
		fmt.Println("Ошибка при кодировании JSON:", err)
		return
	}
	w.Write(respBody)

}
