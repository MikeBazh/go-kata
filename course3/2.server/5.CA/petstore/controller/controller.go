package controller

import (
	"encoding/json"
	"fmt"
	UserModel "go-kata/2.server/5.CA/petstore/dto/user"
	"go-kata/2.server/5.CA/petstore/services"
	"html/template"
	"net/http"
	"strings"
	"time"
)

type Responder interface {
	OutputJSON(w http.ResponseWriter, responseData interface{})

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

func (c *UserController) CreateUser(w http.ResponseWriter, r *http.Request) {
	var requestBody UserModel.User
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		c.responder.ErrorBadRequest(w, err)
		return
	}
	fmt.Println(requestBody)
	// Вызов метода сервиса для регистрации пользователя
	err := c.servicer.CreateUser(requestBody)
	if err != nil {
		// Обработка ошибок сервиса
		c.responder.ErrorInternal(w, err)
		return
	}
	// Отправка успешного ответа клиенту
	c.responder.OutputJSON(w, "Ok")
}

func (c *UserController) GetUserByName(w http.ResponseWriter, r *http.Request) {
	//Извлекаем параметры из URL
	path := r.URL.Path
	parts := strings.Split(path, "/")
	username := parts[len(parts)-1]
	fmt.Println(username)
	//username := r.URL.Query().Get("username")
	// Вызов метода сервиса для регистрации пользователя
	respond, err := c.servicer.GetUserByName(username)
	if err != nil {
		// Обработка ошибок сервиса
		c.responder.ErrorInternal(w, err)
		return
	}
	// Отправка успешного ответа клиенту
	c.responder.OutputJSON(w, respond)
}

func (c *UserController) UpdateUserByName(w http.ResponseWriter, r *http.Request) {
	//Извлекаем параметры из URL
	path := r.URL.Path
	parts := strings.Split(path, "/")
	username := parts[len(parts)-1]
	//username := r.URL.Query().Get("username")
	fmt.Println(username)
	var requestBody UserModel.User
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		c.responder.ErrorBadRequest(w, err)
		return
	}
	fmt.Println(requestBody)
	// Вызов метода сервиса для регистрации пользователя
	respond, err := c.servicer.UpdateUserByName(username, requestBody)
	if err != nil {
		// Обработка ошибок сервиса
		c.responder.ErrorInternal(w, err)
		return
	}
	// Отправка успешного ответа клиенту
	c.responder.OutputJSON(w, respond)
}

func (c *UserController) DeleteUserByName(w http.ResponseWriter, r *http.Request) {
	//Извлекаем параметры из URL
	path := r.URL.Path
	parts := strings.Split(path, "/")
	username := parts[len(parts)-1]
	fmt.Println(username)
	// Вызов метода сервиса для регистрации пользователя
	respond, err := c.servicer.DeleteUserByName(username)
	if err != nil {
		// Обработка ошибок сервиса
		c.responder.ErrorInternal(w, err)
		return
	}
	// Отправка успешного ответа клиенту
	c.responder.OutputJSON(w, respond)
}

func (c *UserController) LoginUser(w http.ResponseWriter, r *http.Request) {
	//Извлекаем параметры из URL
	username := r.URL.Query().Get("username")
	password := r.URL.Query().Get("password")
	fmt.Println(username, password)
	// Вызов метода сервиса для регистрации пользователя
	err := c.servicer.LoginUser(username, password)
	if err != nil {
		// Обработка ошибок сервиса
		c.responder.ErrorInternal(w, err)
		return
	}
	// Отправка успешного ответа клиенту
	c.responder.OutputJSON(w, "Успешный вход")
}

func (c *UserController) LogoutUser(w http.ResponseWriter, r *http.Request) {
	//Извлекаем параметры из URL
	username := r.URL.Query().Get("username")
	//password := r.URL.Query().Get("password")
	fmt.Println(username)
	// Вызов метода сервиса для регистрации пользователя
	err := c.servicer.LogoutUser(username)
	if err != nil {
		// Обработка ошибок сервиса
		c.responder.ErrorInternal(w, err)
		return
	}
	// Отправка успешного ответа клиенту
	c.responder.OutputJSON(w, "Пользователь удален")
}

func (c *UserController) CreateWithArray(w http.ResponseWriter, r *http.Request) {
	//Извлекаем параметры из URL
	//username := r.URL.Query().Get("username")
	var requestBody []UserModel.User
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		c.responder.ErrorBadRequest(w, err)
		return
	}
	fmt.Println(requestBody)
	// Вызов метода сервиса для регистрации пользователя
	err := c.servicer.CreateWithArray(requestBody)
	if err != nil {
		// Обработка ошибок сервиса
		c.responder.ErrorInternal(w, err)
		return
	}
	// Отправка успешного ответа клиенту
	c.responder.OutputJSON(w, "Ok")
}

func (c *UserController) CreateWithList(w http.ResponseWriter, r *http.Request) {
	//Извлекаем параметры из URL
	//username := r.URL.Query().Get("username")
	var requestBody []UserModel.User
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		c.responder.ErrorBadRequest(w, err)
		return
	}
	fmt.Println(requestBody)
	// Вызов метода сервиса для регистрации пользователя
	err := c.servicer.CreateWithList(requestBody)
	if err != nil {
		// Обработка ошибок сервиса
		c.responder.ErrorInternal(w, err)
		return
	}
	// Отправка успешного ответа клиенту
	c.responder.OutputJSON(w, "Ok")
}

//func (c *UserController) GetAuthorsWithBooks(w http.ResponseWriter, r *http.Request) {
//
//	// Вызов метода сервиса для регистрации пользователя
//	respond, err := c.servicer.GetAuthorsWithBooks()
//	if err != nil {
//		// Обработка ошибок сервиса
//		c.responder.ErrorInternal(w, err)
//		return
//	}
//	// Отправка успешного ответа клиенту
//	c.responder.OutputJSON(w, respond)
//}
//
//func (c *UserController) GetUsersWithRentedBooks(w http.ResponseWriter, r *http.Request) {
//	// Вызов метода сервиса для регистрации пользователя
//	respond, err := c.servicer.GetUsersWithRentedBooks()
//	if err != nil {
//		// Обработка ошибок сервиса
//		c.responder.ErrorInternal(w, err)
//		return
//	}
//	// Отправка успешного ответа клиенту
//	c.responder.OutputJSON(w, respond)
//}
//
//func (c *UserController) BookTake(w http.ResponseWriter, r *http.Request) {
//	// Извлекаем параметры из URL
//	userIdStr := r.URL.Query().Get("userID")
//	bookIdStr := r.URL.Query().Get("bookID")
//	userID, err := strconv.Atoi(userIdStr)
//	if err != nil {
//		c.responder.ErrorBadRequest(w, err)
//		return
//	}
//	bookID, err := strconv.Atoi(bookIdStr)
//	if err != nil {
//		c.responder.ErrorBadRequest(w, err)
//		return
//	}
//	// Вызов метода сервиса для регистрации пользователя
//	response, err := c.servicer.BookTake(userID, bookID)
//	if err != nil {
//		// Обработка ошибок сервиса
//		c.responder.ErrorInternal(w, err)
//		return
//	}
//	// Отправка успешного ответа клиенту
//	c.responder.OutputJSON(w, response)
//}
//
//func (c *UserController) ReturnBook(w http.ResponseWriter, r *http.Request) {
//	// Извлекаем параметры из URL
//	userIdStr := r.URL.Query().Get("userID")
//	bookIdStr := r.URL.Query().Get("bookID")
//	userID, err := strconv.Atoi(userIdStr)
//	if err != nil {
//		c.responder.ErrorBadRequest(w, err)
//		return
//	}
//	bookID, err := strconv.Atoi(bookIdStr)
//	if err != nil {
//		c.responder.ErrorBadRequest(w, err)
//		return
//	}
//	// Вызов метода сервиса для регистрации пользователя
//	err = c.servicer.ReturnBook(userID, bookID)
//	if err != nil {
//		// Обработка ошибок сервиса
//		c.responder.ErrorInternal(w, err)
//		return
//	}
//	// Отправка успешного ответа клиенту
//	c.responder.OutputJSON(w, "Ok")
//}
//
//func (c *UserController) AddAuthor(w http.ResponseWriter, r *http.Request) {
//	// Извлекаем параметры из URL
//	authorName := r.URL.Query().Get("AuthorName")
//	// Вызов метода сервиса для регистрации пользователя
//	respond, err := c.servicer.AddAuthor(authorName)
//	if err != nil {
//		// Обработка ошибок сервиса
//		c.responder.ErrorInternal(w, err)
//		return
//	}
//	// Отправка успешного ответа клиенту
//	c.responder.OutputJSON(w, respond)
//}
//
//func (c *UserController) AddBook(w http.ResponseWriter, r *http.Request) {
//	// Извлекаем параметры из URL
//	title := r.URL.Query().Get("title")
//	authorIdStr := r.URL.Query().Get("authorID")
//	authorID, err := strconv.Atoi(authorIdStr)
//	if err != nil {
//		c.responder.ErrorBadRequest(w, err)
//		return
//	}
//	fmt.Println("добавление книги:", title, authorID)
//	// Вызов метода сервиса для регистрации пользователя
//	respond, err := c.servicer.AddBook(title, authorID)
//	if err != nil {
//		// Обработка ошибок сервиса
//		c.responder.ErrorInternal(w, err)
//		return
//	}
//	// Отправка успешного ответа клиенту
//	c.responder.OutputJSON(w, respond)
//}
//
//func (c *UserController) GetBooks(w http.ResponseWriter, r *http.Request) {
//	// Вызов метода сервиса для вывода списка книг
//	respond, err := c.servicer.GetBooks()
//	if err != nil {
//		// Обработка ошибок сервиса
//		c.responder.ErrorInternal(w, err)
//		return
//	}
//	// Отправка успешного ответа клиенту
//	c.responder.OutputJSON(w, respond)
//}

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
