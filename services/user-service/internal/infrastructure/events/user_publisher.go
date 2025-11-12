package events

import (
	"context"
	"encoding/json"
	"go-clinet-locations/shared/contracts"
	"go-clinet-locations/shared/messaging"
	"go-clinet-locations/shared/types"
)

type UserEvenPublisher struct {
	rabbitmq *messaging.RabbitMQ
}

func NewUserEventPublisher(rabbitmq *messaging.RabbitMQ) *UserEvenPublisher {
	return &UserEvenPublisher{
		rabbitmq: rabbitmq,
	}
}

func (p *UserEvenPublisher) PublishUserCreated(ctx context.Context, userLocation *types.UserLocation) error {
	userEventJSON, err := json.Marshal(userLocation)
	if err != nil {
		return err
	}

	return p.rabbitmq.PublishMessage(ctx, messaging.RegisterLocationEventBind, contracts.AmqpMessage{
		OwnerID: userLocation.UserId,
		Data:    userEventJSON,
	})
}
