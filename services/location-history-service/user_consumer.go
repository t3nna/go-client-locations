package main

import (
	"context"
	"encoding/json"
	"go-clinet-locations/shared/contracts"
	"go-clinet-locations/shared/messaging"
	"go-clinet-locations/shared/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"time"

	"github.com/rabbitmq/amqp091-go"
)

type userConsumer struct {
	rabbitmq *messaging.RabbitMQ
	service  LocationsService
}

func NewUserConsumer(rabbitmq *messaging.RabbitMQ, service LocationsService) *userConsumer {
	return &userConsumer{
		rabbitmq: rabbitmq,
		service:  service,
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

		now := time.Now()

		locationRecords, err := c.service.RegisterLocation(ctx, payload.UserId, payload.Coordinate, now)

		if err != nil {
			return status.Errorf(codes.Internal, "failed to Register Location %v", err)
		}
		log.Printf("%+v", locationRecords)

		return nil
	})
}
