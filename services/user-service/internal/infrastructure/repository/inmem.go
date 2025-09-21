package repository

import (
	"context"
	"fmt"
	"go-clinet-locations/services/user-service/internal/domain"
	"go-clinet-locations/shared/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"sync"
)

type inmemRepository struct {
	users map[string]*domain.UserModel
	mu    sync.RWMutex
}

func NewInmemRepository() *inmemRepository {
	return &inmemRepository{
		users: map[string]*domain.UserModel{
			// Magnolia
			"user1": {
				UserName: "Ivan",
				ID:       primitive.NewObjectID(),
				Coordinates: &types.Coordinate{
					Latitude:  51.11822470712269,
					Longitude: 16.990711729269563,
				},
			},
			// Warsaw
			"user2": {
				UserName: "Igor",
				ID:       primitive.NewObjectID(),
				Coordinates: &types.Coordinate{
					Latitude:  52.23553956649786,
					Longitude: 20.984595191389918,
				},
			},

			// Kyiv
			"user3": {
				UserName: "Den",
				ID:       primitive.NewObjectID(),
				Coordinates: &types.Coordinate{
					Latitude:  50.53401932980686,
					Longitude: 31.178889172903055,
				},
			},
			//Biskupin
			"user4": {
				UserName: "Kate",
				ID:       primitive.NewObjectID(),
				Coordinates: &types.Coordinate{
					Latitude:  51.10181370006046,
					Longitude: 17.10312341673202,
				},
			},
			// T6
			"user5": {
				UserName: "Barbara",
				ID:       primitive.NewObjectID(),
				Coordinates: &types.Coordinate{
					Latitude:  51.11956092410769,
					Longitude: 17.05696305051491,
				},
			},
		},
	}
}

func (r *inmemRepository) CreateUser(ctx context.Context, user *domain.UserModel) (*domain.UserModel, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.users[user.ID.String()] = user
	return user, nil
}
func (r *inmemRepository) UpdateUser(ctx context.Context, userName string, coordinates *types.Coordinate) (*domain.UserModel, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	for key, user := range r.users {
		if user.UserName == userName {

			updatedUser := &domain.UserModel{
				ID:          user.ID,
				UserName:    user.UserName,
				Coordinates: coordinates,
			}

			r.users[key] = updatedUser

			return updatedUser, nil
		}
	}

	return nil, fmt.Errorf("failed to find user in DB")
}

func (r *inmemRepository) GetUsers(ctx context.Context) ([]*domain.UserModel, error) {
	// TODO: check if I need to use mutexes on getting

	result := make([]*domain.UserModel, 0, 10)

	for _, value := range r.users {
		result = append(result, value)
	}
	return result, nil
}
