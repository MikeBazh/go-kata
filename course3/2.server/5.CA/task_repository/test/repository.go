package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"go-kata/2.server/5.CA/task_repository/dto"
	_ "go-kata/2.server/5.CA/task_repository/storage"
)

//func NewSQLDB(dbConf config.DB) (*sqlx.DB, error) {
//
//}

//var db *sqlx.DB

func main() {
	const (
		connStr = "user=postgres password=123 dbname=postgres sslmode=disable"
	)
	//// Устанавливаем соединение с базой данных PostgreSQL
	//connStr := "user=postgres password=123 dbname=postgres sslmode=disable"
	//db, err := sql.Open("postgres", connStr)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//defer db.Close()
	//
	//// Вставляем данные в таблицу
	//_, err = db.Exec(`INSERT INTO vacancies (title, company) VALUES ($1, $2)`, "Alice", "25")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//fmt.Println("Данные успешно внесены в таблицу!")
	//user := storage.User{Name: "Mike", Email: "m@mail.com"}

	//user2 := storage.User{Name: "user345", Email: "4563@mail.com"}
	//user3 := storage.User{Name: "John", Email: "John@mail.com"}

	//store := storage.NewUserStorage()
	//////err := store.Create(user)
	////err := store.Create(user2)
	////err = store.Create(user3)
	////if err != nil {
	////	fmt.Println(err)
	////}
	//user, err := store.GetByID(4)
	//if err != nil {
	//	fmt.Println(err)
	//}
	//fmt.Println(user)
	//}

	//type User struct {
	//	ID    int    `json:"id"`
	//	Name  string `json:"name"`
	//	Email string `json:"email"`
	//	//Verified      bool   `json:"verified"`
	//}

	//func (s *UserStorage) List(limit, offset int) ([]dto.User, error) {

	limit, offset := 2, 4
	// Устанавливаем соединение с базой данных PostgreSQL
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		//return []dto.User{}, err
	}
	defer db.Close()
	// SQL-запрос для получения списка пользователей с учетом пагинации
	query := fmt.Sprintf(`SELECT * FROM users ORDER BY id LIMIT %d OFFSET %d`, limit, offset)

	// Выполнение запроса к базе данных
	rows, err := db.Query(query)
	if err != nil {
		fmt.Println(err)
		//return []dto.User{}, err
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
			fmt.Println(err)
			//return []dto.User{}, err
		}
		users = append(users, user)
	}
	fmt.Println(users, err)
}
