//Реализация кэширования истории поиска
//Реализуй слой repository в геосервисе для метода /api/address/search.
//Чтобы снизить издержки бизнеса по использованию сторонних сервисов, мы решили кэшировать ответы
//от сервиса dadata.ru в базе данных.

package main

import (
	"context"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-redis/redis"
	"github.com/joho/godotenv"
	_ "github.com/joho/godotenv"
	"go-kata/2.server/5.CA/geoserviceRabbitKafka/geo/controller"
	"go-kata/2.server/5.CA/geoserviceRabbitKafka/geo/responder"
	"go-kata/2.server/5.CA/geoserviceRabbitKafka/geo/services"
	"go-kata/2.server/5.CA/geoserviceRabbitKafka/geo/storage"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error load .env file: ", err)
	}
	port := os.Getenv("SERVER_PORT")
	redisDB, _ := strconv.Atoi(os.Getenv("REDIS_DB"))

	client := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),     // Адрес сервера Redis
		Password: os.Getenv("REDIS_PASSWORD"), // Пароль, если установлен
		DB:       redisDB,                     // Индекс базы данных
	})

	MessageProtocol := os.Getenv("MESSAGE_PROTOCOL")

	if !(MessageProtocol == "RabbitMQ" || MessageProtocol == "kafka") {
		log.Fatalf("Server error: Protocol must be set to 'RabbitMQ' or 'kafka'")
	}

	// Создание экземпляра сервиса
	//servicer := services.NewService(MessageProtocol)

	// Создание экземпляра сервиса
	storager := storage.NewUserStorage()
	mesenger := services.NewMessanger(MessageProtocol)
	servicer := services.NewService(storager)
	proxyServicer := services.NewServiceProxy(*servicer, *client)

	// Создание экземпляра контроллера
	responder := responder.NewResponder()
	UserController := controller.NewUserController(responder, servicer, mesenger, *proxyServicer)

	r := chi.NewRouter()

	r.Get("/geo/swagger/index.html", UserController.SwaggerUI)
	r.Get("/public/*", func(w http.ResponseWriter, r *http.Request) {
		http.StripPrefix("/public/", http.FileServer(http.Dir("/app/public"))).ServeHTTP(w, r)
		//http.StripPrefix("/public/", http.FileServer(http.Dir("/home/m/GolandProjects/gh_course4/go-kata/course3/2.server/5.CA/GeoAuthUserProxy/geo/public"))).ServeHTTP(w, r)

	})

	r.Post("/api/address/geocode", UserController.SearchByGeo)
	r.Post("/api/address/search", UserController.SearchByQuery)

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
