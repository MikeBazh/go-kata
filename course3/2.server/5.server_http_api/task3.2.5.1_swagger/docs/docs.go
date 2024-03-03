// Package classification геосервис.
//
// Документация по геосервису.
//
//	 Schemes: http
//	 BasePath: /
//	 Version: 1.0.0
//
//	 Consumes:
//	 - application/json
//
//	 Produces:
//	 - application/json
//
//	 Security:
//	 - basic
//
//	SecurityDefinitions:
//	basic:
//	  type: basic
//
// swagger:meta
package docs

import (
	"go-kata/2.server/5.server_http_api/task3.2.5.1_swagger/Dadata"
	"go-kata/2.server/5.server_http_api/task3.2.5.1_swagger/Handlers"
)

// swagger:route POST /api/address/search search SearchRequest
// Поиск по адресу.
//
//  responses:
//   200: SearchResponse

// swagger:parameters SearchRequest
type SearchRequest struct {
	// запрос для поиска по адресу.
	// in: body
	Body Handlers.SearchRequest
}

// Pезультат поиска по адресу.
// swagger:response SearchResponse
type SearchResponse struct {
	// in: body
	Body Dadata.SearchResponse
}

// swagger:route POST /api/address/geocode geocode GeocodeRequest
// Поиск по координатам.
//
//  responses:
//   200: GeocodeResponse

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
	Body Handlers.GeocodeRequest
}

//go:generate swagger generate spec -o ../public/swagger.json --scan-models
