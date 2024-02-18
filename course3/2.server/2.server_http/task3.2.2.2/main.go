// Контейнеризируй сервер, используя пакет net/http и Docker.
// Пример кода для сервера на пакете net/http
package main

import (
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error load .env file")
	}
	port := os.Getenv("SERVER_PORT")
	log.Printf("Сервер запущен на порту %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

//Внимание! Учитывай GOPATH при компиляции го-программы внутри контейнера.
//Критерии приемки
//
//Сервер должен быть контейнеризирован.
//Сервер должен быть доступен на порту 8080.
//Сервер должен настраивать порт через файл .env.
//Должен быть использован пакет dotenv.
//Образ должен быть минимальным с использованием, alpine или scratch.
//Проект должен стартовать командой docker-compose up. docker-compose.yml должен быть в корне проекта.
//Покрытие тестами не требуется.
