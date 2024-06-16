package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

func main() {
	// Настройка параметров подключения к Kafka
	topic := "my-topic"
	partition := 0
	//offset := kafka.LastOffset

	// Создание конфигурации потребителя
	config := kafka.ReaderConfig{
		Brokers:   []string{"kafka:9093"},
		Topic:     topic,
		Partition: partition,
		MinBytes:  10e3,
		MaxBytes:  10e6,
		MaxWait:   1 * time.Second,
	}

	// Создание потребителя Kafka
	reader := kafka.NewReader(config)
	defer reader.Close()

	// Чтение сообщений из топика
	for {
		msg, err := reader.ReadMessage(context.Background())
		if err != nil {
			log.Fatal(err)
		}
		// отправка смс уведомления (просто вывод в терминал)
		fmt.Printf("Отправка email для пользователя: %s, Сообщение: %s\n", msg.Value, "Rate limit exceeded")
	}
}
