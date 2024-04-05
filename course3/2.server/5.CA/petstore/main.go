//В этой задаче требуется добавить контроллер в геосервис, написанный на gohugo.
//Контроллер должен быть реализован в соответствии с принципами чистой архитектуры.

package main

import (
	"context"
	"fmt"
	"github.com/go-chi/chi"
	"go-kata/2.server/5.CA/petstore/storage"
	//"github.com/joho/godotenv"
	"go-kata/2.server/5.CA/petstore/controller"
	"go-kata/2.server/5.CA/petstore/responder"
	"go-kata/2.server/5.CA/petstore/services"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {

	//err := godotenv.Load()
	//if err != nil {
	//	log.Fatal("Error load .env file: ", err)
	//}
	//port := os.Getenv("SERVER_PORT")

	// Создание экземпляра сервиса
	storager := storage.NewLibraryStorage()
	servicer := services.NewService(storager)

	// Создание экземпляра контроллера
	responder := responder.NewResponder()
	UserController := controller.NewUserController(responder, servicer)

	r := chi.NewRouter()

	//r.Post("/api/register", UserController.Register)
	//r.Post("/api/login", UserController.Login)
	r.Get("/swagger/index.html", UserController.SwaggerUI)
	r.Get("/public/*", func(w http.ResponseWriter, r *http.Request) {
		//http.StripPrefix("/public/", http.FileServer(http.Dir("/app/public"))).ServeHTTP(w, r)
		http.StripPrefix("/public/", http.FileServer(http.Dir("/home/m/GolandProjects/GHcourse3/course3/2.server/5.CA/petstore/public"))).ServeHTTP(w, r)
	})

	//r.Group(func(r chi.Router) {
	//r.Use(jwtauth.Verifier(services.TokenAuth))
	//r.Use(jwtauth.Authenticator)
	//r.Use(UserController.UnauthorizedToForbidden)

	r.Post("/api/user", UserController.CreateUser)
	r.Post("/api/user/createWithArray", UserController.CreateWithArray)
	r.Post("/api/user/createWithList", UserController.CreateWithList)
	r.Get("/api/user/login", UserController.LoginUser)
	r.Get("/api/user/logout", UserController.LogoutUser)
	r.Get("/api/user/{username}", UserController.GetUserByName)
	r.Put("/api/user/{username}", UserController.UpdateUserByName)
	r.Delete("/api/user/{username}", UserController.DeleteUserByName)
	//GetUsersWithRentedBooks(w http.ResponseWriter, r *http.Request)
	//r.Delete("/api/users", UserController.GetUsersWithRentedBooks)
	//})

	port := ":8080"
	server := &http.Server{
		//Addr: ":" + port,
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

//Критерии приемки
//Контроллер успешно добавлен в геосервис.
//Контроллер использует интерфейс Responder.
//Контроллер реализован в соответствии с принципами чистой архитектуры.
//Геосервис успешно обрабатывает запросы на /api/address/search и возвращает ожидаемый результат.
//Геосервис успешно обрабатывает запросы на /api/address/geocode и возвращает ожидаемый результат.
//Роуты имеют соответствующую иерархию.
//Присутствует документация swagger для всех эндпоинтов.
//Проект для проверки ментором, просто запускается с помощью команды docker-compose up.
//Покрытие тестами не требуется.
