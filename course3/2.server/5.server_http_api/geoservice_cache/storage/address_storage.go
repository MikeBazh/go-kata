package storage

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

type UserRepository interface {
	CreateSearchHistory(RequestedAddress string) (SearchHistoryID int, err error)
	CreateRespondHistory(SearchResponse []byte, value string) (RespondHistoryID int, err error)
	CreateHistorySearch(search_id, address_id int) error
	CheckHistory(query string) ([][]byte, error)
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

// CreateSearchHistory Добавление данных запроса и ответов от внешнего сервиса
func (s *UserStorage) CreateSearchHistory(RequestedAddress string) (SearchHistoryID int, err error) {
	columns := map[string]string{
		"addresses": "VARCHAR",
	}
	err = CreateTable("search_history", columns)
	if err != nil {
		fmt.Println(err)
		return 0, err
	}
	// Устанавливаем соединение с базой данных PostgreSQL
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		fmt.Println(err)
		return 0, err
	}
	defer db.Close()
	// Вставляем данные в таблицу
	query := "INSERT INTO search_history (addresses) VALUES ($1) RETURNING id"
	err = db.QueryRow(query, RequestedAddress).Scan(&SearchHistoryID)
	if err != nil {
		return 0, err
		fmt.Println(err)
	}
	return SearchHistoryID, err
}

func (s *UserStorage) CreateRespondHistory(SearchResponse []byte, value string) (RespondHistoryID int, err error) {
	columns := map[string]string{
		"addresses":     "JSONB",
		"address_value": "VARCHAR",
	}
	err = CreateTable("address", columns)
	if err != nil {
		fmt.Println(err)
		return 0, err
	}
	// Устанавливаем соединение с базой данных PostgreSQL
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		fmt.Println(err)
		return 0, err
	}
	defer db.Close()
	// Вставляем данные в таблицу
	err = db.QueryRow("INSERT INTO address (addresses, address_value) VALUES ($1, $2) RETURNING id", SearchResponse, value).Scan(&RespondHistoryID)
	if err != nil {
		return 0, err
		fmt.Println(err)
	}
	return RespondHistoryID, err
}

func (s *UserStorage) CreateHistorySearch(search_id, address_id int) error {
	columns := map[string]string{
		"search_history_id": "INTEGER",
		"address_id":        "INTEGER",
	}
	err := CreateTable("history_search_address", columns)
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
	//query := "INSERT INTO history_search_address (search_history_id, address_id) VALUES ($1, $2)", search_id, address_id
	_, err = db.Exec("INSERT INTO history_search_address (search_history_id, address_id) VALUES ($1, $2)", search_id, address_id)
	if err != nil {
		return err
		fmt.Println(err)
	}
	return err
}

func (s *UserStorage) CheckHistory(query string) ([][]byte, error) {
	// Устанавливаем соединение с базой данных PostgreSQL
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		fmt.Println(err)
		return [][]byte{}, err
	}
	defer db.Close()

	var searchID int
	var searchIDs []int
	var QueryFromHistory string
	var QuerysFromHistory []string

	// Выполняем запрос для поиска текущего запроса в истории запросов
	row, err := db.Query(`
	SELECT id, addresses
	FROM search_history
	WHERE similarity(addresses, $1) >0.7
	ORDER BY similarity(addresses, $1) DESC;
	`, query)

	if err != nil {
		fmt.Println(err)
		return [][]byte{}, err
	}
	for row.Next() {
		err = row.Scan(&searchID, &QueryFromHistory)
		if err != nil {
			fmt.Println(err)
			return [][]byte{}, err
		}
		searchIDs = append(searchIDs, searchID)
		QuerysFromHistory = append(QuerysFromHistory, QueryFromHistory)
	}
	defer row.Close()

	var addressID int
	var addressIDs []int
	var address []byte
	var addressList [][]byte

	if len(searchIDs) > 0 {
		fmt.Println("repository: найдено в кэше по запросу: '", QuerysFromHistory[0], "'")
		// Выполняем запрос для поиска ID адресов, связанных с найденными ID запросов в истории
		rows, err := db.Query("SELECT address_id FROM history_search_address WHERE search_history_id = $1", searchID)
		if err != nil {
			fmt.Println(err)
			return [][]byte{}, err
		}
		defer rows.Close()

		// Сканируем результаты запроса
		for rows.Next() {
			err := rows.Scan(&addressID)
			if err != nil {
				fmt.Println(err)
				return [][]byte{}, err
			}

			addressIDs = append(addressIDs, addressID)

			// Выполняем запрос для выбора адресов по полученным ID
			rows2, err := db.Query("SELECT addresses FROM address WHERE id = $1", addressID)
			if err != nil {
				fmt.Println(err)
				return [][]byte{}, err
			}
			defer rows2.Close()

			// Сканируем адреса из запроса
			for rows2.Next() {
				err := rows2.Scan(&address)
				if err != nil {
					fmt.Println(err)
					return [][]byte{}, err
				}
			}
			addressList = append(addressList, address)

		}
		return addressList, err
	} else {
		err = fmt.Errorf("repository: не найдено в кэше")
		return [][]byte{}, err
	}
}

func CreateTable(TableName string, Columns map[string]string) error {
	// Устанавливаем соединение с базой данных PostgreSQL
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return err
	}
	defer db.Close()

	// Проверяем существование таблицы
	var tableExists bool
	err = db.QueryRow("SELECT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = $1)", TableName).Scan(&tableExists)
	if err != nil {
		return err
	}
	// Если таблица не существует, создаем её
	if !tableExists {
		var columnsStr string
		for columnName, columnType := range Columns {
			columnsStr += fmt.Sprintf("%s %s, ", columnName, columnType)
		}
		query := fmt.Sprintf("CREATE TABLE \"%s\" (id SERIAL PRIMARY KEY, %s)", TableName, columnsStr)
		// Удаляем последнюю запятую и пробел из строки столбцов
		query = query[:len(query)-3] + ")"
		_, err = db.Exec(query)
		if err != nil {
			return err
		}
	}
	return nil
}
