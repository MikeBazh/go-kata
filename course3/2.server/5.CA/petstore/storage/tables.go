package storage

import (
	"database/sql"
	"fmt"
)

func CreateTables() (*sql.DB, error) {
	db, err := CreateTableUsersIfNotExists()
	db, err = CreateTablePetsIfNotExists()
	db, err = CreateTableOrdersIfNotExists()
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
	//fmt.Println("таблица пользователей найдена")
	return db, nil
}

func CreateTablePetsIfNotExists() (*sql.DB, error) {
	// Открываем соединение с базой данных
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		fmt.Println("Ошибка при подключении к базе данных:", err)
		return db, err
	}
	// Проверяем существование таблицы
	var tableExists bool
	err = db.QueryRow("SELECT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = $1)", "pets").Scan(&tableExists)
	if err != nil {
		fmt.Println(err)
	}
	// Если таблица не существует, создаем её
	if !tableExists {
		query := `
        CREATE TABLE IF NOT EXISTS pets (
            id SERIAL PRIMARY KEY,
            name VARCHAR(50) NOT NULL,
          status VARCHAR(50) NOT NULL,
            tags JSONB NOT NULL,
            category JSONB NOT NULL                     
        )
    `
		_, err := db.Exec(query)
		if err != nil {
			return db, err
		}
		fmt.Println("таблица pets создана")
	}
	//fmt.Println("таблица pets найдена")
	return db, nil
}

func CreateTableOrdersIfNotExists() (*sql.DB, error) {
	// Открываем соединение с базой данных
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		fmt.Println("Ошибка при подключении к базе данных:", err)
		return db, err
	}
	// Проверяем существование таблицы
	var tableExists bool
	err = db.QueryRow("SELECT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = $1)", "orders").Scan(&tableExists)
	if err != nil {
		fmt.Println(err)
	}
	// Если таблица не существует, создаем её
	if !tableExists {
		query := `
        CREATE TABLE IF NOT EXISTS orders (
            id SERIAL PRIMARY KEY,
            complete VARCHAR(50),
          petId INTEGER,
          quantity INTEGER,
            shipDate TIMESTAMP,
            status VARCHAR(50)                      
        )
    `
		_, err := db.Exec(query)
		if err != nil {
			return db, err
		}
		fmt.Println("таблица orders создана")
	}
	//fmt.Println("таблица orders найдена")
	return db, nil
}
