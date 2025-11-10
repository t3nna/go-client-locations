package main

import (
	"context"
	"go-clinet-locations/shared/messaging"
	"log"

	"github.com/rabbitmq/amqp091-go"
)

type userConsumer struct {
	rabbitmq *messaging.RabbitMQ
}

func NewUserConsumer(rabbitmq *messaging.RabbitMQ) *userConsumer {
	return &userConsumer{
		rabbitmq: rabbitmq,
	}
}

func (c *userConsumer) Listen() error {
	return c.rabbitmq.ConsumeMessages("hello", func(ctx context.Context, msg amqp091.Delivery) error {
		log.Printf("location management received message: %v", msg)

		return nil
	})
}
