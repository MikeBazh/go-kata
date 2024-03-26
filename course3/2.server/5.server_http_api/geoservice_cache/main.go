//Реализация кэширования истории поиска
//Реализуй слой repository в геосервисе для метода /api/address/search.
//Чтобы снизить издержки бизнеса по использованию сторонних сервисов, мы решили кэшировать ответы
//от сервиса dadata.ru в базе данных.

package main

import (
	"context"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/jwtauth"
	"github.com/joho/godotenv"
	"go-kata/2.server/5.server_http_api/geoservice_cache/controller"
	"go-kata/2.server/5.server_http_api/geoservice_cache/responder"
	"go-kata/2.server/5.server_http_api/geoservice_cache/services"
	"go-kata/2.server/5.server_http_api/geoservice_cache/storage"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error load .env file: ", err)
	}
	port := os.Getenv("SERVER_PORT")

	// Создание экземпляра сервиса
	storager := storage.NewUserStorage()
	servicer := services.NewService(storager)

	// Создание экземпляра контроллера
	responder := responder.NewResponder()
	UserController := controller.NewUserController(responder, servicer)

	r := chi.NewRouter()

	r.Post("/api/register", UserController.Register)
	r.Post("/api/login", UserController.Login)
	r.Get("/swagger/index.html", UserController.SwaggerUI)
	r.Get("/public/*", func(w http.ResponseWriter, r *http.Request) {
		http.StripPrefix("/public/", http.FileServer(http.Dir("/app/public"))).ServeHTTP(w, r)
		//http.StripPrefix("/public/", http.FileServer(http.Dir("/home/m/GolandProjects/GHcourse3/course3/2.server/5.server_http_api/geoservice_cache/public"))).ServeHTTP(w, r)

	})

	r.Group(func(r chi.Router) {
		r.Use(jwtauth.Verifier(services.TokenAuth))
		//r.Use(jwtauth.Authenticator)
		r.Use(UserController.UnauthorizedToForbidden)

		r.Post("/api/address/search", UserController.SearchByQuery)
		r.Post("/api/address/geocode", UserController.SearchByGeo)
	})

	//port := ":8080"
	server := &http.Server{
		Addr: ":" + port,
		//Addr:    port,
		Handler: r}

	// Создание канала для получения сигналов остановки
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT)

	// Запуск сервера в отдельной горутине
	go func() {
		log.Printf("server started on port %s ", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Обработка сигнала SIGINT
	sig := <-sigChan
	fmt.Printf("Received signal: %s\n", sig)

	// Создание контекста с таймаутом для graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// Остановка сервера с использованием graceful shutdown
	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	log.Println("Server stopped gracefully")
}
