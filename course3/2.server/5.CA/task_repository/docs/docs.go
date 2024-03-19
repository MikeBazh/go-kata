// Package classification геосервис.
//
// Документация по геосервису.
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
	"go-kata/2.server/5.CA/task_repository/dto"
)

// swagger:route POST /api/users/post request RequestUser
// Добавление нового пользователя.
//
//  responses:
//   200: RequestUserResponse
//   400: BadRequestResponse
//   500: InternalServerErrorResponse
//

// swagger:parameters RequestUser
type AddUser struct {
	// запрос для добавления пользователя.
	// in: body
	// required: true
	Body dto.RequestUser
}

// Pезультат запроса на добавление.
// swagger:response RequestUserResponse
type AddUserResponse string

// swagger:route POST /api/users/update update UpdateUser
// Обновление информации о пользователя.
//
//  responses:
//   200: UpdateUserResponse
//   400: BadRequestResponse
//   500: InternalServerErrorResponse
//

// swagger:parameters UpdateUser
type UpdateUser struct {
	// запрос для обновления пользователя (данные пользователя с именем name будут обновлены).
	// in: body
	// required: true
	Body dto.RequestUser
}

// Pезультат запроса на обновление.
// swagger:response UpdateUserResponse
type UpdateUserResponse string

// Ошибка, неверный формат запроса.
// swagger:response BadRequestResponse
type BadRequestResponseGeo struct {
	// Ошибка, неверный формат запроса.
}

// Внутренняя ошибка сервера 5.
// swagger:response InternalServerErrorResponse
type InternalServerErrorResponseGeo struct {
}

// swagger:route GET /api/users/{id} get GetUser
// Получение пользователя по ID.
//
// parameters:
// +id (path) - ID пользователя
// in: path
// required: true
// name: id
// type: integer
//
//  responses:
//   200: GetUserResponse
//   400: BadRequestResponse
//   500: InternalServerErrorResponse
//

// Pезультат запроса по ID.
// swagger:response GetUserResponse
type GetUserResponse struct {
	// in: body
	// required: true
	Body dto.User
}

// swagger:route DELETE /api/users/{id} get DeleteUser
// Удаление пользователя по ID.
//
// parameters:
// +id (path) - ID пользователя
// in: path
// required: true
// name: id
// type: integer
//
//  responses:
//   200: DeleteUserResponse
//   400: BadRequestResponse
//   500: InternalServerErrorResponse
//

// Pезультат запроса по ID.
// swagger:response DeleteUserResponse
type DeleteUserResponse string

// swagger:route GET /api/users list ListUsers
// Вывод списка пользователей.
//
// parameters:
// +limit - Limit
// in: query
// required: true
// name: limit
// type: integer
//
// +offset (path) - Offset
// in: query
// required: true
// name: offset
// type: integer
//
//  responses:
//   200: RequestUserResponse
//   400: BadRequestResponse
//   500: InternalServerErrorResponse
//

// Pезультат запроса списка пользователей.
// swagger:response ListUsers
type ListUsers []dto.User

//go:generate swagger generate spec -o ../public/swagger.json --scan-models
