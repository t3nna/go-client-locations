package main

import (
	"context"
	"encoding/json"
	"go-clinet-locations/shared/contracts"
	"go-clinet-locations/shared/messaging"
	"go-clinet-locations/shared/types"
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
	return c.rabbitmq.ConsumeMessages(messaging.SaveUserLocationQueue, func(ctx context.Context, msg amqp091.Delivery) error {
		var userEvent contracts.AmqpMessage

		if err := json.Unmarshal(msg.Body, &userEvent); err != nil {
			return err
		}

		var payload types.UserLocation

		if err := json.Unmarshal(userEvent.Data, &payload); err != nil {
			return err
		}

		log.Printf("user data received: %+v", payload)

		return nil
	})
}
