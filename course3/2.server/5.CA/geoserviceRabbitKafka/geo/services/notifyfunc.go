package services

import (
	"encoding/json"
	"fmt"
	"github.com/streadway/amqp"
	"gitlab.com/ptflp/gopubsub/queue"
	"gitlab.com/ptflp/gopubsub/rabbitmq"
	"net/rpc"
)

var (
	rmq queue.MessageQueuer
)

type MessangerKafka struct {
}

type MessangerRabbitMQ struct {
}

type MessangerService interface {
	SendRateLimitExceededMessage(email string) error
}

type NotificationMessage struct {
	Email string
	Mes   string
}

// sendRateLimitExceededMessage отправляет сообщение в RabbitMQ о превышении лимита
func (s *MessangerRabbitMQ) SendRateLimitExceededMessage(email string) error {
	// Установка соединения с RabbitMQ
	conn, err := amqp.Dial("amqp://guest:guest@rabbitmq:5672/")
	if err != nil {
		fmt.Printf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	rmq, err = rabbitmq.NewRabbitMQ(conn)
	if err != nil {
		fmt.Printf("Failed to create RabbitMQ instance: %v", err)
	}
	defer rmq.Close()

	// Создаем обменник для уведомлений о превышении лимита
	if err := rabbitmq.CreateExchange(conn, "rate_limit_exchange", "fanout"); err != nil {
		fmt.Printf("Failed to create exchange: %v", err)
	}

	message := NotificationMessage{
		Email: email,
		Mes:   "Rate limit exceeded",
	}
	body, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	err = rmq.Publish("rate_limit_exchange", body)
	if err != nil {
		return fmt.Errorf("failed to publish message to RabbitMQ: %w", err)
	}
	fmt.Println("Rate limit exceeded message sent to RabbitMQ, user - ", email)
	return nil
}

type SendMessageArgs struct {
	Message []byte
}

// sendRateLimitExceededMessage отправляет сообщение в RabbitMQ о превышении лимита
func (s *MessangerKafka) SendRateLimitExceededMessage(email string) error {
	client, err := rpc.Dial("tcp", "publisherkafka:8070")
	//var reply []*dto.Address
	if err != nil {
		fmt.Println("Ошибка при подключении к серверу:", err)
		return err
	}
	args := &SendMessageArgs{Message: []byte(email)}
	var reply *bool
	err = client.Call("KafkaPublisher.SendMessage", args, reply)
	if err != nil {
		fmt.Println("Ошибка при вызове удаленного метода:", err)
		return err
	}
	return nil
}
