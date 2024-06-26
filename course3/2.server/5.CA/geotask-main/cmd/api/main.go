package main

import (
	"fmt"
	"github.com/joho/godotenv"
	run "gitlab.com/ptflp/geotask/run"

	"log"
	"os"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		return
	}
	// инициализация приложения
	app := run.NewApp()
	// запуск приложения
	err = app.Run()
	// в случае ошибки выводим ее в лог и завершаем работу с кодом 2
	if err != nil {
		log.Println(fmt.Sprintf("error: %s", err))
		os.Exit(2)
	}
}
