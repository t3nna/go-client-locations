package events

import (
	"context"
	"go-clinet-locations/shared/messaging"
)

type UserEvenPublisher struct {
	rabbitmq *messaging.RabbitMQ
}

func NewUserEventPublisher(rabbitmq *messaging.RabbitMQ) *UserEvenPublisher {
	return &UserEvenPublisher{
		rabbitmq: rabbitmq,
	}
}

func (p *UserEvenPublisher) PublishUserCreated(ctx context.Context) error {
	return p.rabbitmq.PublishMessage(ctx, "hello", "hello world")
}
