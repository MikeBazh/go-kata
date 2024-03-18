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
// SecurityDefinitions:
// basic:
//
//	type: basic
//
// bearerAuth:
//
//	type: apiKey
//	in: header
//	name: Authorization
//
// swagger:meta
package docs

import (
	"go-kata/2.server/5.server_http_api/layer-service/Dadata"
	"go-kata/2.server/5.server_http_api/layer-service/internal/controller"
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

// swagger:route POST /api/address/search search SearchRequest
// Поиск по адресу.
//
//  responses:
//   200: SearchResponse
//   400: BadRequestResponse
//   403: Forbidden
//   500: InternalServerErrorResponse
//
//  security:
//  - bearerAuth: []

// swagger:parameters SearchRequest
type SearchRequest struct {
	// запрос для поиска по адресу.
	// in: body
	// required: true
	// swagger:allOf
	// Items:
	//    $ref: "#/definitions/AuthToken"
	Body controller.SearchRequest
}

// Pезультат поиска по адресу.
// swagger:response SearchResponse
type SearchResponse struct {
	// in: body
	// required: true
	// swagger:allOf
	// Items:
	//    $ref: "#/definitions/AuthToken"
	Body Dadata.SearchResponse
}

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

// swagger:route POST /api/address/geocode geocode GeocodeRequest
// Поиск по координатам.
//
//  responses:
//   200: GeocodeResponse
//   400: BadRequestResponse
//   403: Forbidden
//   500: InternalServerErrorResponse
//
//  security:
//  - bearerAuth: []

// Pезультат поиска по координатам.
// swagger:response GeocodeResponse
type GeocodeResponse struct {
	// in: body
	Body Dadata.GeocodeResponse
}

// swagger:parameters GeocodeRequest
type GeocodeRequest struct {
	// запрос для поиска по координатам.
	// in: body
	Body controller.GeocodeRequest
}

// Ошибка, неверный формат запроса.
// swagger:response BadRequestResponse
type BadRequestResponseGeo struct {
	// Ошибка, неверный формат запроса.
}

// Внутренняя ошибка сервера: сервис https://dadata.ru не доступен.
// swagger:response InternalServerErrorResponse
type InternalServerErrorResponseGeo struct {
	// Внутренняя ошибка сервера: сервис https://dadata.ru не доступен.
}

// Пользователь не авторизирован.
// swagger:response Unauthorized
type Unauthorized struct {
}

// Токен авторизации не валиден.
// swagger:response Forbidden
type Forbidden struct{}

// swagger:model AuthToken
type AuthToken struct {
	// Токен для аутентификации пользователя.
	// in: header
	// name: Authorization
	// required: true
	Token string `json:"token"`
}

//go:generate swagger generate spec -o ../public/swagger.json --scan-models
