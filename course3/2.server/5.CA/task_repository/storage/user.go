package storage

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"go-kata/2.server/5.CA/task_repository/dto"
	"log"
)

type UserRepository interface {
	Create(user dto.User) error
	GetByID(id int) (dto.User, error)
	Update(user dto.User) error
	Delete(id int) error
	List(l, o string) ([]dto.User, error)
	CreateTable() error
	// Другие методы, необходимые для работы с пользователями
}

// UserStorage - хранилище пользователей
type UserStorage struct {
	//adapter *adapter.SQLAdapter
	//cache   cache.Cache
}

const (
	//connStr = "user=postgres password=123 dbname=postgres sslmode=disable"
	connStr = "host=db user=postgres password=123 dbname=postgres sslmode=disable"
)

// NewUserStorage - конструктор хранилища пользователей
func NewUserStorage() *UserStorage {
	return &UserStorage{}
}

// Create - создание пользователя в БД
func (s *UserStorage) CreateTable() error {
	// Устанавливаем соединение с базой данных PostgreSQL
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	// Создаем таблицу
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    name VARCHAR,
    email VARCHAR,
    deleted BOOLEAN DEFAULT false
  )`)
	if err != nil {
		log.Fatal(err)
	}
	return err
}

// Create - создание пользователя в БД
func (s *UserStorage) Create(user dto.RequestUser) error {
	err := s.CreateTable()
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println("таблица создана")
	// Устанавливаем соединение с базой данных PostgreSQL
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		fmt.Println(err)
		return err
		//log.Fatal(err)
	}
	defer db.Close()
	// Вставляем данные в таблицу
	_, err = db.Exec(`INSERT INTO users (name, email) VALUES ($1, $2)`, user.Name, user.Email)
	if err != nil {
		return err
		//log.Fatal(err)
	}
	return err
}

// Update - обновление пользователя в БД
func (s *UserStorage) Update(user dto.RequestUser) error {
	// Устанавливаем соединение с базой данных PostgreSQL
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		//fmt.Println(err)
		//return err
		log.Fatal(err)
	}
	defer db.Close()
	// Обновляем данные в таблице
	_, err = db.Exec(`UPDATE users 
SET email = $2 WHERE name = $1`, user.Name, user.Email)
	return err
}

// GetByID - получение пользователя по ID из БД
func (s *UserStorage) GetByID(userID int) (user dto.User, err error) {
	// Устанавливаем соединение с базой данных PostgreSQL
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return dto.User{}, err
	}
	defer db.Close()

	row := db.QueryRow(`SELECT * FROM users WHERE id=$1`, userID)
	err = row.Scan(&user.ID, &user.Name, &user.Email, &user.Deleted)
	//if user.Deleted {
	//	err = fmt.Errorf("user не существует (удален)")
	fmt.Println(err)
	//	return user, err
	//}
	return user, err
}

// Delete(id string) error
func (s *UserStorage) Delete(userID int) error {
	// Устанавливаем соединение с базой данных PostgreSQL
	db, err := sql.Open("postgres", connStr)
	//if err != nil {
	//	return dto.User{}, err
	//}
	defer db.Close()
	//var user dto.User
	//row := db.QueryRow(`SELECT * FROM users WHERE id=$1`, userID)
	//err = row.Scan(&user.ID, &user.Name, &user.Email, &user.Deleted)
	//if user.Deleted {
	//	err = fmt.Errorf("пользователь уже удален")
	//	fmt.Println(err)
	//	return err
	//}
	if err != nil {
		fmt.Println(err)
		return err
	}
	_, err = db.Exec(`UPDATE users SET deleted = true WHERE id=$1`, userID)
	fmt.Println(err)
	//err = row.Scan(&user.ID, &user.Name, &user.Email)
	//if err != nil {
	//	return dto.User{}, err
	//}
	return err
}

func (s *UserStorage) List(limit, offset int) ([]dto.User, error) {
	// Устанавливаем соединение с базой данных PostgreSQL
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return []dto.User{}, err
	}
	defer db.Close()
	// SQL-запрос для получения списка пользователей с учетом пагинации
	query := fmt.Sprintf(`SELECT * FROM users ORDER BY id LIMIT %d OFFSET %d`, limit, offset)

	// Выполнение запроса к базе данных
	rows, err := db.Query(query)
	if err != nil {
		fmt.Println(err)
		return []dto.User{}, err
	}
	defer rows.Close()

	// Список пользователей
	users := []dto.User{}

	// Итерация по результатам запроса и добавление пользователей в список
	for rows.Next() {
		var user dto.User
		err = rows.Scan(&user.ID, &user.Name, &user.Email, &user.Deleted)
		//http.Error(w, fmt.Sprintf("Ошибка сканирования строки результата: %v", err), http.StatusInternalServerError)
		if err != nil {
			return []dto.User{}, err
		}
		users = append(users, user)
	}
	return users, err
}
