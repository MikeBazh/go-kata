//Ознакомься с материалами по документации swagger и go-swagger.
//Реализуй все хендлеры Swagger Petstore.

package main

import (
	"context"
	"fmt"
	"github.com/go-chi/chi"
	"go-kata/2.server/5.CA/GeoserviceNginx/user/controller"
	"go-kata/2.server/5.CA/GeoserviceNginx/user/responder"
	"go-kata/2.server/5.CA/GeoserviceNginx/user/services"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// Создание экземпляра сервиса
	servicer := services.NewService()

	// Создание экземпляра контроллера
	responder := responder.NewResponder()
	UserController := controller.NewUserController(responder, *servicer)

	r := chi.NewRouter()

	r.Post("/api/user/list", UserController.ListUsers)

	r.Post("/api/user/profile", UserController.Profile)

	r.Get("/user/swagger/index.html", UserController.SwaggerUI)
	r.Get("/public/*", func(w http.ResponseWriter, r *http.Request) {
		http.StripPrefix("/public/", http.FileServer(http.Dir("/app/public"))).ServeHTTP(w, r)
		//http.StripPrefix("/public/", http.FileServer(http.Dir("/home/m/GolandProjects/gh_course4/go-kata/course3/2.server/5.CA/geoservice_user/public"))).ServeHTTP(w, r)

	})

	port := "8080"
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
