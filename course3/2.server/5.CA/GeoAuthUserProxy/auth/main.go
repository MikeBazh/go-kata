package main

import (
	"context"
	"fmt"
	"github.com/go-chi/chi"
	_ "github.com/joho/godotenv"
	"go-kata/2.server/5.CA/GeoAuthUserProxy/auth/controller"
	"go-kata/2.server/5.CA/GeoAuthUserProxy/auth/responder"
	"go-kata/2.server/5.CA/GeoAuthUserProxy/auth/services"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

func main() {

	// Создание экземпляра сервиса
	servicer := services.NewService()

	// Создание экземпляра контроллера
	responder := responder.NewResponder()
	UserController := controller.NewUserController(responder, servicer)

	r := chi.NewRouter()

	r.Post("/api/auth/register", UserController.Register)

	r.Post("/api/auth/login", UserController.Login)

	//r.Post("/api/login", func(w http.ResponseWriter, r *http.Request) {
	//	conn, err := grpc.Dial(authAddress, grpc.WithInsecure())
	//	if err != nil {
	//		log.Fatalf("Failed to connect: %v", err)
	//	}
	//	defer conn.Close()
	//
	//	// Создаем клиент
	//	UserClient := grpcUser.NewUserServiceClient(conn)
	//	fmt.Println("запрос в сервис User...")
	//	// Вызовите gRPC-сервис для проверки токена
	//	var resp *grpcUser.GetUserByEmailResponse
	//	resp, err = UserClient.GetUserByEmail(context.Background(), &grpcUser.GetUserByEmailRequest{Email: "admin"})
	//	if err != nil {
	//		// Обработайте ошибку
	//		http.Error(w, "Error", http.StatusInternalServerError)
	//	}
	//	fmt.Println("user: resp:", resp)
	//})

	//r.Post("/api/register", func(w http.ResponseWriter, r *http.Request) {
	//	conn, err := grpc.Dial(authAddress, grpc.WithInsecure())
	//	if err != nil {
	//		log.Fatalf("Failed to connect: %v", err)
	//	}
	//	defer conn.Close()
	//
	//	// Создаем клиент
	//	UserClient := grpcUser.NewUserServiceClient(conn)
	//	fmt.Println("запрос в сервис User...")
	//	// Вызовите gRPC-сервис для проверки токена
	//	var resp *grpcUser.CreateUserResponse
	//	resp, err = UserClient.CreateUser(context.Background(), &grpcUser.CreateUserRequest{Email: "email", Password: "password"})
	//	if err != nil {
	//		// Обработайте ошибку
	//		http.Error(w, "Error", http.StatusInternalServerError)
	//	}
	//	fmt.Println("user: resp:", resp)
	//})

	r.Get("/auth/swagger/index.html", UserController.SwaggerUI)
	r.Get("/public/*", func(w http.ResponseWriter, r *http.Request) {
		http.StripPrefix("/public/", http.FileServer(http.Dir("/app/public"))).ServeHTTP(w, r)
		//http.StripPrefix("/public/", http.FileServer(http.Dir("/home/m/GolandProjects/gh_course4/go-kata/course3/2.server/5.CA/GeoAuthUserProxy/auth/public"))).ServeHTTP(w, r)
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

func getTokenFromHeader(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", fmt.Errorf("Authorization header is missing")
	}

	authParts := strings.Split(authHeader, " ")
	if len(authParts) != 2 || strings.ToLower(authParts[0]) != "bearer" {
		return "", fmt.Errorf("Invalid Authorization header format")
	}
	return authParts[1], nil
}
