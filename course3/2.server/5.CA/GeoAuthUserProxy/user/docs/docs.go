// Package classification User.
//
// Документация по сервису User.
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
	"go-kata/2.server/5.CA/GeoAuthUserProxy/user/dto"
)

// Ошибка, неверный формат запроса.
// swagger:response BadRequestResponse
type BadRequestResponse struct {
	// Ошибка, неверный формат запроса.
}

// Внутренняя ошибка сервера
// swagger:response InternalServerErrorResponse
type InternalServerErrorResponse struct {
	// Внутренняя ошибка сервера
}

// swagger:route POST /api/user/list user ListUsers
// Вывод списка пользователей.
//
//  responses:
//   200: ListUsers
//   400: BadRequestResponse
//   403: Forbidden
//   500: InternalServerErrorResponse
//

// Pезультат запроса списка пользователей.
// swagger:response ListUsers
type ListUsers []dto.User

// swagger:route POST /api/user/profile user GetUserResponse
// Вывод пользователя по ID.
//
// parameters:
// +UserEmail - UserEmail
// in: query
// required: true
// name: UserEmail
// type: string
//
//  responses:
//   200: GetUserResponse
//   400: BadRequestResponse
//   403: Forbidden
//   500: InternalServerErrorResponse
//

// Вывод профиля пользователя.
// swagger:response GetUserResponse
type GetUserResponse string

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

// Токен авторизации не валиден.
// swagger:response Forbidden
type Forbidden struct{}

//go:generate swagger generate spec -o ../public/swagger.json --scan-models
