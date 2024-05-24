package controller

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	pb "go-kata/2.server/5.CA/GeoAuthUserProxy/geo/dto/auth/grpc"
	"go-kata/2.server/5.CA/GeoAuthUserProxy/geo/services"
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

	//ErrorUnauthorized(w http.ResponseWriter, err error)
	ErrorBadRequest(w http.ResponseWriter, err error)
	ErrorForbidden(w http.ResponseWriter, err error)
	ErrorInternal(w http.ResponseWriter, err error)
}

type UserController struct {
	responder     Responder
	servicer      services.Servicer
	proxyservicer services.ServiceProxy
}

func NewUserController(responder Responder, servicer services.Servicer, proxyservicer services.ServiceProxy) *UserController {
	return &UserController{
		responder:     responder,
		servicer:      servicer,
		proxyservicer: proxyservicer,
	}
}

func (c *UserController) SearchByQuery(w http.ResponseWriter, r *http.Request) {
	if !Validation(r) {
		c.responder.ErrorForbidden(w, errors.New("validation needed"))
		return
	}

	// Предварительная обработка запроса
	var requestBody SearchRequest
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		c.responder.ErrorBadRequest(w, err)
		return
	}

	// Вызов метода сервиса для регистрации пользователя

	respond, err := c.proxyservicer.SearchByQuery(requestBody.Query)
	//respond, err := c.servicer.SearchByQuery(requestBody.Query)
	if err != nil {
		// Обработка ошибок сервиса
		c.responder.ErrorInternal(w, err)
		return
	}
	// Отправка успешного ответа клиенту
	c.responder.OutputJSON(w, respond)
}

func (c *UserController) SearchByGeo(w http.ResponseWriter, r *http.Request) {
	if !Validation(r) {
		c.responder.ErrorForbidden(w, errors.New("validation needed"))
		return
	}
	// Предварительная обработка запроса
	var requestBody GeocodeRequest
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		c.responder.ErrorBadRequest(w, err)
		return
	}

	// Вызов метода сервиса для регистрации пользователя
	respond, err := c.proxyservicer.SearchByGeo(requestBody.Lat, requestBody.Lng)
	if err != nil {
		// Обработка ошибок сервиса
		c.responder.ErrorInternal(w, err)
		return
	}
	// Отправка успешного ответа клиенту
	c.responder.OutputJSON(w, respond)
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