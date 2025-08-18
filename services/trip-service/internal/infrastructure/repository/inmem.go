package repository

import (
	"context"
	"fmt"
	"go-clinet-locations/services/trip-service/internal/domain"
	"go-clinet-locations/shared/types"
	"sync"
)

type inmemRepository struct {
	users map[string]*domain.UserModel
	mu    sync.RWMutex
}

func newInmemRepo() *inmemRepository {
	return &inmemRepository{
		users: map[string]*domain.UserModel{
			"userId1234": {
				UserName: "Ivan",
				UserId:   "userId1234",
				Coordinates: &types.Coordinate{
					Latitude:  51.10861618949566,
					Longitude: 17.03187985482019,
				},
			},
		},
	}
}

func (r *inmemRepository) CreateUser(ctx context.Context, user *domain.UserModel) (*domain.UserModel, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.users[user.UserId] = user
	return user, nil
}
func (r *inmemRepository) UpdateUser(ctx context.Context, userName string, coordinates *types.Coordinate) (*domain.UserModel, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	for key, user := range r.users {
		if user.UserName == userName {

			updatedUser := &domain.UserModel{
				UserId:      user.UserId,
				UserName:    user.UserName,
				Coordinates: coordinates,
			}

			r.users[key] = updatedUser

			return updatedUser, nil
		}
	}

	return nil, fmt.Errorf("failde to find user in DB")
}
