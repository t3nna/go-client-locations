package repository

import (
	"context"
	"fmt"
	"go-clinet-locations/services/user-service/internal/domain"
	"go-clinet-locations/shared/db"
	"go-clinet-locations/shared/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"sync"
)

type mongoRepository struct {
	db *mongo.Database
	mu sync.RWMutex
}

func NewMongoRepository(db *mongo.Database) *mongoRepository {
	return &mongoRepository{db: db}
}

// TODO: implement UserRepository interface

func (r *mongoRepository) CreateUser(ctx context.Context, user *domain.UserModel) (*domain.UserModel, error) {
	result, err := r.db.Collection(db.UserCollection).InsertOne(ctx, user)

	if err != nil {
		return nil, err
	}

	user.ID = result.InsertedID.(primitive.ObjectID)

	return user, nil
}
func (r *mongoRepository) UpdateUser(ctx context.Context, userName string, coordinates *types.Coordinate) (*domain.UserModel, error) {
	collection := r.db.Collection(db.UserCollection)
	filter := bson.M{"userName": userName}
	update := bson.M{"$set": bson.M{"coordinates": bson.M{"latitude": coordinates.Latitude, "longitude": coordinates.Longitude}}}

	result := collection.FindOneAndUpdate(ctx, filter, update, options.FindOneAndUpdate().SetReturnDocument(options.After))
	if result.Err() != nil {
		return nil, fmt.Errorf("failed to update user: %v", result.Err())
	}

	var updatedUser domain.UserModel
	if err := result.Decode(&updatedUser); err != nil {
		return nil, fmt.Errorf("failed to decode updated user: %v", err)
	}

	return &updatedUser, nil
}

func (r *mongoRepository) GetUsers(ctx context.Context) ([]*domain.UserModel, error) {
	collection := r.db.Collection(db.UserCollection)

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("failed to get users: %v", err)
	}
	defer cursor.Close(ctx)

	var users []*domain.UserModel
	for cursor.Next(ctx) {
		var user domain.UserModel
		if err := cursor.Decode(&user); err != nil {
			return nil, fmt.Errorf("failed to decode user: %v", err)
		}
		users = append(users, &domain.UserModel{
			ID:       user.ID,
			UserName: user.UserName,
			Coordinates: &types.Coordinate{
				Latitude:  user.Coordinates.Latitude,
				Longitude: user.Coordinates.Longitude,
			},
		})
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor error: %v", err)
	}

	return users, nil

}
