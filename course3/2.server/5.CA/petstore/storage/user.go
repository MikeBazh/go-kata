package storage

import (
	"fmt"
	_ "github.com/lib/pq"
	//UserModel "go-kata/2.server/5.CA/library/dto/user"
	UserModel "go-kata/2.server/5.CA/petstore/dto/user"
)

type LibraryRepository interface {
	CreateUser(UserModel.User) error
	CreateWithArray([]UserModel.User) error
	CreateWithList([]UserModel.User) error
	GetUserByName(name string) (UserModel.User, error)
	UpdateUserByName(name string, newUser UserModel.User) (UserModel.User, error)
	DeleteUserByName(name string) (UserModel.User, error)
	LoginUser(name, password string) error
	LogoutUser(name string) error
}

type LibraryStorage struct {
}

const (
	//connStr = "host=db user=postgres password=123 dbname=postgres sslmode=disable"
	connStr = "user=postgres password=123 dbname=postgres sslmode=disable"
)

// NewLibraryStorage - конструктор хранилища пользователей
func NewLibraryStorage() *LibraryStorage {
	return &LibraryStorage{}
}

func (ls *LibraryStorage) GetUserByName(name string) (UserModel.User, error) {
	db, err := CreateTableUsersIfNotExists()
	if err != nil {
		fmt.Println("Ошибка при подключении к базе данных:", err)
		return UserModel.User{}, err
	}
	var user UserModel.User
	query := "SELECT id, username, firstname, lastname, email, phone, userStatus FROM users WHERE username = $1"
	row := db.QueryRow(query, name)
	err = row.Scan(&user.Id, &user.Username, &user.FirstName, &user.LastName, &user.Email, &user.Phone, &user.UserStatus)
	user.Password = "*****"
	if err != nil {
		return UserModel.User{}, err
	}
	return user, nil
}

func (ls *LibraryStorage) UpdateUserByName(name string, newUser UserModel.User) (UserModel.User, error) {
	db, err := CreateTableUsersIfNotExists()
	if err != nil {
		fmt.Println("Ошибка при подключении к базе данных:", err)
		return UserModel.User{}, err
	}
	//var updatedUser UserModel.User
	query := "UPDATE users SET firstname = $1, lastname=$2, email = $3, phone = $4, userStatus=$5, password=$6  WHERE username = $7 RETURNING id"
	var id int
	err = db.QueryRow(query, newUser.FirstName, newUser.LastName, newUser.Email, newUser.Phone, newUser.UserStatus, newUser.Password, name).Scan(&id)
	if err != nil {
		return UserModel.User{}, err
	}
	fmt.Println("user", id, "updated")
	return newUser, nil
}

func (ls *LibraryStorage) DeleteUserByName(name string) (UserModel.User, error) {
	db, err := CreateTableUsersIfNotExists()
	if err != nil {
		fmt.Println("Ошибка при подключении к базе данных:", err)
		return UserModel.User{}, err
	}
	var deletedUser UserModel.User
	query := "DELETE FROM users WHERE username = $1 RETURNING id, username"
	row := db.QueryRow(query, name)
	err = row.Scan(&deletedUser.Id, &deletedUser.Username)
	if err != nil {
		return UserModel.User{}, err
	}
	fmt.Println("user", deletedUser.Username, "deleted")
	return deletedUser, nil
}

func (ls *LibraryStorage) LoginUser(name, password string) error {
	db, err := CreateTableUsersIfNotExists()
	if err != nil {
		fmt.Println("Ошибка при подключении к базе данных:", err)
		return err
	}
	//m:=make(map[string]string)
	var UserPassword string
	query := "SELECT password FROM users WHERE username = $1"
	err = db.QueryRow(query, name).Scan(&UserPassword)
	if err != nil {
		return err
	}
	var ID string
	if UserPassword == password {
		// Implement login logic here
		query := "UPDATE users SET userStatus = 'LoggedIn' WHERE username = $1 RETURNING id"
		row := db.QueryRow(query, name)
		err = row.Scan(&ID)
		if err != nil {
			return err
		}
		fmt.Println("user ", ID, "logged in")
	}
	return nil
}

func (ls *LibraryStorage) LogoutUser(name string) error {
	db, err := CreateTableUsersIfNotExists()
	if err != nil {
		fmt.Println("Ошибка при подключении к базе данных:", err)
		return err
	}
	var ID string
	query := "UPDATE users SET userStatus = null WHERE username = $1 RETURNING id"
	row := db.QueryRow(query, name)
	err = row.Scan(ID)
	if err != nil {
		return err
	}
	fmt.Println("user ", ID, "logged out")
	return nil
}

func (ls *LibraryStorage) CreateWithArray(users []UserModel.User) error {
	db, err := CreateTableUsersIfNotExists()
	if err != nil {
		fmt.Println("Ошибка при подключении к базе данных:", err)
		return err
	}
	for _, user := range users {
		query := "INSERT INTO users (username, firstname, lastname, email, phone, userStatus, password) VALUES ($1, $2, $3, $4, $5, $6, $7)"
		err = db.QueryRow(query, user.Username, user.FirstName, user.LastName, user.Email, user.Phone, user.UserStatus, user.Password).Err()
		//err = row.Scan(&user.Id, &user.Username, &user.FirstName, &user.LastName, &user.Email, &user.Phone, &user.UserStatus)
		if err != nil {
			return err
		}
	}
	return nil
}

func (ls *LibraryStorage) CreateWithList(users []UserModel.User) error {
	db, err := CreateTableUsersIfNotExists()
	if err != nil {
		fmt.Println("Ошибка при подключении к базе данных:", err)
		return err
	}
	for _, user := range users {
		query := "INSERT INTO users (username, firstname, lastname, email, phone, userStatus, password) VALUES ($1, $2, $3, $4, $5, $6, $7)"
		err = db.QueryRow(query, user.Username, user.FirstName, user.LastName, user.Email, user.Phone, user.UserStatus, user.Password).Err()
		//err = row.Scan(&user.Id, &user.Username, &user.FirstName, &user.LastName, &user.Email, &user.Phone, &user.UserStatus)
		if err != nil {
			return err
		}
	}
	return nil
}

func (ls *LibraryStorage) CreateUser(user UserModel.User) error {
	db, err := CreateTableUsersIfNotExists()
	if err != nil {
		fmt.Println("Ошибка при подключении к базе данных:", err)
		return err
	}
	query := "INSERT INTO users (username, firstname, lastname, email, phone, userStatus, password) VALUES ($1, $2, $3, $4, $5, $6, $7)"
	err = db.QueryRow(query, user.Username, user.FirstName, user.LastName, user.Email, user.Phone, user.UserStatus, user.Password).Err()
	//err = row.Scan(&user.Id, &user.Username, &user.FirstName, &user.LastName, &user.Email, &user.Phone, &user.UserStatus)
	if err != nil {
		return err
	}
	return nil
}

//// GetUsersWithRentedBooks возвращает список пользователей с их арендованными книгами
//func (s *LibraryStorage) GetUsersWithRentedBooks() ([]UserModel.User, error) {
//	db, err := CreateTables()
//	if err != nil {
//		return []UserModel.User{}, err
//	}
//	query := `
//        SELECT
//            u.id AS user_id,
//            u.name AS user_name,
//            b.id AS book_id,
//            b.title AS book_title,
//            a.name AS author_name
//        FROM
//            users u
//        LEFT JOIN
//            books b ON b.id = ANY(u.rented_books)
//        LEFT JOIN
//            authors a ON b.author_id = a.id
//    `
//
//	rows, err := db.Query(query)
//	if err != nil {
//		return nil, fmt.Errorf("ошибка при выполнении запроса: %v", err)
//	}
//	defer rows.Close()
//
//	usersMap := make(map[int]*UserModel.User)
//	for rows.Next() {
//		var userID int
//		var userName string
//		var bookID sql.NullInt64
//		var bookTitle sql.NullString
//		var authorName sql.NullString
//
//		err := rows.Scan(&userID, &userName, &bookID, &bookTitle, &authorName)
//		if err != nil {
//			return nil, fmt.Errorf("ошибка при сканировании строк запроса: %v", err)
//		}
//
//		user, ok := usersMap[userID]
//		if !ok {
//			user = &UserModel.User{
//				ID:          userID,
//				Name:        userName,
//				RentedBooks: []UserModel.RentedBook{},
//			}
//			usersMap[userID] = user
//		}
//
//		if bookID.Valid && bookTitle.Valid && authorName.Valid {
//			book := UserModel.RentedBook{
//				ID:       int(bookID.Int64),
//				Title:    bookTitle.String,
//				Author:   authorName.String,
//				IsRented: true,
//			}
//			user.RentedBooks = append(user.RentedBooks, book)
//		}
//	}
//
//	users := make([]UserModel.User, 0, len(usersMap))
//	for _, user := range usersMap {
//		users = append(users, *user)
//	}
//
//	return users, nil
//}
