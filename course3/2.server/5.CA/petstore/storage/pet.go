package storage

import (
	"encoding/json"
	"fmt"
	PetModel "go-kata/2.server/5.CA/petstore/dto/pet"
)

//
//// GetUserBook позволяет пользователю получить книгу по ее ID
//func (s *LibraryStorage) GetUserBook(userID, bookID int) (string, error) {
//	db, err := CreateTables()
//	if err != nil {
//		return "", err
//	}
//	defer db.Close()
//	// Проверяем, существует ли пользователь с указанным ID
//	var userExists bool
//	err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM Users WHERE id = $1)", userID).Scan(&userExists)
//	if err != nil {
//		return "", fmt.Errorf("ошибка при проверке существования пользователя: %v", err)
//	}
//	if !userExists {
//		return fmt.Sprintf("Ошибка: пользователь с ID %d не найден", userID), nil
//	}
//
//	// Проверяем, существует ли книга с указанным ID и не арендована ли она уже
//	var bookExists, bookAvailable bool
//	err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM Books WHERE id = $1)", bookID).Scan(&bookExists)
//	if err != nil {
//		return "", fmt.Errorf("ошибка при проверке существования книги: %v", err)
//	}
//	if !bookExists {
//		return fmt.Sprintf("Ошибка: книга с ID %d не найдена", bookID), nil
//	}
//
//	err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM Books WHERE id = $1 AND rented_by IS NULL)", bookID).Scan(&bookAvailable)
//	if err != nil {
//		return "", fmt.Errorf("ошибка при проверке существования книги: %v", err)
//	}
//	if !bookAvailable {
//		return fmt.Sprintf("Ошибка: книга с ID %d уже арендована", bookID), nil
//	}
//
//	// Обновляем запись о книге в базе данных, указывая ID пользователя, который ее получил
//	_, err = db.Exec("UPDATE Books SET rented_by = $1 WHERE id = $2", userID, bookID)
//	if err != nil {
//		return "", fmt.Errorf("ошибка при обновлении записи о книге: %v", err)
//	}
//	// Обновляем столбец rented_books в таблице Users
//	_, err = db.Exec("UPDATE Users SET rented_books = array_append(rented_books, $1) WHERE id = $2", bookID, userID)
//	if err != nil {
//		return "", fmt.Errorf("ошибка при обновлении столбца rented_books: %v", err)
//	}
//	return "Ок", nil
//}
//
//// ReturnBook позволяет пользователю вернуть книгу по ее ID
//func (s *LibraryStorage) ReturnBook(userID, bookID int) error {
//	db, err := CreateTables()
//	if err != nil {
//		return err
//	}
//	defer db.Close()
//	// Проверяем, существует ли пользователь с указанным ID
//	var userExists bool
//	err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM Users WHERE id = $1)", userID).Scan(&userExists)
//	if err != nil {
//		return fmt.Errorf("ошибка при проверке существования пользователя: %v", err)
//	}
//	if !userExists {
//		return fmt.Errorf("пользователь с ID %d не найден", userID)
//	}
//
//	// Проверяем, арендована ли книга с указанным ID пользователем
//	var bookRented bool
//	err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM Books WHERE id = $1 AND rented_by = $2)", bookID, userID).Scan(&bookRented)
//	if err != nil {
//		return fmt.Errorf("ошибка при проверке арендованной книги: %v", err)
//	}
//	if !bookRented {
//		return fmt.Errorf("книга с ID %d не арендована пользователем с ID %d", bookID, userID)
//	}
//
//	// Обновляем запись о книге в базе данных, убирая информацию об аренде
//	_, err = db.Exec("UPDATE Books SET rented_by = NULL WHERE id = $1", bookID)
//	if err != nil {
//		return fmt.Errorf("ошибка при обновлении записи о книге: %v", err)
//	}
//
//	// Обновляем столбец rented_books в таблице Users, убирая ID возвращенной книги
//	_, err = db.Exec("UPDATE users SET rented_books = array_remove(rented_books, $1) WHERE id = $2", bookID, userID)
//	if err != nil {
//		return fmt.Errorf("ошибка при обновлении столбца rented_books: %v", err)
//	}
//
//	return nil
//}

func (s *LibraryStorage) AddPet(pet PetModel.Pet) error {
	db, err := CreateTables()
	if err != nil {
		return err
	}
	// Выполняем запрос к базе данных для добавления
	//var tags []byte
	tags, err := json.Marshal(pet.Tags)
	category, err := json.Marshal(pet.Category)
	if err != nil {
		return err
	}
	_, err = db.Exec("INSERT INTO pets (name, status, tags, category) VALUES ($1, $2, $3, $4)", pet.Name, pet.Status, tags, category)
	if err != nil {
		return fmt.Errorf("ошибка при добавлении pet: %v", err)
	}
	return nil
}

func (s *LibraryStorage) UpdatePet(pet PetModel.Pet) error {
	db, err := CreateTables()
	if err != nil {
		return err
	}
	if err != nil {
		return err
	}
	//tags, err := json.Marshal(pet.Tags)
	//category, err := json.Marshal(pet.Category)
	// Выполняем запрос к базе данных для обновления
	_, err = db.Exec("UPDATE pets SET name=$1, status=$2 WHERE id=$3", pet.Name, pet.Status, pet.Id)
	//fmt.Println(id)
	if err != nil {
		return fmt.Errorf("ошибка при обновлении pet: %v", err)
	}
	return nil
}

func (s *LibraryStorage) FindPetById(id int) (pet PetModel.Pet, err error) {
	db, err := CreateTables()
	if err != nil {
		return PetModel.Pet{}, err
	}
	// Выполняем запрос к базе данных для обновления
	row, err := db.Query("SELECT * FROM pets WHERE id=$1", id)
	var tags Tags
	var category Category
	var tagsJson []byte
	var categoryJson []byte
	for row.Next() {
		err = row.Scan(&pet.Id, &pet.Name, &pet.Status, &tagsJson, &categoryJson)
		//fmt.Println("1", err)
		err = json.Unmarshal(tagsJson, &tags)
		//fmt.Println("2", err)
		err = json.Unmarshal(categoryJson, &category)
		//fmt.Println("3", err)
		if err != nil {
			return PetModel.Pet{}, fmt.Errorf("ошибка при сканировании строк запроса: %v", err)
		}
		pet.Tags = tags
		pet.Category = category
		fmt.Println(pet)
	}
	return pet, nil
}

func (s *LibraryStorage) DeletePet(id int) error {
	db, err := CreateTables()
	if err != nil {
		return err
	}
	// Выполняем запрос к базе данных для обновления
	var ID int
	query := "DELETE FROM pets WHERE id = $1 RETURNING id"
	row := db.QueryRow(query, id)
	err = row.Scan(&ID)
	if err != nil {
		return err
	}
	if ID == 0 {
		return fmt.Errorf("not found")
	}
	fmt.Println("pet deleted")
	return nil
}

func (s *LibraryStorage) FindPetByStatus(status string) (pets []PetModel.Pet, err error) {
	db, err := CreateTables()
	if err != nil {
		return []PetModel.Pet{}, err
	}
	var pet PetModel.Pet
	// Выполняем запрос к базе данных для обновления
	//db.Query("SELECT * FROM pets WHERE status=$1", pet.Name)
	rows, err := db.Query("SELECT * FROM pets WHERE status=$1", status)
	if err != nil {
		return []PetModel.Pet{}, fmt.Errorf("ошибка при запросе pet: %v", err)
	}
	for rows.Next() {
		var tags Tags
		var category Category
		var tagsJson []byte
		var categoryJson []byte
		err := rows.Scan(&pet.Id, &pet.Name, &pet.Status, &tagsJson, &categoryJson)
		//fmt.Println("1", err)
		err = json.Unmarshal(tagsJson, &tags)
		//fmt.Println("2", err)
		err = json.Unmarshal(categoryJson, &category)
		//fmt.Println("3", err)
		if err != nil {
			return []PetModel.Pet{}, fmt.Errorf("ошибка при сканировании строк запроса: %v", err)
		}
		pet.Tags = tags
		pet.Category = category
		fmt.Println(pet)
		pets = append(pets, pet)
	}
	return pets, nil
}

type Tags []struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type Category struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

//query := "UPDATE users SET firstname = $1, lastname=$2, email = $3, phone = $4, userStatus=$5, password=$6  WHERE username = $7 RETURNING id"
//var id int
//err = db.QueryRow(query, newUser.FirstName, newUser.LastName, newUser.Email, newUser.Phone, newUser.UserStatus, newUser.Password, name).Scan(&id)
//if err != nil {
//return UserModel.User{}, err
//}
//fmt.Println("user", id, "updated")

//// GetBooksWithAuthors возвращает список всех книг с информацией об авторах.
//func (s *LibraryStorage) GetBooksWithAuthors() ([]BookModel.Book, error) {
//	db, err := CreateTables()
//	if err != nil {
//		return []BookModel.Book{}, err
//	}
//	query := `
//   SELECT
//     b.id AS book_id,
//     b.title AS book_title,
//     a.id AS author_id,
//     a.name AS author_name
//   FROM
//     books b
//   LEFT JOIN
//     authors a ON b.author_id = a.id
// `
//	rows, err := db.Query(query)
//	if err != nil {
//		return nil, fmt.Errorf("ошибка при выполнении запроса: %v", err)
//	}
//	defer rows.Close()
//
//	var books []BookModel.Book
//	for rows.Next() {
//		var book BookModel.Book
//		var author BookModel.Author
//		err := rows.Scan(&book.ID, &book.Title, &author.ID, &author.Name)
//		if err != nil {
//			return nil, fmt.Errorf("ошибка при сканировании строк запроса: %v", err)
//		}
//		book.Author = author
//		books = append(books, book)
//	}
//
//	return books, nil
//}
