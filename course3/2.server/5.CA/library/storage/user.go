package storage

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	model "go-kata/2.server/5.CA/library/dto"
	BookModel "go-kata/2.server/5.CA/library/dto/book"
	UserModel "go-kata/2.server/5.CA/library/dto/user"
)

type LibraryRepository interface {
	GetUserBook(userID, bookID int) (string, error)
	ReturnBook(userID, bookID int) error
	GetAuthorsWithBooks() ([]model.Author, error)
	AddAuthor(name string) (string, error)
	AddBook(title string, authorID int) (string, error)
	GetBooksWithAuthors() ([]BookModel.Book, error)
	GetUsersWithRentedBooks() ([]UserModel.User, error)
}

type LibraryStorage struct {
}

const (
	connStr = "host=db user=postgres password=123 dbname=postgres sslmode=disable"
	//connStr = "user=postgres password=123 dbname=postgres sslmode=disable"
)

// NewLibraryStorage - конструктор хранилища пользователей
func NewLibraryStorage() *LibraryStorage {
	return &LibraryStorage{}
}

// GetUsersWithRentedBooks возвращает список пользователей с их арендованными книгами
func (s *LibraryStorage) GetUsersWithRentedBooks() ([]UserModel.User, error) {
	db, err := CreateTables()
	if err != nil {
		return []UserModel.User{}, err
	}
	query := `
        SELECT
            u.id AS user_id,
            u.name AS user_name,
            b.id AS book_id,
            b.title AS book_title,
            a.name AS author_name
        FROM
            users u
        LEFT JOIN
            books b ON b.id = ANY(u.rented_books)
        LEFT JOIN
            authors a ON b.author_id = a.id
    `

	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("ошибка при выполнении запроса: %v", err)
	}
	defer rows.Close()

	usersMap := make(map[int]*UserModel.User)
	for rows.Next() {
		var userID int
		var userName string
		var bookID sql.NullInt64
		var bookTitle sql.NullString
		var authorName sql.NullString

		err := rows.Scan(&userID, &userName, &bookID, &bookTitle, &authorName)
		if err != nil {
			return nil, fmt.Errorf("ошибка при сканировании строк запроса: %v", err)
		}

		user, ok := usersMap[userID]
		if !ok {
			user = &UserModel.User{
				ID:          userID,
				Name:        userName,
				RentedBooks: []UserModel.RentedBook{},
			}
			usersMap[userID] = user
		}

		if bookID.Valid && bookTitle.Valid && authorName.Valid {
			book := UserModel.RentedBook{
				ID:       int(bookID.Int64),
				Title:    bookTitle.String,
				Author:   authorName.String,
				IsRented: true,
			}
			user.RentedBooks = append(user.RentedBooks, book)
		}
	}

	users := make([]UserModel.User, 0, len(usersMap))
	for _, user := range usersMap {
		users = append(users, *user)
	}

	return users, nil
}
