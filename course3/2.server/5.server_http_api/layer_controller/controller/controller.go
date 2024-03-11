package controller

import (
	"encoding/json"
	"go-kata/2.server/5.server_http_api/layer_controller/Handlers"
	"net/http"
)

type Responder interface {
	OutputJSON(w http.ResponseWriter, responseData interface{})

	ErrorUnauthorized(w http.ResponseWriter, err error)
	ErrorBadRequest(w http.ResponseWriter, err error)
	ErrorForbidden(w http.ResponseWriter, err error)
	ErrorInternal(w http.ResponseWriter, err error)
}

type UserController struct {
	responder Responder
	servicer  Handlers.Servicer
}

func NewUserController(responder Responder, servicer Handlers.Servicer) *UserController {
	return &UserController{
		responder: responder,
		servicer:  servicer,
	}
}

func (c *UserController) Register(w http.ResponseWriter, r *http.Request) {
	// Предварительная обработка запроса
	var requestBody RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		c.responder.ErrorBadRequest(w, err)
		return
	}

	// Вызов метода сервиса для регистрации пользователя
	if err := c.servicer.RegisterUser(requestBody.Login, requestBody.Password); err != nil {
		// Обработка ошибок сервиса
		c.responder.ErrorInternal(w, err)
		return
	}

	// Отправка успешного ответа клиенту
	c.responder.OutputJSON(w, "Пользователь зарегестрирован")
}

// Другие методы контроллера...
