package main

import (
	"context"
	"fmt"
	"github.com/segmentio/kafka-go"
	"log"
	"net"
	"net/rpc"
	"os"
)

type KafkaPublisher struct {
}

func main() {
	KafkaPublisher := new(KafkaPublisher)
	rpc.Register(KafkaPublisher)

	// Запуск RPC сервера.
	rpc.HandleHTTP()

	// Запуск сервера на порту 8070.
	listener, err := net.Listen("tcp", ":8070")
	if err != nil {
		log.Fatal("Error starting RPC server:", err)
	}
	log.Println("RPC server started on port 8070")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal("Ошибка при принятии соединения:", err)
		}
		go rpc.ServeConn(conn)
	}
	//SendMessage("", []byte("Hello"))
}

type SendMessageArgs struct {
	Message []byte
}

func (g *KafkaPublisher) SendMessage(args *SendMessageArgs, reply *bool) error {
	topic := "my-topic"
	kafkaBroker := os.Getenv("KAFKA_BROKER_ADDRESS")
	// Создание конфигурации издателя
	writer := kafka.Writer{
		Addr:     kafka.TCP(kafkaBroker),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}

	err := writer.WriteMessages(context.Background(),
		kafka.Message{Value: args.Message},
	)
	if err != nil {
		fmt.Println(err)
		*reply = false
		return err
	}
	*reply = true
	return nil
}
