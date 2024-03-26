//В этой задаче требуется добавить контроллер в геосервис, написанный на gohugo.
//Контроллер должен быть реализован в соответствии с принципами чистой архитектуры.

package main

import (
	"context"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/joho/godotenv"
	"go-kata/2.server/5.CA/library/storage"
	//"github.com/joho/godotenv"
	"go-kata/2.server/5.CA/library/controller"
	"go-kata/2.server/5.CA/library/responder"
	"go-kata/2.server/5.CA/library/services"
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

	//r.Post("/api/register", UserController.Register)
	//r.Post("/api/login", UserController.Login)
	r.Get("/swagger/index.html", UserController.SwaggerUI)
	r.Get("/public/*", func(w http.ResponseWriter, r *http.Request) {
		http.StripPrefix("/public/", http.FileServer(http.Dir("/app/public"))).ServeHTTP(w, r)
		//http.StripPrefix("/public/", http.FileServer(http.Dir("/home/m/GolandProjects/GHcourse3/course3/2.server/5.CA/library/public"))).ServeHTTP(w, r)
	})

	r.Group(func(r chi.Router) {
		//r.Use(jwtauth.Verifier(services.TokenAuth))
		//r.Use(jwtauth.Authenticator)
		//r.Use(UserController.UnauthorizedToForbidden)

		r.Get("/api/authors", UserController.GetAuthorsWithBooks)
		r.Post("/api/book/return", UserController.ReturnBook)
		r.Post("/api/book/take", UserController.BookTake)
		r.Post("/api/book/add", UserController.AddBook)
		r.Get("/api/books", UserController.GetBooks)
		r.Post("/api/authors", UserController.AddAuthor)
		//GetUsersWithRentedBooks(w http.ResponseWriter, r *http.Request)
		r.Get("/api/users", UserController.GetUsersWithRentedBooks)
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
