package storage

import (
	"database/sql"
	"fmt"
	model "go-kata/2.server/5.CA/library/dto"
)

// GetAuthorsWithBooks возвращает список авторов с их книгами
func (s *LibraryStorage) GetAuthorsWithBooks() ([]model.Author, error) {
	//func (s *LibraryStorage) GetAuthorsWithBooks(db *sql.DB) ([]Author, error) {
	db, err := CreateTables()
	if err != nil {
		return []model.Author{}, err
	}
	defer db.Close()
	rows, err := db.Query(`
    SELECT
      a.id AS author_id,
      a.name AS author_name,
      b.id AS book_id,
      b.title AS book_title
    FROM
      Authors a
    LEFT JOIN
      Books b ON a.id = b.author_id
  `)
	if err != nil {
		return nil, fmt.Errorf("ошибка при выполнении запроса к базе данных: %v", err)
	}
	defer rows.Close()

	authorsMap := make(map[int]*model.Author)
	for rows.Next() {
		var authorID int
		var authorName string
		var bookID sql.NullInt64
		var bookTitle sql.NullString
		if err := rows.Scan(&authorID, &authorName, &bookID, &bookTitle); err != nil {
			return nil, fmt.Errorf("ошибка при сканировании строк запроса: %v", err)
		}

		author, ok := authorsMap[authorID]
		if !ok {
			author = &model.Author{
				ID:    authorID,
				Name:  authorName,
				Books: []model.Book{},
			}
			authorsMap[authorID] = author
		}

		if bookID.Valid && bookTitle.Valid {
			author.Books = append(author.Books, model.Book{
				ID:    int(bookID.Int64),
				Title: bookTitle.String,
			})
		}
	}

	authors := make([]model.Author, 0, len(authorsMap))
	for _, author := range authorsMap {
		authors = append(authors, *author)
	}
	return authors, nil
}

func (s *LibraryStorage) AddAuthor(name string) (string, error) {
	db, err := CreateTables()
	if err != nil {
		return "", err
	}
	defer db.Close()
	// Проверяем, существует ли уже автор с таким именем
	var exists bool
	err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM Authors WHERE name = $1)", name).Scan(&exists)
	if err != nil {
		return "", fmt.Errorf("ошибка при проверке существования автора: %v", err)
	}
	if exists {
		return fmt.Sprintf("автор с именем '%s' уже существует", name), nil
	}

	// Выполняем запрос к базе данных для добавления нового автора
	_, err = db.Exec("INSERT INTO Authors (name) VALUES ($1)", name)
	if err != nil {
		return "", fmt.Errorf("ошибка при добавлении автора: %v", err)
	}

	return "Ок", nil
}

// AuthorExists проверяет существование автора по его ID
func AuthorExists(db *sql.DB, authorID int) (bool, error) {
	var exists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM Authors WHERE id = $1)", authorID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("ошибка при проверке существования автора: %v", err)
	}
	return exists, nil
}
