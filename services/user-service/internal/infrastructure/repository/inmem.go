package repository

import (
	"context"
	"fmt"
	"go-clinet-locations/services/user-service/internal/domain"
	"go-clinet-locations/shared/types"
	"sync"
)

type inmemRepository struct {
	users map[string]*domain.UserModel
	mu    sync.RWMutex
}

func NewInmemRepository() *inmemRepository {
	return &inmemRepository{
		users: map[string]*domain.UserModel{
			"user1": {
				UserName: "Ivan",
				UserId:   "user1",
				Coordinates: &types.Coordinate{
					Latitude:  51.10861618949566,
					Longitude: 17.03187985482019,
				},
			},
			"user2": {
				UserName: "Ivan",
				UserId:   "user2",
				Coordinates: &types.Coordinate{
					Latitude:  51.5,
					Longitude: 17.7,
				},
			},
			"user3": {
				UserName: "Ivan",
				UserId:   "user3",
				Coordinates: &types.Coordinate{
					Latitude:  51.9,
					Longitude: 17.0,
				},
			},
			"user4": {
				UserName: "Ivan",
				UserId:   "user4",
				Coordinates: &types.Coordinate{
					Latitude:  51.32,
					Longitude: 17.88,
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

func (r *inmemRepository) GetUsers(ctx context.Context) ([]*domain.UserModel, error) {
	// TODO: check if I need to use mutexes on getting

	result := make([]*domain.UserModel, len(r.users))

	for _, value := range r.users {
		result = append(result, value)
	}
	return result, nil
}
