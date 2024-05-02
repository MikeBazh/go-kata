package main

import (
	"context"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/joho/godotenv"
	"go-kata/2.server/5.server_http_api/geoservise_jsonrpc/internal/controller"
	"go-kata/2.server/5.server_http_api/geoservise_jsonrpc/internal/services"
	"go-kata/2.server/5.server_http_api/geoservise_jsonrpc/responder"
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
	protocol := os.Getenv("RPC_PROTOCOL")

	if !(protocol == "json-rpc" || protocol == "rpc") {
		log.Fatalf("Server error: Protocol must be set to 'json-rpc' or 'rpc'")
	}

	// Создание экземпляра сервиса
	servicer := services.NewService(protocol)

	// Создание экземпляра контроллера
	responder := responder.NewResponder()
	UserController := controller.NewUserController(responder, servicer)

	r := chi.NewRouter()

	r.Post("/api/register", UserController.Register)
	r.Post("/api/login", UserController.Login)
	r.Get("/swagger/index.html", UserController.SwaggerUI)
	r.Get("/public/*", func(w http.ResponseWriter, r *http.Request) {
		http.StripPrefix("/public/", http.FileServer(http.Dir("/app/public"))).ServeHTTP(w, r)
	})

	r.Group(func(r chi.Router) {
		//r.Use(jwtauth.Verifier(services.TokenAuth))
		//r.Use(jwtauth.Authenticator)
		r.Use(UserController.UnauthorizedToForbidden)

		r.Post("/api/address/search", UserController.SearchByQuery)
		r.Post("/api/address/geocode", UserController.SearchByGeo)
	})

	port := ":8080"
	server := &http.Server{
		Addr:    port,
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
