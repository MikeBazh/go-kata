package storage

import (
	"database/sql"
	"fmt"
)

func CreateTables() (*sql.DB, error) {
	//// Открываем соединение с базой данных
	//db, err := sql.Open("postgres", connStr)
	//if err != nil {
	//	fmt.Println("Ошибка при подключении к базе данных:", err)
	//	return db, err
	//}
	//
	//// Проверяем существование таблицы Authors
	//var tableAuthorsExists bool
	//err = db.QueryRow("SELECT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = $1)", "authors").Scan(&tableAuthorsExists)
	//if err != nil {
	//	fmt.Println(err)
	//}
	//// Если таблица не существует, создаем её
	//if !tableAuthorsExists {
	//
	//	// Создание таблицы Authors
	//	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS Authors (
	//  id SERIAL PRIMARY KEY,
	//  name VARCHAR(100) NOT NULL
	//)`)
	//	if err != nil {
	//		fmt.Println("Ошибка при создании таблицы Authors:", err)
	//		return db, err
	//	}
	//	// Заполнение таблицы Authors данными
	//	for i := 0; i < 10; i++ {
	//		name := gofakeit.Name()
	//		_, err := db.Exec("INSERT INTO Authors (name) VALUES ($1)", name)
	//		if err != nil {
	//			fmt.Println("Ошибка при добавлении автора:", err)
	//		}
	//	}
	//	fmt.Println("заполнена таблица authors")
	//}
	//
	//// Проверяем существование таблицы Users
	//var tableUsersExists bool
	//err = db.QueryRow("SELECT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = $1)", "users").Scan(&tableUsersExists)
	//if err != nil {
	//	fmt.Println(err)
	//}
	//// Если таблица не существует, создаем её
	//if !tableUsersExists {
	//	// Создание таблицы Users
	//	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS Users (
	//  id SERIAL PRIMARY KEY,
	//  name VARCHAR(100) NOT NULL,
	//  rented_books INTEGER[] DEFAULT '{}'::INTEGER[]
	//)`)
	//	if err != nil {
	//		fmt.Println("Ошибка при создании таблицы Users:", err)
	//		return db, err
	//	}
	//	// Заполнение таблицы Users данными
	//	for i := 0; i < 51; i++ {
	//		name := gofakeit.FirstName()
	//		_, err := db.Exec("INSERT INTO Users (name) VALUES ($1)", name)
	//		if err != nil {
	//			fmt.Println("Ошибка при добавлении пользователя:", err)
	//		}
	//	}
	//	fmt.Println("заполнена таблица users")
	//}
	//
	//// Проверяем существование таблицы Books
	//var tableBooksExists bool
	//err = db.QueryRow("SELECT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = $1)", "books").Scan(&tableBooksExists)
	//if err != nil {
	//	fmt.Println(err)
	//}
	//// Если таблица не существует, создаем её
	//if !tableBooksExists {
	//	// Создание таблицы Books
	//	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS Books (
	//  id SERIAL PRIMARY KEY,
	//  title VARCHAR(255) NOT NULL,
	//  author_id INTEGER REFERENCES Authors(id),
	//  rented_by INTEGER REFERENCES Users(id) DEFAULT NULL
	//)`)
	//	if err != nil {
	//		fmt.Println("Ошибка при создании таблицы Books:", err)
	//		return db, err
	//	}
	//	// Заполнение таблицы Books данными
	//	rand.Seed(100)
	//	for i := 0; i < 100; i++ {
	//		gofakeit.Seed(0)
	//		title := gofakeit.Sentence(3)
	//		authorID := rand.Intn(9) + 1 // Выбираем случайный ID автора
	//		_, err := db.Exec("INSERT INTO Books (title, author_id) VALUES ($1, $2)", title, authorID)
	//		if err != nil {
	//			fmt.Println("Ошибка при добавлении книги:", err)
	//		}
	//	}
	//	fmt.Println("заполнена таблица books")
	//}
	////fmt.Println("наличие таблиц при запуске: ", tableAuthorsExists, tableUsersExists, tableBooksExists)
	////fmt.Println(err)
	//return db, nil
	db, err := CreateTableUsersIfNotExists()
	if err != nil {
		fmt.Println("Ошибка при подключении к базе данных:", err)
		return db, err
	}
	return db, err
}

func CreateTableUsersIfNotExists() (*sql.DB, error) {
	// Открываем соединение с базой данных
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		fmt.Println("Ошибка при подключении к базе данных:", err)
		return db, err
	}
	fmt.Println("here")
	// Проверяем существование таблицы Authors
	var tableUsersExists bool
	err = db.QueryRow("SELECT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = $1)", "users").Scan(&tableUsersExists)
	if err != nil {
		fmt.Println(err)
	}
	// Если таблица не существует, создаем её
	if !tableUsersExists {
		query := `
        CREATE TABLE IF NOT EXISTS users (
            id SERIAL PRIMARY KEY,
            username VARCHAR(50) UNIQUE NOT NULL,
            firstname VARCHAR(50) NOT NULL,
            lastname VARCHAR(50) NOT NULL,
            email VARCHAR(50) NOT NULL,
            password VARCHAR(100) NOT NULL,
        	phone VARCHAR(100) NOT NULL,
        	userStatus VARCHAR(100) NOT NULL                         
        )
    `
		_, err := db.Exec(query)
		if err != nil {
			return db, err
		}
		fmt.Println("таблица пользователей создана")
	}
	fmt.Println("таблица пользователей найдена")
	return db, nil
}
