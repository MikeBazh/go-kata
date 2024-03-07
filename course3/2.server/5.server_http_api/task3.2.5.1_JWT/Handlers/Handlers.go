package Handlers

import (
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-chi/jwtauth"
	"go-kata/2.server/5.server_http_api/task3.2.5.1_JWT/Dadata"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

type SearchRequest struct {
	Query string `json:"query"`
}

type GeocodeRequest struct {
	Lat string `json:"lat"`
	Lng string `json:"lng"`
}

type RegisterRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
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

var users = make(map[string]string)
var TokenAuth *jwtauth.JWTAuth

func HandlerRegister(w http.ResponseWriter, r *http.Request) {
	var NewRequest RegisterRequest
	err := json.NewDecoder(r.Body).Decode(&NewRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	//users[NewRequest.Login] = NewRequest.Password

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(NewRequest.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	users[NewRequest.Login] = string(hashedPassword)

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte("Пользователь зарегистрирован"))
}

func HandlerLogin(w http.ResponseWriter, r *http.Request) {
	var NewRequest RegisterRequest
	err := json.NewDecoder(r.Body).Decode(&NewRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	savedPassword, exists := users[NewRequest.Login]
	err = bcrypt.CompareHashAndPassword([]byte(savedPassword), []byte(NewRequest.Password))
	if !exists || err != nil {
		http.Error(w, "Ошибка: Пользователь не существует или пароль не совпадает", http.StatusOK)
		return
	}

	// Генерация JWT токена
	_, tokenString, err := TokenAuth.Encode(jwt.MapClaims{"sub": NewRequest.Login})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Отправка токена в ответе
	response := "Bearer " + tokenString
	//response["token"] = tokenString
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
	//json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
}

func UnauthorizedToForbidden(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, claims, err := jwtauth.FromContext(r.Context())
		if err != nil || claims == nil {
			w.WriteHeader(http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}
