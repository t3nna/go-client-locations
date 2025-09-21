package main

import (
	"context"
	"go-clinet-locations/shared/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type LocationRecord struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	Coordinate *types.Coordinate  `bson:"coordinate"`
	Timestamp  time.Time          `bson:"timestamp"`
}

type LocationsService interface {
	RegisterLocation(ctx context.Context, userId string, coords *types.Coordinate, timestamp time.Time) ([]*LocationRecord, error)
	CalculateDistance(ctx context.Context, userId string, startDate time.Time, endDate time.Time) (*DistanceRecord, error)
}
