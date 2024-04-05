package services

import (
	"fmt"
	"github.com/go-chi/jwtauth"
	//model "go-kata/2.server/5.CA/library/dto"
	//BookModel "go-kata/2.server/5.CA/library/dto/book"
	UserModel "go-kata/2.server/5.CA/petstore/dto/user"
	"go-kata/2.server/5.CA/petstore/storage"
)

var users = make(map[string]string)

var TokenAuth *jwtauth.JWTAuth = jwtauth.New("HS256", []byte("secret"), nil)

type Servicer interface {
	//BookTake(userID int, bookID int) (string, error)
	//ReturnBook(userID int, bookID int) error
	//GetAuthorsWithBooks() ([]model.Author, error)
	//AddAuthor(name string) (string, error)
	//AddBook(title string, authorID int) (string, error)
	//GetBooks() ([]BookModel.Book, error)
	//GetUsersWithRentedBooks() ([]UserModel.User, error)
	//
	CreateUser(UserModel.User) error
	GetUserByName(name string) (UserModel.User, error)
	UpdateUserByName(name string, newUser UserModel.User) (UserModel.User, error)
	DeleteUserByName(name string) (UserModel.User, error)
	LoginUser(name, password string) error
	LogoutUser(name string) error
	CreateWithArray([]UserModel.User) error
	CreateWithList([]UserModel.User) error
	//
	//ReturnInventories(props Props)
	//PlaceOrder(order Order)
	//GetOrder(orderID int) (order Order)
	//DeleteOrder(orderID int)
	////
	//PetAdd(pet Pet)
	//PetUpdate(pet Pet)
	//PetFindByStatus(string)
	//PetGetByID(int)
	//PetUpdateByID(pet Pet)
	//PetDelete(pet Pet)
}

type Service struct {
	UserStorage storage.LibraryRepository
}

func NewService(UserStorage storage.LibraryRepository) *Service {
	return &Service{
		UserStorage: UserStorage}
}

func (s *Service) CreateUser(user UserModel.User) error {
	//var User UserModel.User
	//var Users []User
	err := s.UserStorage.CreateUser(user)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func (s *Service) GetUserByName(name string) (UserModel.User, error) {
	//var user UserModel.User
	User, err := s.UserStorage.GetUserByName(name)
	if err != nil {
		fmt.Println(err)
		return UserModel.User{}, err
	}
	return User, nil
}

func (s *Service) UpdateUserByName(name string, newUser UserModel.User) (UserModel.User, error) {
	//var user UserModel.User
	User, err := s.UserStorage.UpdateUserByName(name, newUser)
	if err != nil {
		fmt.Println(err)
		return UserModel.User{}, err
	}
	return User, nil
}

func (s *Service) DeleteUserByName(name string) (UserModel.User, error) {
	User, err := s.UserStorage.DeleteUserByName(name)
	if err != nil {
		fmt.Println(err)
		return UserModel.User{}, err
	}
	return User, nil
}

func (s *Service) LoginUser(name, password string) error {
	err := s.UserStorage.LoginUser(name, password)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func (s *Service) LogoutUser(name string) error {
	// Implement logout logic here
	err := s.UserStorage.LogoutUser(name)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func (s *Service) CreateWithArray(users []UserModel.User) error {
	// Implement logout logic here
	err := s.UserStorage.CreateWithArray(users)
	if err != nil {
		fmt.Println("Ошибка сервис:", err)
		return err
	}
	return nil
}

func (s *Service) CreateWithList(users []UserModel.User) error {
	// Implement logout logic here
	err := s.UserStorage.CreateWithList(users)
	if err != nil {
		fmt.Println("Ошибка сервис:", err)
		return err
	}
	return nil
}

//
//func (s *Service) BookTake(userID int, bookID int) (string, error) {
//	str, err := s.UserStorage.GetUserBook(userID, bookID)
//	if err != nil {
//		fmt.Println(err)
//		return "", err
//	}
//	return str, nil
//}
//
//func (s *Service) ReturnBook(userID int, bookID int) error {
//	err := s.UserStorage.ReturnBook(userID, bookID)
//	if err != nil {
//		fmt.Println(err)
//		return err
//	}
//	return nil
//}
//
//func (s *Service) GetAuthorsWithBooks() ([]model.Author, error) {
//	respond, err := s.UserStorage.GetAuthorsWithBooks()
//	if err != nil {
//		fmt.Println(err)
//		return respond, err
//	}
//	return respond, nil
//}
//
//func (s *Service) GetUsersWithRentedBooks() ([]UserModel.User, error) {
//	respond, err := s.UserStorage.GetUsersWithRentedBooks()
//	if err != nil {
//		fmt.Println(err)
//		return respond, err
//	}
//	return respond, nil
//}
//
//func (s *Service) AddAuthor(name string) (string, error) {
//	respond, err := s.UserStorage.AddAuthor(name)
//	if err != nil {
//		fmt.Println(err)
//		return "", err
//	}
//	return respond, nil
//}
//
//func (s *Service) AddBook(title string, authorID int) (string, error) {
//	respond, err := s.UserStorage.AddBook(title, authorID)
//	if err != nil {
//		fmt.Println(err)
//		return "", err
//	}
//	return respond, nil
//}
//
//func (s *Service) GetBooks() ([]BookModel.Book, error) {
//	respond, err := s.UserStorage.GetBooksWithAuthors()
//	if err != nil {
//		fmt.Println(err)
//		return respond, err
//	}
//	return respond, nil
//}
