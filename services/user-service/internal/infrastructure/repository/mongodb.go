package repository

import (
	"context"
	"go-clinet-locations/services/user-service/internal/domain"
	"go-clinet-locations/shared/types"
	"go.mongodb.org/mongo-driver/mongo"
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
	return user, nil
}
func (r *mongoRepository) UpdateUser(ctx context.Context, userName string, coordinates *types.Coordinate) (*domain.UserModel, error) {
	return nil, nil
}

func (r *mongoRepository) GetUsers(ctx context.Context) ([]*domain.UserModel, error) {
	return nil, nil
}
