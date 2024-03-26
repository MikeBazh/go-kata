package storage

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	model "go-kata/2.server/5.CA/library/dto"
	BookModel "go-kata/2.server/5.CA/library/dto/book"
	UserModel "go-kata/2.server/5.CA/library/dto/user"
)

type UserRepository interface {
	GetUserBook(userID, bookID int) (string, error)
	ReturnBook(userID, bookID int) error
	GetAuthorsWithBooks() ([]model.Author, error)
	AddAuthor(name string) (string, error)
	AddBook(title string, authorID int) (string, error)
	GetBooksWithAuthors() ([]BookModel.Book, error)
	GetUsersWithRentedBooks() ([]UserModel.User, error)
	// Другие методы, необходимые для работы с БД
}

// UserStorage - хранилище пользователей
type UserStorage struct {
}

const (
	connStr = "host=db user=postgres password=123 dbname=postgres sslmode=disable"
	//connStr = "user=postgres password=123 dbname=postgres sslmode=disable"
)

// NewUserStorage - конструктор хранилища пользователей
func NewUserStorage() *UserStorage {
	return &UserStorage{}
}

// GetUserBook позволяет пользователю получить книгу по ее ID
func (s *UserStorage) GetUserBook(userID, bookID int) (string, error) {
	db, err := CreateTables()
	if err != nil {
		return "", err
	}
	defer db.Close()
	// Проверяем, существует ли пользователь с указанным ID
	var userExists bool
	err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM Users WHERE id = $1)", userID).Scan(&userExists)
	if err != nil {
		return "", fmt.Errorf("ошибка при проверке существования пользователя: %v", err)
	}
	if !userExists {
		return fmt.Sprintf("Ошибка: пользователь с ID %d не найден", userID), nil
	}

	// Проверяем, существует ли книга с указанным ID и не арендована ли она уже
	var bookExists, bookAvailable bool
	err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM Books WHERE id = $1)", bookID).Scan(&bookExists)
	if err != nil {
		return "", fmt.Errorf("ошибка при проверке существования книги: %v", err)
	}
	if !bookExists {
		return fmt.Sprintf("Ошибка: книга с ID %d не найдена", bookID), nil
	}

	err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM Books WHERE id = $1 AND rented_by IS NULL)", bookID).Scan(&bookAvailable)
	if err != nil {
		return "", fmt.Errorf("ошибка при проверке существования книги: %v", err)
	}
	if !bookAvailable {
		return fmt.Sprintf("Ошибка: книга с ID %d уже арендована", bookID), nil
	}

	// Обновляем запись о книге в базе данных, указывая ID пользователя, который ее получил
	_, err = db.Exec("UPDATE Books SET rented_by = $1 WHERE id = $2", userID, bookID)
	if err != nil {
		return "", fmt.Errorf("ошибка при обновлении записи о книге: %v", err)
	}
	// Обновляем столбец rented_books в таблице Users
	_, err = db.Exec("UPDATE Users SET rented_books = array_append(rented_books, $1) WHERE id = $2", bookID, userID)
	if err != nil {
		return "", fmt.Errorf("ошибка при обновлении столбца rented_books: %v", err)
	}
	return "Ок", nil
}

// ReturnBook позволяет пользователю вернуть книгу по ее ID
func (s *UserStorage) ReturnBook(userID, bookID int) error {
	db, err := CreateTables()
	if err != nil {
		return err
	}
	// Открываем соединение с базой данных
	//db, err := sql.Open("postgres", connStr)
	//if err != nil {
	//	fmt.Println("Ошибка при подключении к базе данных:", err)
	//	return err
	//}
	defer db.Close()
	// Проверяем, существует ли пользователь с указанным ID
	var userExists bool
	err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM Users WHERE id = $1)", userID).Scan(&userExists)
	if err != nil {
		return fmt.Errorf("ошибка при проверке существования пользователя: %v", err)
	}
	if !userExists {
		return fmt.Errorf("пользователь с ID %d не найден", userID)
	}

	// Проверяем, арендована ли книга с указанным ID пользователем
	var bookRented bool
	err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM Books WHERE id = $1 AND rented_by = $2)", bookID, userID).Scan(&bookRented)
	if err != nil {
		return fmt.Errorf("ошибка при проверке арендованной книги: %v", err)
	}
	if !bookRented {
		return fmt.Errorf("книга с ID %d не арендована пользователем с ID %d", bookID, userID)
	}

	// Обновляем запись о книге в базе данных, убирая информацию об аренде
	_, err = db.Exec("UPDATE Books SET rented_by = NULL WHERE id = $1", bookID)
	if err != nil {
		return fmt.Errorf("ошибка при обновлении записи о книге: %v", err)
	}

	// Обновляем столбец rented_books в таблице Users, убирая ID возвращенной книги
	_, err = db.Exec("UPDATE users SET rented_books = array_remove(rented_books, $1) WHERE id = $2", bookID, userID)
	if err != nil {
		return fmt.Errorf("ошибка при обновлении столбца rented_books: %v", err)
	}

	return nil
}

// GetAuthorsWithBooks возвращает список авторов с их книгами
func (s *UserStorage) GetAuthorsWithBooks() ([]model.Author, error) {
	//func (s *UserStorage) GetAuthorsWithBooks(db *sql.DB) ([]Author, error) {
	db, err := CreateTables()
	if err != nil {
		return []model.Author{}, err
	}
	// Открываем соединение с базой данных
	//db, err := sql.Open("postgres", connStr)
	//if err != nil {
	//	fmt.Println("Ошибка при подключении к базе данных:", err)
	//	return []model.Author{}, err
	//}
	defer db.Close()
	rows, err := db.Query(`
    SELECT
      a.id AS author_id,
      a.name AS author_name,
      b.id AS book_id,
      b.title AS book_title
    FROM
      Authors a
    LEFT JOIN
      Books b ON a.id = b.author_id
  `)
	if err != nil {
		return nil, fmt.Errorf("ошибка при выполнении запроса к базе данных: %v", err)
	}
	defer rows.Close()

	authorsMap := make(map[int]*model.Author)
	for rows.Next() {
		var authorID int
		var authorName string
		var bookID sql.NullInt64
		var bookTitle sql.NullString
		if err := rows.Scan(&authorID, &authorName, &bookID, &bookTitle); err != nil {
			return nil, fmt.Errorf("ошибка при сканировании строк запроса: %v", err)
		}

		author, ok := authorsMap[authorID]
		if !ok {
			author = &model.Author{
				ID:    authorID,
				Name:  authorName,
				Books: []model.Book{},
			}
			authorsMap[authorID] = author
		}

		if bookID.Valid && bookTitle.Valid {
			author.Books = append(author.Books, model.Book{
				ID:    int(bookID.Int64),
				Title: bookTitle.String,
			})
		}
	}

	authors := make([]model.Author, 0, len(authorsMap))
	for _, author := range authorsMap {
		authors = append(authors, *author)
	}
	return authors, nil
}

func (s *UserStorage) AddAuthor(name string) (string, error) {
	db, err := CreateTables()
	if err != nil {
		return "", err
	}
	// Открываем соединение с базой данных
	//db, err := sql.Open("postgres", connStr)
	//if err != nil {
	//	fmt.Println("Ошибка при подключении к базе данных:", err)
	//	return err
	//}
	defer db.Close()
	// Проверяем, существует ли уже автор с таким именем
	var exists bool
	err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM Authors WHERE name = $1)", name).Scan(&exists)
	if err != nil {
		return "", fmt.Errorf("ошибка при проверке существования автора: %v", err)
	}
	if exists {
		return fmt.Sprintf("автор с именем '%s' уже существует", name), nil
	}

	// Выполняем запрос к базе данных для добавления нового автора
	_, err = db.Exec("INSERT INTO Authors (name) VALUES ($1)", name)
	if err != nil {
		return "", fmt.Errorf("ошибка при добавлении автора: %v", err)
	}

	return "Ок", nil
}

// AuthorExists проверяет существование автора по его ID
func AuthorExists(db *sql.DB, authorID int) (bool, error) {
	var exists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM Authors WHERE id = $1)", authorID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("ошибка при проверке существования автора: %v", err)
	}
	return exists, nil
}

// AddBook добавляет новую книгу с указанным ID автора
func (s *UserStorage) AddBook(title string, authorID int) (string, error) {
	db, err := CreateTables()
	if err != nil {
		return "", err
	}
	// Открываем соединение с базой данных
	//db, err := sql.Open("postgres", connStr)
	//if err != nil {
	//	fmt.Println("Ошибка при подключении к базе данных:", err)
	//	return err
	//}
	// Проверяем существование автора
	exists, err := AuthorExists(db, authorID)
	if err != nil {
		return "", err
	}
	if !exists {
		return fmt.Sprintf("Ошибка: автора с ID %d не существует", authorID), nil
	}

	// Выполняем запрос к базе данных для добавления новой книги
	_, err = db.Exec("INSERT INTO Books (title, author_id) VALUES ($1, $2)", title, authorID)
	if err != nil {
		return "", fmt.Errorf("ошибка при добавлении книги: %v", err)
	}

	return "Ок", nil
}

// GetBooksWithAuthors возвращает список всех книг с информацией об авторах.
func (s *UserStorage) GetBooksWithAuthors() ([]BookModel.Book, error) {
	db, err := CreateTables()
	if err != nil {
		return []BookModel.Book{}, err
	}
	query := `
    SELECT
      b.id AS book_id,
      b.title AS book_title,
      a.id AS author_id,
      a.name AS author_name
    FROM
      books b
    LEFT JOIN
      authors a ON b.author_id = a.id
  `
	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("ошибка при выполнении запроса: %v", err)
	}
	defer rows.Close()

	var books []BookModel.Book
	for rows.Next() {
		var book BookModel.Book
		var author BookModel.Author
		err := rows.Scan(&book.ID, &book.Title, &author.ID, &author.Name)
		if err != nil {
			return nil, fmt.Errorf("ошибка при сканировании строк запроса: %v", err)
		}
		book.Author = author
		books = append(books, book)
	}

	return books, nil
}

// GetUsersWithRentedBooks возвращает список пользователей с их арендованными книгами
func (s *UserStorage) GetUsersWithRentedBooks() ([]UserModel.User, error) {
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
