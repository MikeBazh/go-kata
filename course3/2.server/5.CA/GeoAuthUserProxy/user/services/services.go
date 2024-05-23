package services

import (
	"context"
	"fmt"
	"go-kata/2.server/5.CA/GeoAuthUserProxy/user/dto"
	grpcUser "go-kata/2.server/5.CA/GeoAuthUserProxy/user/dto/user"
	"google.golang.org/grpc"
)

const (
	dbUserAddress = "grpcUser:50051"
)

type Servicer interface {
	ListUsers() ([]dto.User, error)
}

type Service struct {
}

func NewService() *Service {
	return &Service{}
}

func (s *Service) ListUsers() ([]dto.User, error) {

	conn, err := grpc.Dial(dbUserAddress, grpc.WithInsecure())
	// Создаем клиент
	UserClient := grpcUser.NewUserServiceClient(conn)
	fmt.Println("запрос в сервис User...")
	// Вызовите gRPC-сервис для проверки токена
	var resp *grpcUser.ListUsersResponse
	resp, err = UserClient.ListUsers(context.Background(), &grpcUser.ListUsersRequest{Limit: 100, Offset: 0})
	if err != nil {
		fmt.Println("Error:", err)
		return []dto.User{}, err
	}
	//fmt.Println("user: resp:", resp)
	var users []dto.User
	for _, u := range resp.Users {
		users = append(users, dto.User{
			ID:       int(u.Id),
			Email:    u.Email,
			Password: u.Password,
		})
	}
	return users, nil
}

func (s *Service) Profile(email string) (dto.User, error) {

	conn, err := grpc.Dial(dbUserAddress, grpc.WithInsecure())
	// Создаем клиент
	UserClient := grpcUser.NewUserServiceClient(conn)
	fmt.Println("запрос в сервис User...")
	// Вызовите gRPC-сервис для проверки токена
	var resp *grpcUser.GetUserByEmailResponse
	resp, err = UserClient.GetUserByEmail(context.Background(), &grpcUser.GetUserByEmailRequest{Email: email})
	if err != nil {
		fmt.Println("Error:", err)
		return dto.User{}, err
	}
	//fmt.Println("user: resp:", resp)
	user := dto.User{
		ID:       int(resp.User.Id),
		Email:    resp.User.Email,
		Password: resp.User.Password,
	}
	return user, nil
}
