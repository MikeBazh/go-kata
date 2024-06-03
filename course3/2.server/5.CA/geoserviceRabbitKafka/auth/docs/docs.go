// Package classification Auth.
//
// Документация по Auth.
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
// SecurityDefinitions:
// basic:
//
// swagger:meta
package docs

import (
	"go-kata/2.server/5.CA/GeoAuthUserProxy/auth/controller"
)

// swagger:route POST /api/register register RegisterRequest
// Регистрация нового пользователя.
//
//  responses:
//   200: RegisterResponse

// swagger:parameters RegisterRequest
type RegisterRequest struct {
	// запрос для регистрации пользователя.
	// in: body
	Body controller.RegisterRequest
}

// успешная регистрация нового пользователя.
// swagger:response RegisterResponse
// Ответ на запрос регистрации нового пользователя.
type RegisterResponse string

// swagger:route POST /api/login login LoginRequest
// Аутентификация пользователя.
//
//  responses:
//   200: LoginResponse

// swagger:parameters LoginRequest
type LoginRequest struct {
	// запрос для аутентификации пользователя.
	// in: body
	Body controller.RegisterRequest
}

// получен токен авторизации (в формате: Bearer + токен) или ошибка входа.
// swagger:response LoginResponse
// Ответ на запрос аутентификации пользователя.
type LoginResponse string

// Внутренняя ошибка сервера.
// swagger:response InternalServerErrorResponse
type InternalServerErrorResponseGeo struct {
	// Внутренняя ошибка сервера: сервис https://dadata.ru не доступен.
}

// swagger:model AuthToken
type AuthToken struct {
	// Токен для аутентификации пользователя.
	// in: header
	// name: Authorization
	// required: true
	Token string `json:"token"`
}

//go:generate swagger generate spec -o ../public/swagger.json --scan-models
