//Ознакомься с материалами по документации swagger и go-swagger.
//Реализуй все хендлеры Swagger Petstore.

package main

import (
	"context"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/jwtauth"
	"github.com/joho/godotenv"
	"go-kata/2.server/5.CA/petstore/controller"
	"go-kata/2.server/5.CA/petstore/responder"
	"go-kata/2.server/5.CA/petstore/services"
	"go-kata/2.server/5.CA/petstore/storage"
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
	storager := storage.NewLibraryStorage()
	servicer := services.NewService(storager)

	// Создание экземпляра контроллера
	responder := responder.NewResponder()
	UserController := controller.NewUserController(responder, servicer)

	r := chi.NewRouter()

	r.Get("/swagger/index.html", UserController.SwaggerUI)
	r.Get("/public/*", func(w http.ResponseWriter, r *http.Request) {
		http.StripPrefix("/public/", http.FileServer(http.Dir("/app/public"))).ServeHTTP(w, r)
		//http.StripPrefix("/public/", http.FileServer(http.Dir("/home/m/GolandProjects/GHcourse3/course3/2.server/5.CA/petstore/public"))).ServeHTTP(w, r)
	})

	r.Group(func(r chi.Router) {
		r.Use(jwtauth.Verifier(services.TokenAuth))
		//r.Use(jwtauth.Authenticator)
		//r.Use(UserController.UnauthorizedToForbidden)
		r.Delete("/api/user/{username}", UserController.DeleteUserByName)

		r.Get("/api/pet/findByStatus", UserController.FindPetByStatus)
		r.Post("/api/pet", UserController.AddPet)
		r.Put("/api/pet", UserController.UpdatePet)
		r.Delete("/api/pet/{petId}", UserController.DeletePet)
		r.Post("/api/pet/{petId}", UserController.UpdatePetWithData)
		r.Get("/api/pet/{petId}", UserController.FindPetById)

		r.Get("/api/store/inventory", UserController.Inventory)
	})

	r.Post("/api/user", UserController.CreateUser)
	r.Post("/api/user/createWithArray", UserController.CreateWithArray)
	r.Post("/api/user/createWithList", UserController.CreateWithList)
	r.Get("/api/user/login", UserController.LoginUser)
	r.Get("/api/user/logout", UserController.LogoutUser)
	r.Get("/api/user/{username}", UserController.GetUserByName)
	r.Put("/api/user/{username}", UserController.UpdateUserByName)

	r.Post("/api/store/order", UserController.AddOrder)
	r.Get("/api/store/order/{orderId}", UserController.FindOrderById)
	r.Delete("/api/store/order/{orderId}", UserController.DeleteOrder)

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
//Все хендлеры реализованы с применением всех 3 слоев: repository, service, handler (Controller).
//Присутствует документация Swagger.
//Документация корректно работает в Swagger.
//Защищенные роуты доступны только при наличии токена.
//Загрузку файла реализовывать не требуется
//Нужно создать документацию только на 200 коды ответа
//Будет плюсом, если ты реализуешь для каждой сущности свой service.
//Функционал покрыт тестами на 90%.
