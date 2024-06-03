package services

import (
	"context"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-chi/jwtauth"
	"github.com/go-redis/redis"
	grpcUser "go-kata/2.server/5.CA/GeoAuthUserProxy/auth/dto/user/grpc"
	"google.golang.org/grpc"
	"log"
)

// var users = make(map[string]string)
const (
	userAddress = "grpcUser:50051"
	//userAddress = "localhost:50051"
)

var TokenAuth *jwtauth.JWTAuth = jwtauth.New("HS256", []byte("secret"), nil)

type Servicer interface {
	RegisterUser(login, password string) error
	LoginUser(login, password string) (string, error)
}

type Service struct {
}

func NewService() *Service {
	return &Service{}
}

type ServiceProxy struct {
	Service Servicer
	cache   redis.Client
}

func (s *Service) RegisterUser(login, password string) error {
	conn, err := grpc.Dial(userAddress, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	// Создаем клиент
	UserClient := grpcUser.NewUserServiceClient(conn)
	fmt.Println("запрос в сервис User...")
	// Вызовите gRPC-сервис для проверки токена
	//var resp *grpcUser.CreateUserResponse
	_, err = UserClient.CreateUser(context.Background(), &grpcUser.CreateUserRequest{Email: login, Password: password})
	if err != nil {
		fmt.Println(err)
		return err
	}

	//hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	//if err != nil {
	//	return err
	//}
	//users[login] = string(hashedPassword)
	return nil
}

func (s *Service) LoginUser(login, password string) (string, error) {
	conn, err := grpc.Dial(userAddress, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	// Создаем клиент
	UserClient := grpcUser.NewUserServiceClient(conn)
	fmt.Println("запрос в сервис User...")
	// Вызовите gRPC-сервис для проверки токена
	var resp *grpcUser.GetUserByEmailResponse
	resp, err = UserClient.GetUserByEmail(context.Background(), &grpcUser.GetUserByEmailRequest{Email: login})
	//fmt.Println(err)
	//savedPassword, exists := users[login]
	//err := bcrypt.CompareHashAndPassword([]byte(savedPassword), []byte(password))

	if err != nil {
		fmt.Println(err)
		return "Ошибка: Пользователь не существует или пароль не совпадает", nil
	}

	if password != resp.User.Password {
		return "Ошибка: Пользователь не существует или пароль не совпадает", nil
	}

	// Генерация JWT токена
	_, tokenString, err := TokenAuth.Encode(jwt.MapClaims{"sub": login})
	//fmt.Println("tokenString: ", tokenString)
	if err != nil {
		return "", fmt.Errorf("ошибка генерации JWT токена")
	}
	// Отправка токена в ответе
	response := "Bearer " + tokenString
	return response, nil
}
