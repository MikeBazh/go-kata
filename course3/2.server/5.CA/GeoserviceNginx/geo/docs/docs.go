// Package classification Geo.
//
// Документация по Geo.
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
//	type: basic
// bearerAuth:
//	type: apiKey
//	in: header
//	name: Authorization
// swagger:meta

package docs

import (
	"go-kata/2.server/5.CA/GeoserviceNginx/geo/Dadata"
	"go-kata/2.server/5.CA/GeoserviceNginx/geo/controller"
)

// swagger:route POST /api/address/search geo SearchRequest
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

// swagger:route POST /api/address/geocode geo GeocodeRequest
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

// Токен авторизации не валиден.
// swagger:response Forbidden
type Forbidden struct{}

//go:generate swagger generate spec -o ../public/swagger.json --scan-models
