package main

import (
	"context"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-redis/redis"
	"github.com/joho/godotenv"
	_ "github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go-kata/2.server/5.server_http_api/geoservice_metrics/controller"
	"go-kata/2.server/5.server_http_api/geoservice_metrics/responder"
	"go-kata/2.server/5.server_http_api/geoservice_metrics/services"
	"go-kata/2.server/5.server_http_api/geoservice_metrics/storage"
	"log"
	"net/http"
	Pprof "net/http/pprof"
	"os"
	"os/signal"
	"strconv"
	_ "strconv"
	"syscall"
	"time"

	_ "github.com/prometheus/client_golang/prometheus"
	_ "github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error load .env file: ", err)
	}
	port := os.Getenv("SERVER_PORT")
	redisDB, _ := strconv.Atoi(os.Getenv("REDIS_DB"))

	client := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       redisDB,
	})

	//client := redis.NewClient(&redis.Options{
	//	Addr:     "localhost:6379",
	//	Password: "",
	//	DB:       0, // Индекс базы данных
	//})

	// Создание экземпляра сервиса
	storager := storage.NewUserStorage()
	servicer := services.NewService(storager)
	proxyServicer := services.NewServiceProxy(*servicer, *client)

	// Создание экземпляра контроллера
	responder := responder.NewResponder()
	UserController := controller.NewUserController(responder, servicer, *proxyServicer)

	prometheus.MustRegister(controller.RegisterTotal, controller.RegisterDuration, controller.LoginTotal, controller.LoginDuration)
	prometheus.MustRegister(controller.SearchByQueryTotal, controller.SearchByQueryDuration, controller.SearchByGeoTotal, controller.SearchByGeoDuration)
	prometheus.MustRegister(services.SearchByQueryExternalDuration, services.SearchByQueryDatabaseDuration)
	prometheus.MustRegister(services.SearchByQueryCacheDuration, services.SearchByGeoExternalDuration, services.SearchByGeoCacheDuration)

	r := chi.NewRouter()
	//r.Use(middleware.Logger)
	r.Handle("/metrics", promhttp.Handler())

	r.Post("/api/register", UserController.Register)
	r.Post("/api/login", UserController.Login)
	r.Get("/swagger/index.html", UserController.SwaggerUI)
	r.Get("/public/*", func(w http.ResponseWriter, r *http.Request) {
		http.StripPrefix("/public/", http.FileServer(http.Dir("/app/public"))).ServeHTTP(w, r)
		//http.StripPrefix("/public/", http.FileServer(http.Dir("/home/m/GolandProjects/GHcourse3/course3/2.server/5.server_http_api/geoservice_cache/public"))).ServeHTTP(w, r)
	})

	r.Group(func(r chi.Router) {
		//r.Use(jwtauth.Verifier(services.TokenAuth))
		//r.Use(jwtauth.Authenticator)
		r.Use(UserController.UnauthorizedToForbidden)

		r.Post("/api/address/search", UserController.SearchByQuery)
		r.Post("/api/address/geocode", UserController.SearchByGeo)

		r.Get("/mycustompath/pprof/allocs", Pprof.Handler("allocs").ServeHTTP)
		r.Get("/mycustompath/pprof/block", Pprof.Handler("block").ServeHTTP)
		r.Get("/mycustompath/pprof/cmdline", Pprof.Cmdline)
		r.Get("/mycustompath/pprof/goroutine", Pprof.Handler("goroutine").ServeHTTP)
		r.Get("/mycustompath/pprof/heap", Pprof.Handler("heap").ServeHTTP)
		r.Get("/mycustompath/pprof/mutex", Pprof.Handler("mutex").ServeHTTP)
		r.Get("/mycustompath/pprof/profile", Pprof.Profile)
		r.Get("/mycustompath/pprof/threadcreate", Pprof.Handler("threadcreate").ServeHTTP)
		r.Get("/mycustompath/pprof/trace", Pprof.Trace)
	})

	r.Handle("/mycustompath/pprof/", http.HandlerFunc(Pprof.Index))

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

	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
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
