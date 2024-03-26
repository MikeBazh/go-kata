package services

import (
	"fmt"
	"github.com/go-chi/jwtauth"
	model "go-kata/2.server/5.CA/library/dto"
	BookModel "go-kata/2.server/5.CA/library/dto/book"
	UserModel "go-kata/2.server/5.CA/library/dto/user"
	"go-kata/2.server/5.CA/library/storage"
)

var users = make(map[string]string)

var TokenAuth *jwtauth.JWTAuth = jwtauth.New("HS256", []byte("secret"), nil)

type Servicer interface {
	BookTake(userID int, bookID int) (string, error)
	ReturnBook(userID int, bookID int) error
	GetAuthorsWithBooks() ([]model.Author, error)
	AddAuthor(name string) (string, error)
	AddBook(title string, authorID int) (string, error)
	GetBooks() ([]BookModel.Book, error)
	GetUsersWithRentedBooks() ([]UserModel.User, error)
}

type Service struct {
	UserStorage storage.UserRepository
}

func NewService(UserStorage storage.UserRepository) *Service {
	return &Service{
		UserStorage: UserStorage}
}

func (s *Service) BookTake(userID int, bookID int) (string, error) {
	str, err := s.UserStorage.GetUserBook(userID, bookID)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	return str, nil
}

func (s *Service) ReturnBook(userID int, bookID int) error {
	err := s.UserStorage.ReturnBook(userID, bookID)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func (s *Service) GetAuthorsWithBooks() ([]model.Author, error) {
	respond, err := s.UserStorage.GetAuthorsWithBooks()
	if err != nil {
		fmt.Println(err)
		return respond, err
	}
	return respond, nil
}

func (s *Service) GetUsersWithRentedBooks() ([]UserModel.User, error) {
	respond, err := s.UserStorage.GetUsersWithRentedBooks()
	if err != nil {
		fmt.Println(err)
		return respond, err
	}
	return respond, nil
}

func (s *Service) AddAuthor(name string) (string, error) {
	respond, err := s.UserStorage.AddAuthor(name)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	return respond, nil
}

func (s *Service) AddBook(title string, authorID int) (string, error) {
	respond, err := s.UserStorage.AddBook(title, authorID)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	return respond, nil
}

func (s *Service) GetBooks() ([]BookModel.Book, error) {
	respond, err := s.UserStorage.GetBooksWithAuthors()
	if err != nil {
		fmt.Println(err)
		return respond, err
	}
	return respond, nil
}
