package controller

import (
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	_ "github.com/dgrijalva/jwt-go"
	orderModel "go-kata/2.server/5.CA/petstore/dto/order"
	PetModel "go-kata/2.server/5.CA/petstore/dto/pet"
	_ "go-kata/2.server/5.CA/petstore/dto/pet"
	UserModel "go-kata/2.server/5.CA/petstore/dto/user"
	"go-kata/2.server/5.CA/petstore/services"
	"html/template"
	"net/http"
	"strconv"
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
	_, err := c.servicer.DeleteUserByName(username)
	if err != nil {
		// Обработка ошибок сервиса
		c.responder.ErrorInternal(w, err)
		return
	}
	// Отправка успешного ответа клиенту
	//c.responder.OutputJSON(w, respond)
}

func (c *UserController) LoginUser(w http.ResponseWriter, r *http.Request) {
	//Извлекаем параметры из URL
	username := r.URL.Query().Get("username")
	password := r.URL.Query().Get("password")
	fmt.Println(username, password)
	// Вызов метода сервиса для регистрации пользователя
	token, err := c.servicer.LoginUser(username, password)
	if err != nil {
		// Обработка ошибок сервиса
		c.responder.ErrorInternal(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Expires-After", "2024-04-05T12:00:00Z") // Пример даты и времени в UTC
	w.Header().Set("X-Rate-Limit", "1000")                    // Пример количества вызовов в час
	//json.NewEncoder(w).Encode(token)
	w.Header().Set("Authorization", "Bearer "+token)
	// Отправка успешного ответа клиенту
	c.responder.OutputJSON(w, "Успешный вход")
}

func (c *UserController) LogoutUser(w http.ResponseWriter, r *http.Request) {
	//Извлекаем параметры из URL
	token := r.Header.Get("Authorization")
	token = strings.TrimPrefix(token, "Bearer ")
	//fmt.Println(token)
	token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJ1c2VyIn0.8jVjALlPRYpE03sMD8kuqG9D4RSih5NjiISNZ-wO3oY"
	//tokenString := "your_jwt_token_here"
	claims, err := DecodeJWTToken(token)
	if err != nil {
		fmt.Println("Ошибка декодирования токена:", err)
		return
	}

	sub := claims["sub"].(string)
	fmt.Println("Пользователь:", sub)

	// Вызов метода сервиса для регистрации пользователя
	err = c.servicer.LogoutUser(sub)
	if err != nil {
		// Обработка ошибок сервиса
		c.responder.ErrorInternal(w, err)
		return
	}
	// Отправка успешного ответа клиенту
	c.responder.OutputJSON(w, "Пользователь разлогинен")
}

func DecodeJWTToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	})
	if err != nil {
		return nil, fmt.Errorf("ошибка декодирования токена: %v", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("недействительный токен")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("ошибка получения утверждений из токена")
	}
	// здесь можно Blacklist.Add
	return claims, nil
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

func (c *UserController) FindPetByStatus(w http.ResponseWriter, r *http.Request) {
	status := r.URL.Query().Get("status")
	fmt.Println(status)
	pets, err := c.servicer.FindPetByStatus(status)
	if err != nil {
		// Обработка ошибок сервиса
		c.responder.ErrorInternal(w, err)
		return
	}
	// Отправка успешного ответа клиенту
	c.responder.OutputJSON(w, pets)
}

func (c *UserController) AddPet(w http.ResponseWriter, r *http.Request) {
	var requestBody PetModel.Pet
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		c.responder.ErrorBadRequest(w, err)
		return
	}
	fmt.Println(requestBody)
	err := c.servicer.AddPet(requestBody)
	if err != nil {
		// Обработка ошибок сервиса
		c.responder.ErrorInternal(w, err)
		return
	}
	// Отправка успешного ответа клиенту
	c.responder.OutputJSON(w, "Ok")
}

func (c *UserController) UpdatePet(w http.ResponseWriter, r *http.Request) {
	var requestBody PetModel.Pet
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		c.responder.ErrorBadRequest(w, err)
		return
	}
	fmt.Println(requestBody)
	err := c.servicer.UpdatePet(requestBody)
	if err != nil {
		// Обработка ошибок сервиса
		c.responder.ErrorInternal(w, err)
		return
	}
	// Отправка успешного ответа клиенту
	c.responder.OutputJSON(w, "Ok")
}

func (c *UserController) UpdatePetWithData(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	name := r.Form.Get("name")
	//fmt.Println(name)
	status := r.Form.Get("status")
	//fmt.Println(status)
	path := r.URL.Path
	parts := strings.Split(path, "/")
	petId := parts[len(parts)-1]
	id, err := strconv.Atoi(petId)
	//fmt.Println(id)
	if err != nil {
		fmt.Println(err)
	}
	//fmt.Println(id)
	err = c.servicer.UpdatePetWithData(id, name, status)
	if err != nil {
		// Обработка ошибок сервиса
		c.responder.ErrorInternal(w, err)
		return
	}
	// Отправка успешного ответа клиенту
	c.responder.OutputJSON(w, "Ok")
}

func (c *UserController) DeletePet(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("api_key")
	fmt.Println(name)
	path := r.URL.Path
	parts := strings.Split(path, "/")
	petId := parts[len(parts)-1]
	id, err := strconv.Atoi(petId)
	fmt.Println(id)
	if err != nil {
		fmt.Println(err)
	}
	err = c.servicer.DeletePet(id)
	if err != nil {
		// Обработка ошибок сервиса
		c.responder.ErrorInternal(w, err)
		return
	}
	// Отправка успешного ответа клиенту
	c.responder.OutputJSON(w, "Ok")
}

func (c *UserController) FindPetById(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	parts := strings.Split(path, "/")
	petId := parts[len(parts)-1]
	id, err := strconv.Atoi(petId)
	fmt.Println(id)
	pets, err := c.servicer.FindPetById(id)
	if err != nil {
		// Обработка ошибок сервиса
		c.responder.ErrorInternal(w, err)
		return
	}
	// Отправка успешного ответа клиенту
	c.responder.OutputJSON(w, pets)
}

func (c *UserController) Inventory(w http.ResponseWriter, r *http.Request) {
	//Извлекаем параметры из URL

	// Вызов метода сервиса для регистрации пользователя
	props, err := c.servicer.Inventory()
	if err != nil {
		// Обработка ошибок сервиса
		c.responder.ErrorInternal(w, err)
		return
	}
	// Отправка успешного ответа клиенту
	c.responder.OutputJSON(w, props)
}

func (c *UserController) AddOrder(w http.ResponseWriter, r *http.Request) {
	var requestBody orderModel.Order
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		c.responder.ErrorBadRequest(w, err)
		return
	}
	fmt.Println(requestBody)
	err := c.servicer.AddOrder(requestBody)
	if err != nil {
		// Обработка ошибок сервиса
		c.responder.ErrorInternal(w, err)
		return
	}
	// Отправка успешного ответа клиенту
	c.responder.OutputJSON(w, "Ok")
}

func (c *UserController) FindOrderById(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	parts := strings.Split(path, "/")
	fmt.Println(parts)
	petId := parts[len(parts)-1]
	fmt.Println(petId)
	id, err := strconv.Atoi(petId)
	fmt.Println(id)
	order, err := c.servicer.FindOrderById(id)
	if err != nil {
		// Обработка ошибок сервиса
		c.responder.ErrorInternal(w, err)
		return
	}
	// Отправка успешного ответа клиенту
	c.responder.OutputJSON(w, order)
}

func (c *UserController) DeleteOrder(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	parts := strings.Split(path, "/")
	petId := parts[len(parts)-1]
	id, err := strconv.Atoi(petId)
	fmt.Println(id)
	if err != nil {
		fmt.Println(err)
	}
	err = c.servicer.DeleteOrder(id)
	if err != nil {
		// Обработка ошибок сервиса
		c.responder.ErrorInternal(w, err)
		return
	}
	// Отправка успешного ответа клиенту
	c.responder.OutputJSON(w, "Ok")
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
