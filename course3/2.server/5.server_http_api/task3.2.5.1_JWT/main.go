// Авторизация go-chi с использованием jwt
// В этой задаче необходимо реализовать авторизацию с использованием библиотеки jwtauth из go-chi.
// Потребуется создать две конечных точки: /api/login и /api/register — которые будут отвечать за процессы входа и
// регистрации соответственно.
// В качестве сервиса используй геосервис из предыдущего задания. При обращении к эндпоинтам /api/address/search и
// /api/address/geocode должна происходить проверка токена. Если токен валидный, то запрос должен быть обработан, если нет,
// то должен возвращаться статус код 403.
// При регистрации пользователя  храни его в памяти. При входе пользователя проверяй, что пользователь существует и пароль
// совпадает с тем, что был указан при регистрации. Если пользователь не существует или пароль не совпадает, то возвращай
// статус код 200 и описание ошибки в теле ответа.

package main

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/jwtauth"
	"go-kata/2.server/5.server_http_api/task3.2.5.1_JWT/Handlers"
	"html/template"
	"net/http"
	"time"
)

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

func swaggerUI(w http.ResponseWriter, r *http.Request) {
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

func main() {

	r := chi.NewRouter()
	// Создание экземпляра аутентификации JWT
	Handlers.TokenAuth = jwtauth.New("HS256", []byte("secret"), nil)

	r.Post("/api/register", Handlers.HandlerRegister)
	r.Post("/api/login", Handlers.HandlerLogin)
	r.Get("/swagger/index.html", swaggerUI)
	r.Get("/public/*", func(w http.ResponseWriter, r *http.Request) {
		http.StripPrefix("/public/", http.FileServer(http.Dir("./2.server/5.server_http_api/task3.2.5.1_JWT/public"))).ServeHTTP(w, r)
	})

	r.Group(func(r chi.Router) {
		r.Use(jwtauth.Verifier(Handlers.TokenAuth))
		//r.Use(jwtauth.Authenticator)
		r.Use(Handlers.UnauthorizedToForbidden)

		r.Post("/api/address/search", Handlers.HandlerSearchByQuery)
		r.Post("/api/address/geocode", Handlers.HandlerSearchByGeo)
		http.ListenAndServe(":8080", r)
	})
}

//Критерии приемки
//Запись пароля в память должна быть защищена с помощью bcrypt.
//Должна быть реализована  конечная точка /api/login, которая обрабатывает процесс входа пользователя.
//Должна быть реализована конечная точка /api/register, которая обрабатывает процесс регистрации пользователя.
//Должна быть использована библиотека jwtauth из go-chi для генерации и проверки JWT-токенов.
//В ответе на успешную авторизацию должен возвращаться JWT-токен.
//Проверка токена должна происходить в middleware.
//Все эндпоинты должна быть документированы с помощью swagger.

//Функционал должен быть покрыт тестами на 90%.
//Решение расположи в отдельном проекте

//"200": {
//"description": "Пользователь зарегистрирован."

//"RegisterResponse": {
//"description": "успешная регистрация нового пользователя.",
//"schema": {
//"type": "string",
//"example": "Пользователь зарегистрирован"
