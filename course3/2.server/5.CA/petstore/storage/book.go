package storage

//
//import (
//	"fmt"
//	//BookModel "go-kata/2.server/5.CA/library/dto/book"
//)
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
//
//// AddBook добавляет новую книгу с указанным ID автора
//func (s *LibraryStorage) AddBook(title string, authorID int) (string, error) {
//	db, err := CreateTables()
//	if err != nil {
//		return "", err
//	}
//	// Проверяем существование автора
//	exists, err := AuthorExists(db, authorID)
//	if err != nil {
//		return "", err
//	}
//	if !exists {
//		return fmt.Sprintf("Ошибка: автора с ID %d не существует", authorID), nil
//	}
//
//	// Выполняем запрос к базе данных для добавления новой книги
//	_, err = db.Exec("INSERT INTO Books (title, author_id) VALUES ($1, $2)", title, authorID)
//	if err != nil {
//		return "", fmt.Errorf("ошибка при добавлении книги: %v", err)
//	}
//
//	return "Ок", nil
//}
//
//// GetBooksWithAuthors возвращает список всех книг с информацией об авторах.
//func (s *LibraryStorage) GetBooksWithAuthors() ([]BookModel.Book, error) {
//	db, err := CreateTables()
//	if err != nil {
//		return []BookModel.Book{}, err
//	}
//	query := `
//    SELECT
//      b.id AS book_id,
//      b.title AS book_title,
//      a.id AS author_id,
//      a.name AS author_name
//    FROM
//      books b
//    LEFT JOIN
//      authors a ON b.author_id = a.id
//  `
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
