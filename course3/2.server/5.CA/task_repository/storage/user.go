package storage

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"go-kata/2.server/5.CA/task_repository/dto"
	"log"
)

type UserRepository interface {
	Create(user dto.RequestUser) error
	GetByID(id int) (dto.User, error)
	Update(user dto.RequestUser) error
	Delete(id int) error
	List(limit, offset int) ([]dto.User, error)
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

func NewUserStorage() UserRepository {
	return &UserStorage{}
}

// Create - создание пользователя в БД
func (s *UserStorage) Create(user dto.RequestUser) error {
	err := CreateTable()
	if err != nil {
		fmt.Println(err)
		return err
	}
	// Устанавливаем соединение с базой данных PostgreSQL
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer db.Close()
	// Вставляем данные в таблицу
	_, err = db.Exec(`INSERT INTO users (name, email) VALUES ($1, $2)`, user.Name, user.Email)
	if err != nil {
		return err
		//log.Fatal(err)
	}
	fmt.Println("пользователь", user.Name, " добавлен")
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
	if err != nil {
		return err
	}
	//fmt.Println("данные пользователя", user.Name, " обновлены")
	return nil
}

// GetByID - получение пользователя по ID из БД
func (s *UserStorage) GetByID(userID int) (user dto.User, err error) {
	err = CreateTable()
	if err != nil {
		fmt.Println(err)
		return dto.User{}, err
	}
	// Устанавливаем соединение с базой данных
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return dto.User{}, err
	}
	defer db.Close()

	row := db.QueryRow(`SELECT * FROM users WHERE id=$1`, userID)
	err = row.Scan(&user.ID, &user.Name, &user.Email, &user.Deleted)
	return user, nil
}

// Delete(id string) error
func (s *UserStorage) Delete(userID int) error {
	// Устанавливаем соединение с базой данных PostgreSQL
	db, err := sql.Open("postgres", connStr)
	defer db.Close()
	if err != nil {
		fmt.Println(err)
		return err
	}
	_, err = db.Exec(`UPDATE users SET deleted = true WHERE id=$1`, userID)
	return err
}

func (s *UserStorage) List(limit, offset int) ([]dto.User, error) {
	err := CreateTable()
	if err != nil {
		fmt.Println(err)
		return []dto.User{}, err
	}
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

	users := []dto.User{}
	for rows.Next() {
		var user dto.User
		err = rows.Scan(&user.ID, &user.Name, &user.Email, &user.Deleted)
		if err != nil {
			return []dto.User{}, err
		}
		users = append(users, user)
	}
	return users, err
}

// Create - создание таблицы и добавление пользователей в БД
func CreateTable() error {
	// Устанавливаем соединение с базой данных PostgreSQL
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
            name VARCHAR,
            email VARCHAR,
            deleted BOOLEAN DEFAULT false
        )`)
		if err != nil {
			fmt.Println(err)
		}
		// Добавляем пользователей
		_, err = db.Exec("INSERT INTO users (name, email) VALUES ($1, $2)", "TestUser1", "email1@mail.com")
		if err != nil {
			fmt.Println(err)
		}
		_, err = db.Exec("INSERT INTO users (name, email) VALUES ($1, $2)", "TestUser2", "email2@mail.com")
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("Таблица создана и тестовые пользователи добавлены.")
	}
	return nil
}
