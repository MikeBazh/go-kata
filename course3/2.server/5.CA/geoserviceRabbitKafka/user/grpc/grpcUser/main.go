package main

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq" // Импортируем драйвер PostgreSQL
	"go-kata/2.server/5.CA/GeoAuthUserProxy/user/grpc/grpcUser/dto"
	grpcUser "go-kata/2.server/5.CA/GeoAuthUserProxy/user/grpc/grpcUser/grpc"
	"google.golang.org/grpc"
	"log"
	"net"
)

const (
	connStr = "host=dbUser user=postgres password=123 dbname=postgres sslmode=disable"
	//connStr = "host=172.17.0.2 user=postgres password=123/ dbname=postgres sslmode=disable"
)

type UserStorage struct {
	grpcUser.UnimplementedUserServiceServer
}

func (s *UserStorage) CreateUser(ctx context.Context, req *grpcUser.CreateUserRequest) (*grpcUser.CreateUserResponse, error) {
	err := CreateTable()
	if err != nil {
		return nil, err
	}

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer db.Close()

	_, err = db.Exec(`INSERT INTO users (email, password) VALUES ($1, $2)`, req.Email, req.Password)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &grpcUser.CreateUserResponse{Message: "User created successfully"}, nil
}

func (s *UserStorage) GetUserByEmail(ctx context.Context, req *grpcUser.GetUserByEmailRequest) (*grpcUser.GetUserByEmailResponse, error) {
	err := CreateTable()
	if err != nil {
		return nil, err
	}

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	defer db.Close()
	resp := &grpcUser.GetUserByEmailResponse{}
	var user dto.User
	row := db.QueryRow(`SELECT * FROM users WHERE email=$1`, req.Email)
	err = row.Scan(&user.ID, &user.Email, &user.Password)
	fmt.Println(err)
	if err != nil {
		return resp, err
	}
	//resp:= &grpcUser.GetUserByEmailResponse{}
	resp = &grpcUser.GetUserByEmailResponse{User: &grpcUser.User{
		Id:       int32(user.ID),
		Email:    user.Email,
		Password: user.Password,
	}}
	return resp, nil
}

func (s *UserStorage) ListUsers(ctx context.Context, req *grpcUser.ListUsersRequest) (*grpcUser.ListUsersResponse, error) {
	err := CreateTable()
	if err != nil {
		return nil, err
	}

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	query := fmt.Sprintf(`SELECT * FROM users ORDER BY id LIMIT %d OFFSET %d`, req.Limit, req.Offset)
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*grpcUser.User
	for rows.Next() {
		var user dto.User
		err = rows.Scan(&user.ID, &user.Email, &user.Password)
		if err != nil {
			return nil, err
		}
		users = append(users, &grpcUser.User{
			Id:       int32(user.ID),
			Email:    user.Email,
			Password: user.Password,
		})
	}

	return &grpcUser.ListUsersResponse{Users: users}, nil
}

func CreateTable() error {
	fmt.Println(" Устанавливаем соединение с базой данных PostgreSQL")
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()
	// Проверяем существование таблицы
	var tableExists bool
	err = db.QueryRow("SELECT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = $1)", "users").Scan(&tableExists)
	if err != nil {
		fmt.Println(err)
	}
	// Если таблица не существует, создаем её
	if !tableExists {
		_, err = db.Exec(`CREATE TABLE users (
            id SERIAL PRIMARY KEY,
            email VARCHAR,
            password VARCHAR
        )`)
		if err != nil {
			fmt.Println(err)
		}
		// Добавляем пользователей
		_, err = db.Exec("INSERT INTO users (email, password) VALUES ($1, $2)", "TestUser1@mail.com", "123")
		if err != nil {
			fmt.Println(err)
		}
		_, err = db.Exec("INSERT INTO users (email, password) VALUES ($1, $2)", "TestUser2@mail.com", "321")
		if err != nil {
			fmt.Println(err)
		}
		_, err = db.Exec("INSERT INTO users (email, password) VALUES ($1, $2)", "admin", "admin")
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("Таблица создана и тестовые пользователи добавлены.")
	}
	return nil
}

func main() {
	//запуск gRPC сервера
	lis, err := net.Listen("tcp", ":50051")
	log.Println("gRPC server started on port :50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	grpcUser.RegisterUserServiceServer(s, &UserStorage{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
