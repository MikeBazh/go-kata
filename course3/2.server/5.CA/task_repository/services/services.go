package services

import (
	"errors"
	"fmt"
	"go-kata/2.server/5.CA/task_repository/dto"
	"go-kata/2.server/5.CA/task_repository/storage"
)

//var users = make(map[string]string)
//
//var TokenAuth *jwtauth.JWTAuth = jwtauth.New("HS256", []byte("secret"), nil)

type Servicer interface {
	AddUser(user dto.RequestUser) error
	GetUserByID(id int) (dto.User, error)
	UpdateUser(user dto.RequestUser) error
	DeleteByID(id int) error
	List(limit, offset int) ([]dto.User, error)
}

type Service struct {
	UserStorage storage.UserRepository
}

func NewService(UserStorage storage.UserRepository) *Service {
	return &Service{
		UserStorage: UserStorage}
}

func (s *Service) AddUser(user dto.RequestUser) error {
	//UserStorage := storage.NewUserStorage()
	err := s.UserStorage.Create(user)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) GetUserByID(id int) (dto.User, error) {
	//UserStorage := storage.NewUserStorage()
	user, err := s.UserStorage.GetByID(id)
	if err != nil {
		fmt.Println(err)
	}
	return user, err
}

func (s *Service) UpdateUser(user dto.RequestUser) error {
	//UserStorage := storage.NewUserStorage()
	err := s.UserStorage.Update(user)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("данные пользователя обновлены")
	return nil
}

func (s *Service) DeleteByID(id int) error {
	//UserStorage := storage.NewUserStorage()
	user, err := s.UserStorage.GetByID(id)
	if err != nil || user.Deleted {
		err = errors.New("пользователь уже удален или не существует")
		fmt.Println(err)
		return err
	}
	err = s.UserStorage.Delete(id)
	return err
}

func (s *Service) List(limit, offset int) ([]dto.User, error) {
	//UserStorage := storage.NewUserStorage()
	UsersReceived, err := s.UserStorage.List(limit, offset)
	return UsersReceived, err
}
