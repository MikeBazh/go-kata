package controller

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	pb "go-kata/2.server/5.CA/geoserviceRabbitKafka/geo/dto/auth/grpc"
	"go-kata/2.server/5.CA/geoserviceRabbitKafka/geo/responder"
	"go-kata/2.server/5.CA/geoserviceRabbitKafka/geo/services"
	"google.golang.org/grpc"
	"html/template"
	"net/http"
	"strings"
	"time"
)

const (
	authAddress = "grpcAuth:50051"
)

type UserController struct {
	responder     responder.Responder
	servicer      services.Servicer
	proxyservicer services.ServiceProxy
}

func NewUserController(responder responder.Responder, servicer services.Servicer, proxyservicer services.ServiceProxy) *UserController {
	return &UserController{
		responder:     responder,
		servicer:      servicer,
		proxyservicer: proxyservicer,
	}
}

func (c *UserController) SearchByQuery(w http.ResponseWriter, r *http.Request) {
	valid, email := Validation(r)
	if !valid {
		c.responder.ErrorForbidden(w, errors.New("validation needed"))
		return
	}

	// Предварительная обработка запроса
	var requestBody SearchRequest
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		c.responder.ErrorBadRequest(w, err)
		return
	}

	//respond, err := c.proxyservicer.SearchByQuery(requestBody.Query)
	respond, err, limited := c.servicer.SearchByQuery(requestBody.Query)
	if limited {
		err = c.servicer.SendRateLimitExceededMessage(email)
		if err != nil {
			fmt.Println("Error sending rate limit message")
		}
		// Обработка ошибок сервиса
		c.responder.ErrorTooManyRequests(w, err)
		return
	}
	if err != nil {
		// Обработка ошибок сервиса
		c.responder.ErrorInternal(w, err)
		return
	}
	// Отправка успешного ответа клиенту
	c.responder.OutputJSON(w, respond)
}

func (c *UserController) SearchByGeo(w http.ResponseWriter, r *http.Request) {
	valid, email := Validation(r)
	if !valid {
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
	//respond, err := c.proxyservicer.SearchByGeo(requestBody.Lat, requestBody.Lng)
	respond, err, limited := c.servicer.SearchByGeo(requestBody.Lat, requestBody.Lng)
	if limited {
		err = c.servicer.SendRateLimitExceededMessage(email)
		if err != nil {
			fmt.Println("Error sending rate limit message")
		}
		// Обработка ошибок сервиса
		c.responder.ErrorTooManyRequests(w, err)
		return
	}

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

func Validation(r *http.Request) (bool, string) {
	token, err := getTokenFromHeader(r)
	if err != nil {
		fmt.Println("ошибка чтения токена: ", err)
		return false, ""
	}
	email, err := extractLoginFromToken(token)
	if err != nil {
		fmt.Printf("Failed to extraxt email: %v", err)
		return false, ""
	}

	fmt.Println("geo main: токен получен, email:", email)
	conn, err := grpc.Dial(authAddress, grpc.WithInsecure())
	if err != nil {
		fmt.Printf("Failed to connect: %v", err)
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
		return false, email
	}
	//fmt.Println("geo: resp:", resp)
	if !resp.Valid {
		return false, email
	}
	return true, email
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

func extractLoginFromToken(tokenString string) (string, error) {
	// Удаляем префикс "Bearer " из токена
	//tokenString = tokenString[len("Bearer "):]
	//fmt.Println(tokenString)
	// Декодируем токен
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil // Замените "secret" на ваш секретный ключ для подписи токена
	})

	if err != nil {
		return "", fmt.Errorf("ошибка декодирования JWT токена: %v", err)
	}

	// Извлекаем claims
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if sub, ok := claims["sub"].(string); ok {
			return sub, nil
		}
		return "", fmt.Errorf("токен не содержит поле 'sub'")
	}

	return "", fmt.Errorf("недействительный токен")
}
