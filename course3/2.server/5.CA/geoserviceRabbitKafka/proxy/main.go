package main

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

func main() {
	// Создание прокси-сервера для каждого сервиса
	userProxy := httputil.NewSingleHostReverseProxy(&url.URL{
		Scheme: "http",
		Host:   "user:8080",
	})

	authProxy := httputil.NewSingleHostReverseProxy(&url.URL{
		Scheme: "http",
		Host:   "auth:8080",
	})

	geoProxy := httputil.NewSingleHostReverseProxy(&url.URL{
		Scheme: "http",
		Host:   "geo:8080",
	})

	// Обработчик запросов к документации Swagger для сервиса user
	http.HandleFunc("/user/swagger/", func(w http.ResponseWriter, r *http.Request) {
		r.URL.Path = "/user/swagger/index.html"
		userProxy.ServeHTTP(w, r)
	})

	// Обработчик запросов к документации Swagger для сервиса auth
	http.HandleFunc("/auth/swagger/", func(w http.ResponseWriter, r *http.Request) {
		r.URL.Path = "/auth/swagger/index.html"
		authProxy.ServeHTTP(w, r)
	})

	// Обработчик запросов к документации Swagger для сервиса geo
	http.HandleFunc("/geo/swagger/", func(w http.ResponseWriter, r *http.Request) {
		r.URL.Path = "/geo/swagger/index.html"
		geoProxy.ServeHTTP(w, r)
	})

	// Обработчик запросов к файлу swagger.json
	http.HandleFunc("/public/swagger.json", func(w http.ResponseWriter, r *http.Request) {
		// Прокси на основе предыдущего запроса
		if strings.HasPrefix(r.Referer(), "http://localhost:8080/user/") {
			userProxy.ServeHTTP(w, r)
		} else if strings.HasPrefix(r.Referer(), "http://localhost:8080/auth/") {
			authProxy.ServeHTTP(w, r)
		} else if strings.HasPrefix(r.Referer(), "http://localhost:8080/geo/") {
			geoProxy.ServeHTTP(w, r)
		} else {
			http.Error(w, "Page not found", http.StatusNotFound)
		}
	})

	// Обработчик запросов API
	http.HandleFunc("/api/", func(w http.ResponseWriter, r *http.Request) {
		// Прокси на основе пути запроса
		if strings.HasPrefix(r.URL.Path, "/api/user/") {
			userProxy.ServeHTTP(w, r)
		} else if strings.HasPrefix(r.URL.Path, "/api/auth/") {
			authProxy.ServeHTTP(w, r)
		} else if strings.HasPrefix(r.URL.Path, "/api/address/") {
			geoProxy.ServeHTTP(w, r)
		} else {
			http.Error(w, "Page not found", http.StatusNotFound)
		}
	})

	// Запуск прокси-сервера на порту 8080
	http.ListenAndServe(":8080", nil)
}
