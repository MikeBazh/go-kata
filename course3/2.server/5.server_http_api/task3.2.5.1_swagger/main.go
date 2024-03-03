// Создание документации к геосервису с помощью go-swagger
// Задание
// Напиши документацию к геосервису из предыдущего задания с использованием go-swagger.
// Документация должна описывать API эндпоинты /api/address/search и /api/address/geocode.
// Сервис должен принимать POST запросы с параметром query и возвращать информацию о городе,
//в котором находится данный адрес.
// Ознакомься с материалами по документации swagger и go-swagger.

package main

import (
	"github.com/go-chi/chi"
	"go-kata/2.server/5.server_http_api/task3.2.5.1_swagger/Handlers"
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

	r.Get("/swagger/index.html", swaggerUI)

	r.Get("/public/*", func(w http.ResponseWriter, r *http.Request) {
		http.StripPrefix("/public/", http.FileServer(http.Dir("./2.server/5.server_http_api/task3.2.5.1_swagger/public"))).ServeHTTP(w, r)
	})

	r.Post("/api/address/search", Handlers.HandlerSearchByQuery)
	r.Post("/api/address/geocode", Handlers.HandlerSearchByGeo)
	http.ListenAndServe(":8080", r)
}

//Критерии приемки
//Документация должна быть доступна по адресу http://localhost:8080/swagger/index.html.
//Документация должна описывать API эндпоинты /api/address/search и /api/address/geocode.
//Функционал должен быть покрыт тестами 90%.
//Решение расположи в отдельном проекте
