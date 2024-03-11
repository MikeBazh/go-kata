//Реализуй API-сервер на языке программирования Golang с применением graceful shutdown.
//	Сервер должен обрабатывать HTTP-запросы и корректно завершать свою работу при получении сигнала остановки.

// Пример кода для реализации API-сервера с graceful shutdown

package main

import (
	"context"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/jwtauth"
	"github.com/ptflp/godecoder"
	"go-kata/2.server/5.server_http_api/layer_controller/Handlers"
	"go-kata/2.server/5.server_http_api/layer_controller/controller"
	"go-kata/2.server/5.server_http_api/layer_controller/responder"
	"go.uber.org/zap"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
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

	// Создание экземпляра сервиса
	servicer := Handlers.NewService()
	var decoder godecoder.Decoder
	var logger *zap.Logger
	responder := responder.NewResponder(decoder, logger)
	// Создание экземпляра контроллера
	UserController := controller.NewUserController(responder, servicer)

	r := chi.NewRouter()
	// Создание экземпляра аутентификации JWT
	Handlers.TokenAuth = jwtauth.New("HS256", []byte("secret"), nil)

	r.Post("/api/register", UserController.Register)
	//r.Post("/api/login", Handlers.HandlerLogin)
	r.Get("/swagger/index.html", swaggerUI)
	r.Get("/public/*", func(w http.ResponseWriter, r *http.Request) {
		http.StripPrefix("/public/", http.FileServer(http.Dir("./2.server/5.server_http_api/layer_controller/public"))).ServeHTTP(w, r)
	})

	r.Group(func(r chi.Router) {
		r.Use(jwtauth.Verifier(Handlers.TokenAuth))
		//r.Use(jwtauth.Authenticator)
		//r.Use(Handlers.UnauthorizedToForbidden)

		//r.Post("/api/address/search", Handlers.HandlerSearchByQuery)
		//r.Post("/api/address/geocode", Handlers.HandlerSearchByGeo)
	})

	port := ":8080"
	server := &http.Server{
		Addr:    port,
		Handler: r}

	// Создание канала для получения сигналов остановки
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT)

	// Запуск сервера в отдельной горутине
	go func() {
		log.Printf("server started on port %s ", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Обработка сигнала SIGINT
	sig := <-sigChan
	fmt.Printf("Received signal: %s\n", sig)

	// Создание контекста с таймаутом для graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// Остановка сервера с использованием graceful shutdown
	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	log.Println("Server stopped gracefully")
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
