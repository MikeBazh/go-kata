// Package classification библиотека.
//
// Документация по сервису библиотеки.
//
//	Schemes: http
//	BasePath: /
//	Version: 1.0.0
//
//	Consumes:
//	- application/json
//
//	Produces:
//	- application/json
//
//	Security:
//	- basic
//
// swagger:meta
package docs

import (
	model "go-kata/2.server/5.CA/library/dto"
	BookModel "go-kata/2.server/5.CA/library/dto/book"
	UserModel "go-kata/2.server/5.CA/library/dto/user"
)

// Ошибка, неверный формат запроса.
// swagger:response BadRequestResponse
type BadRequestResponse struct {
	// Ошибка, неверный формат запроса.
}

// Внутренняя ошибка сервера: сервис https://dadata.ru не доступен.
// swagger:response InternalServerErrorResponse
type InternalServerErrorResponse struct {
	// Внутренняя ошибка сервера: сервис https://dadata.ru не доступен.
}

// Пользователь не авторизирован.
// swagger:response Unauthorized

// Токен авторизации не валиден.
// swagger:response Forbidden

// swagger:route POST /api/book/take book BookTakeResponse
// Получение книги пользователем.
//
// parameters:
// +userID - userID
// in: query
// required: true
// name: userID
// type: string
//
// +bookID (path) - bookID
// in: query
// required: true
// name: bookID
// type: string
//
//  responses:
//   200: BookTakeResponse
//   400: BadRequestResponse
//   500: InternalServerErrorResponse
//

// Pезультат запроса получения книги пользователем ("ОК" или ошибка).
// swagger:response BookTakeResponse
type BookTakeResponse string

// swagger:route POST /api/book/return book ReturnBookResponse
// Возврат книги пользователем.
//
// parameters:
// +userID - userID
// in: query
// required: true
// name: userID
// type: string
//
// +bookID (path) - bookID
// in: query
// required: true
// name: bookID
// type: string
//
//  responses:
//   200: BookReturnResponse
//   400: BadRequestResponse
//   500: InternalServerErrorResponse
//

// Pезультат запроса возврата книги пользователем ("ОК" или ошибка).
// swagger:response BookReturnResponse
type BookReturnResponse string

// swagger:route POST /api/book/add book AddBookResponse
//  Добавление книги.
//
// parameters:
// +title - title
// in: query
// required: true
// name: title
// type: string
//
// +authorID (path) - authorID
// in: query
// required: true
// name: authorID
// type: string
//
//  responses:
//   200: AddBookResponse
//   400: BadRequestResponse
//   500: InternalServerErrorResponse
//

// Pезультат запроса на добавление книги ("ОК" или ошибка).
// swagger:response AddBookResponse
type AddBookResponse string

// swagger:route GET /api/authors authors ListAuthors
// Вывод списка авторов.
//
//  responses:
//   200: ListAuthors
//   400: BadRequestResponse
//   500: InternalServerErrorResponse
//

// Pезультат запроса списка авторов.
// swagger:response ListAuthors
type ListAuthors []model.Author

// swagger:route GET /api/users users ListUsers
// Вывод списка пользователей.
//
//  responses:
//   200: ListUsers
//   400: BadRequestResponse
//   500: InternalServerErrorResponse
//

// Pезультат запроса списка авторов.
// swagger:response ListUsers
type ListUsers []UserModel.User

// swagger:route GET /api/books book ListBooks
// Вывод списка книг.
//
//  responses:
//   200: ListBooks
//   400: BadRequestResponse
//   500: InternalServerErrorResponse
//

// Pезультат запроса списка книг.
// swagger:response ListBooks
type ListBooks []BookModel.Book

// swagger:route POST /api/authors authors AddAuthorsResponse
// Добавление автора.
//
// parameters:
// +AuthorName - Author's Name
// in: query
// required: true
// name: AuthorName
// type: string
//
//  responses:
//   200: AddAuthorsResponse
//   400: BadRequestResponse
//   500: InternalServerErrorResponse
//

// Pезультат запроса добавления автора ("ОК" или ошибка).
// swagger:response AddAuthorsResponse
type AddAuthorsResponse string

// Ошибка, неверный формат запроса.
// swagger:response BadRequestResponse
type BadRequestResponseGeo struct {
	// Ошибка, неверный формат запроса.
}

// Внутренняя ошибка сервера.
// swagger:response InternalServerErrorResponse
type InternalServerErrorResponseGeo struct {
	// Внутренняя ошибка сервера.
}

//go:generate swagger generate spec -o ../public/swagger.json --scan-models
