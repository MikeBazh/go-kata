package main

import (
	"encoding/json"
	"fmt"
	"github.com/streadway/amqp"
	"gitlab.com/ptflp/gopubsub/rabbitmq"
	"log"
)

type NotificationMessage struct {
	Email string `json:"email"`
	Mes   string `json:"mes"`
}

func main() {
	// Подключение к RabbitMQ
	conn, err := amqp.Dial("amqp://guest:guest@rabbitmq:5672/")
	if err != nil {
		fmt.Println("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	rmq, err := rabbitmq.NewRabbitMQ(conn)
	if err != nil {
		fmt.Println("Failed to create RabbitMQ instance: %v", err)
	}
	defer rmq.Close()

	// Создаем обменник для уведомлений о превышении лимита
	if err := rabbitmq.CreateExchange(conn, "rate_limit_exchange", "fanout"); err != nil {
		fmt.Printf("Failed to create exchange: %v", err)
	}

	// Подписка на сообщения из обменника
	messages, err := rmq.Subscribe("rate_limit_exchange")
	if err != nil {
		fmt.Println("Failed to subscribe to exchange: %v", err)
	}

	// Канал для завершения работы
	forever := make(chan bool)

	go func() {
		for msg := range messages {
			handleMessage(msg.Data)
			err := rmq.Ack(&msg)
			if err != nil {
				fmt.Println("Failed to acknowledge message: %v", err)
			}
		}
	}()

	log.Println("Waiting for messages. To exit press CTRL+C")
	<-forever
}

func handleMessage(body []byte) {
	var notification NotificationMessage
	err := json.Unmarshal(body, &notification)
	if err != nil {
		fmt.Println("Error decoding JSON: %v", err)
		return
	}

	// отправка смс уведомления (просто вывод в терминал)
	fmt.Printf("Отправка email для пользователя: %s, Сообщение: %s\n", notification.Email, notification.Mes)
}
