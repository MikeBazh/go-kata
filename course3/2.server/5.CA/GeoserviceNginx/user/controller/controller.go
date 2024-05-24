package controller

import (
	"context"
	"errors"
	"fmt"
	_ "github.com/dgrijalva/jwt-go"
	_ "go-kata/2.server/5.CA/GeoserviceNginx/user/dto"
	pb "go-kata/2.server/5.CA/GeoserviceNginx/user/dto/auth"
	"go-kata/2.server/5.CA/GeoserviceNginx/user/services"
	"google.golang.org/grpc"
	"html/template"
	"log"
	"net/http"
	"strings"
	"time"
)

const (
	authAddress = "grpcAuth:50051"
)

type Responder interface {
	OutputJSON(w http.ResponseWriter, responseData interface{})

	ErrorBadRequest(w http.ResponseWriter, err error)
	ErrorForbidden(w http.ResponseWriter, err error)
	ErrorInternal(w http.ResponseWriter, err error)
}

type UserController struct {
	responder Responder
	serviser  services.Service
}

func NewUserController(responder Responder, serviser services.Service) *UserController {
	return &UserController{
		responder: responder,
		serviser:  serviser,
	}
}

func (c *UserController) ListUsers(w http.ResponseWriter, r *http.Request) {
	if !Validation(r) {
		c.responder.ErrorForbidden(w, errors.New("validation needed"))
		return
	}

	// Вызов метода сервиса для регистрации пользователя
	users, err := c.serviser.ListUsers()
	if err != nil {
		// Обработка ошибок сервиса
		c.responder.ErrorInternal(w, err)
		return
	}
	// Отправка успешного ответа клиенту
	c.responder.OutputJSON(w, users)
}

func (c *UserController) Profile(w http.ResponseWriter, r *http.Request) {
	if !Validation(r) {
		c.responder.ErrorForbidden(w, errors.New("validation needed"))
		return
	}
	//Извлекаем параметры из URL
	email := r.URL.Query().Get("UserEmail")
	fmt.Println(email)
	// Вызов метода сервиса для регистрации пользователя
	user, err := c.serviser.Profile(email)
	if err != nil {
		// Обработка ошибок сервиса
		c.responder.ErrorInternal(w, err)
		return
	}
	// Отправка успешного ответа клиенту
	c.responder.OutputJSON(w, user)
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

func Validation(r *http.Request) bool {
	token, err := getTokenFromHeader(r)
	if err != nil {
		fmt.Println("ошибка чтения токена: ", err)
		return false
	}
	fmt.Println("geo main: токен получен")
	conn, err := grpc.Dial(authAddress, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	// Создаем клиент
	tokenClient := pb.NewTokenValidationServiceClient(conn)
	fmt.Println("запрос в сервис auth...")
	// Вызовите gRPC-сервис для проверки токена
	var resp *pb.TokenResponse
	resp, err = tokenClient.ValidateToken(context.Background(), &pb.TokenRequest{Token: token})
	if err != nil {
		// Обработайте ошибку
		fmt.Println("ошибка валидации: ", err)
		return false
	}
	fmt.Println("geo: resp:", resp)
	if !resp.Valid {
		return false
	}
	return true
}

func getTokenFromHeader(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", fmt.Errorf("Authorization header is missing")
	}

	authParts := strings.Split(authHeader, " ")
	if len(authParts) != 2 || strings.ToLower(authParts[0]) != "bearer" {
		return "", fmt.Errorf("Invalid Authorization header format")
	}
	return authParts[1], nil
}
