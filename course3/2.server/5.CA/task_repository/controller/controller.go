package controller

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi"
	"go-kata/2.server/5.CA/task_repository/dto"
	"go-kata/2.server/5.CA/task_repository/services"
	"html/template"
	"net/http"
	"strconv"
	"time"
)

type Responder interface {
	OutputJSON(w http.ResponseWriter, responseData interface{})

	//ErrorUnauthorized(w http.ResponseWriter, err error)
	ErrorBadRequest(w http.ResponseWriter, err error)
	ErrorForbidden(w http.ResponseWriter, err error)
	ErrorInternal(w http.ResponseWriter, err error)
}

type UserController struct {
	responder Responder
	servicer  services.Servicer
}

func NewUserController(responder Responder, servicer services.Servicer) *UserController {
	return &UserController{
		responder: responder,
		servicer:  servicer,
	}
}

func (c *UserController) SwaggerUI(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	tmpl, err := template.New("swagger").Parse(swaggerTemplate)
	if err != nil {
		return
	}
	err = tmpl.Execute(w, struct {
		Time int64
	}{
		Time: time.Now().Unix(),
	})
	if err != nil {
		return
	}
}

func (c *UserController) AddUser(w http.ResponseWriter, r *http.Request) {
	// Предварительная обработка запроса
	var requestBody dto.RequestUser
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		c.responder.ErrorBadRequest(w, err)
		return
	}

	// Вызов метода сервиса для регистрации пользователя
	err := c.servicer.AddUser(requestBody)
	if err != nil {
		// Обработка ошибок сервиса
		c.responder.ErrorInternal(w, err)
		return
	}
	// Отправка успешного ответа клиенту
	c.responder.OutputJSON(w, "пользователь добавлен")
}

func (c *UserController) UpdateUser(w http.ResponseWriter, r *http.Request) {
	// Предварительная обработка запроса
	var requestBody dto.RequestUser
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		c.responder.ErrorBadRequest(w, err)
		return
	}

	// Вызов метода сервиса для регистрации пользователя
	err := c.servicer.UpdateUser(requestBody)
	if err != nil {
		// Обработка ошибок сервиса
		c.responder.ErrorInternal(w, err)
		return
	}
	// Отправка успешного ответа клиенту
	c.responder.OutputJSON(w, "данные обновлены")
}

func (c *UserController) GetUserByID(w http.ResponseWriter, r *http.Request) {
	// Извлекаем ID из URL пути
	idStr := chi.URLParam(r, "id")
	userID, err := strconv.Atoi(idStr)
	if err != nil {
		c.responder.ErrorBadRequest(w, err)
		return
	}
	//fmt.Println("ID: ", userID)
	// Вызов метода сервиса для регистрации пользователя
	user, err := c.servicer.GetUserByID(userID)
	if err != nil {
		c.responder.ErrorInternal(w, err)
		return
	}
	// Отправка успешного ответа клиенту
	c.responder.OutputJSON(w, user)
}

func (c *UserController) DeleteUserByID(w http.ResponseWriter, r *http.Request) {
	// Извлекаем ID из URL пути
	idStr := chi.URLParam(r, "id")
	userID, err := strconv.Atoi(idStr)
	if err != nil {
		c.responder.ErrorBadRequest(w, err)
		return
	}
	//fmt.Println("ID: ", userID)
	// Вызов метода сервиса для регистрации пользователя
	err = c.servicer.DeleteByID(userID)
	if err != nil {
		//fmt.Println("пользователь был удален ранее или не существует")
		c.responder.OutputJSON(w, "пользователь был удален ранее или не существует")
		return
	}
	// Отправка успешного ответа клиенту
	c.responder.OutputJSON(w, "пользователь c ID "+idStr+" удален")
}

func (c *UserController) List(w http.ResponseWriter, r *http.Request) {
	// Извлекаем ID из URL
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		c.responder.ErrorBadRequest(w, err)
		return
	}
	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		c.responder.ErrorBadRequest(w, err)
		return
	}
	// Вызов метода сервиса для вывода пользователей
	users, err := c.servicer.List(limit, offset)
	if err != nil {
		fmt.Println("неверные параметры List")
		c.responder.ErrorBadRequest(w, err)
		return
	}
	// Отправка успешного ответа клиенту
	c.responder.OutputJSON(w, users)
}

const (
	swaggerTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <meta http-equiv="X-UA-Compatible" content="ie=edge">
    <script src="//unpkg.com/swagger-ui-dist@3/swagger-ui-standalone-preset.js"></script>
    <!-- <script src="https://cdnjs.cloudflare.com/ajax/libs/swagger-ui/3.22.1/swagger-ui-standalone-preset.js"></script> -->
    <script src="//unpkg.com/swagger-ui-dist@3/swagger-ui-bundle.js"></script>
    <!-- <script src="https://cdnjs.cloudflare.com/ajax/libs/swagger-ui/3.22.1/swagger-ui-bundle.js"></script> -->
    <link rel="stylesheet" href="//unpkg.com/swagger-ui-dist@3/swagger-ui.css" />
    <!-- <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/swagger-ui/3.22.1/swagger-ui.css" /> -->
	<style>
		body {
			margin: 0;
		}
	</style>
    <title>Swagger</title>
</head>
<body>
    <div id="swagger-ui"></div>
    <script>
        window.onload = function() {
          SwaggerUIBundle({
            url: "/public/swagger.json?{{.Time}}",
            dom_id: '#swagger-ui',
            presets: [
              SwaggerUIBundle.presets.apis,
              SwaggerUIStandalonePreset
            ],
            layout: "StandaloneLayout"
          })
        }
    </script>
</body>
</html>
`
)
