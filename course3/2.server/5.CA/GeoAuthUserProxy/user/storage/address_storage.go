package storage

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"go-kata/2.server/5.CA/geoservice_user/dto"
)

type UserRepository interface {
	Create(user dto.RequestUser) error
	GetByEmail(userEmail string) (user dto.User, err error)
	List(limit, offset int) ([]dto.User, error)
}

type UserStorage struct {
}

const (
	//connStr = "host=db user=postgres password=123 dbname=postgres sslmode=disable"
	connStr = "host=172.17.0.2 user=postgres password=123/ dbname=postgres sslmode=disable"
)

// NewUserStorage - конструктор хранилища пользователей
//func NewUserStorage() *UserStorage {
//	return &UserStorage{}
//}

func NewUserStorage() UserRepository {
	return &UserStorage{}
}

// CreateSearchHistory Добавление данных запроса и ответов от внешнего сервиса
//func (s *UserStorage) CreateSearchHistory(RequestedAddress string) (SearchHistoryID int, err error) {
//	columns := map[string]string{
//		"addresses": "VARCHAR",
//	}
//	err = CreateTable("search_history", columns)
//	if err != nil {
//		fmt.Println(err)
//		return 0, err
//	}
//	// Устанавливаем соединение с базой данных PostgreSQL
//	db, err := sql.Open("postgres", connStr)
//	if err != nil {
//		fmt.Println(err)
//		return 0, err
//	}
//	defer db.Close()
//	// Вставляем данные в таблицу
//	query := "INSERT INTO search_history (addresses) VALUES ($1) RETURNING id"
//	err = db.QueryRow(query, RequestedAddress).Scan(&SearchHistoryID)
//	if err != nil {
//		return 0, err
//		fmt.Println(err)
//	}
//	return SearchHistoryID, err
//}
//
//func (s *UserStorage) CreateRespondHistory(SearchResponse []byte, value string) (RespondHistoryID int, err error) {
//	columns := map[string]string{
//		"addresses":     "JSONB",
//		"address_value": "VARCHAR",
//	}
//	err = CreateTable("address", columns)
//	if err != nil {
//		fmt.Println(err)
//		return 0, err
//	}
//	// Устанавливаем соединение с базой данных PostgreSQL
//	db, err := sql.Open("postgres", connStr)
//	if err != nil {
//		fmt.Println(err)
//		return 0, err
//	}
//	defer db.Close()
//	// Вставляем данные в таблицу
//	err = db.QueryRow("INSERT INTO address (addresses, address_value) VALUES ($1, $2) RETURNING id", SearchResponse, value).Scan(&RespondHistoryID)
//	if err != nil {
//		return 0, err
//		fmt.Println(err)
//	}
//	return RespondHistoryID, err
//}
//
//func (s *UserStorage) CreateHistorySearch(search_id, address_id int) error {
//	columns := map[string]string{
//		"search_history_id": "INTEGER",
//		"address_id":        "INTEGER",
//	}
//	err := CreateTable("history_search_address", columns)
//	if err != nil {
//		fmt.Println(err)
//		return err
//	}
//	// Устанавливаем соединение с базой данных PostgreSQL
//	db, err := sql.Open("postgres", connStr)
//	if err != nil {
//		fmt.Println(err)
//		return err
//	}
//	defer db.Close()
//	// Вставляем данные в таблицу
//	//query := "INSERT INTO history_search_address (search_history_id, address_id) VALUES ($1, $2)", search_id, address_id
//	_, err = db.Exec("INSERT INTO history_search_address (search_history_id, address_id) VALUES ($1, $2)", search_id, address_id)
//	if err != nil {
//		return err
//		fmt.Println(err)
//	}
//	return err
//}
//
//func (s *UserStorage) CheckHistory(query string) ([][]byte, error) {
//	// Устанавливаем соединение с базой данных PostgreSQL
//	db, err := sql.Open("postgres", connStr)
//	if err != nil {
//		fmt.Println(err)
//		return [][]byte{}, err
//	}
//	defer db.Close()
//
//	var searchID int
//	var searchIDs []int
//	var QueryFromHistory string
//	var QuerysFromHistory []string
//
//	columns := map[string]string{
//		"addresses": "VARCHAR",
//	}
//	err = CreateTable("search_history", columns)
//	if err != nil {
//		fmt.Println(err)
//		return [][]byte{}, err
//	}
//
//	// Выполняем запрос для поиска текущего запроса в истории запросов
//	row, err := db.Query(`
//	SELECT id, addresses
//	FROM search_history
//	WHERE similarity(addresses, $1) >0.7
//	ORDER BY similarity(addresses, $1) DESC;
//	`, query)
//
//	if err != nil {
//		fmt.Println(err)
//		return [][]byte{}, err
//	}
//	for row.Next() {
//		err = row.Scan(&searchID, &QueryFromHistory)
//		if err != nil {
//			fmt.Println(err)
//			return [][]byte{}, err
//		}
//		searchIDs = append(searchIDs, searchID)
//		QuerysFromHistory = append(QuerysFromHistory, QueryFromHistory)
//	}
//	defer row.Close()
//
//	var addressID int
//	var addressIDs []int
//	var address []byte
//	var addressList [][]byte
//
//	if len(searchIDs) > 0 {
//		fmt.Println("repository: найдено в БД по запросу: '", QuerysFromHistory[0], "'")
//		// Выполняем запрос для поиска ID адресов, связанных с найденными ID запросов в истории
//		rows, err := db.Query("SELECT address_id FROM history_search_address WHERE search_history_id = $1", searchID)
//		if err != nil {
//			fmt.Println(err)
//			return [][]byte{}, err
//		}
//		defer rows.Close()
//
//		// Сканируем результаты запроса
//		for rows.Next() {
//			err := rows.Scan(&addressID)
//			if err != nil {
//				fmt.Println(err)
//				return [][]byte{}, err
//			}
//
//			addressIDs = append(addressIDs, addressID)
//
//			// Выполняем запрос для выбора адресов по полученным ID
//			rows2, err := db.Query("SELECT addresses FROM address WHERE id = $1", addressID)
//			if err != nil {
//				fmt.Println(err)
//				return [][]byte{}, err
//			}
//			defer rows2.Close()
//
//			// Сканируем адреса из запроса
//			for rows2.Next() {
//				err := rows2.Scan(&address)
//				if err != nil {
//					fmt.Println(err)
//					return [][]byte{}, err
//				}
//			}
//			addressList = append(addressList, address)
//
//		}
//		return addressList, err
//	} else {
//		err = fmt.Errorf("repository: не найдено в БД")
//		return [][]byte{}, err
//	}
//}

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
	_, err = db.Exec(`INSERT INTO users (email, password) VALUES ($1, $2)`, user.Email, user.Password)
	if err != nil {
		return err
		//log.Fatal(err)
	}
	fmt.Println("пользователь", user.Email, " добавлен")
	return err
}

// GetByEmail - получение пользователя по ID из БД
func (s *UserStorage) GetByEmail(userEmail string) (user dto.User, err error) {
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

	row := db.QueryRow(`SELECT * FROM users WHERE Email=$1`, userEmail)
	err = row.Scan(&user.ID, &user.Email, &user.Password)
	return user, nil
}

//// Delete(id string) error
//func (s *UserStorage) Delete(userID int) error {
//	// Устанавливаем соединение с базой данных PostgreSQL
//	db, err := sql.Open("postgres", connStr)
//	defer db.Close()
//	if err != nil {
//		fmt.Println(err)
//		return err
//	}
//	_, err = db.Exec(`UPDATE users SET deleted = true WHERE id=$1`, userID)
//	return err
//}

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
		err = rows.Scan(&user.ID, &user.Email, &user.Password)
		if err != nil {
			return []dto.User{}, err
		}
		users = append(users, user)
	}
	return users, err
}

// Create - создание таблицы и добавление пользователей в БД
func CreateTable() error {
	fmt.Println(" Устанавливаем соединение с базой данных PostgreSQL")
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
            email VARCHAR,
            password VARCHAR
        )`)
		if err != nil {
			fmt.Println(err)
		}
		// Добавляем пользователей
		_, err = db.Exec("INSERT INTO users (email, password) VALUES ($1, $2)", "TestUser1@mail.com", "123")
		if err != nil {
			fmt.Println(err)
		}
		_, err = db.Exec("INSERT INTO users (email, password) VALUES ($1, $2)", "TestUser2@mail.com", "321")
		if err != nil {
			fmt.Println(err)
		}
		_, err = db.Exec("INSERT INTO users (email, password) VALUES ($1, $2)", "admin", "admin")
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("Таблица создана и тестовые пользователи добавлены.")
	}
	return nil
}
