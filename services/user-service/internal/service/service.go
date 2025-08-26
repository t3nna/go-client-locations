package service

import (
	"context"
	"github.com/google/uuid"
	"go-clinet-locations/services/user-service/internal/domain"
	"go-clinet-locations/shared/types"
	"go-clinet-locations/shared/util"
	"log"
)

type service struct {
	repo domain.UserRepository
}

func NewService(repo domain.UserRepository) *service {
	return &service{
		repo: repo,
	}
}

func (s *service) CreateUser(ctx context.Context, user *domain.UserModel) (*domain.UserModel, error) {
	newUUID, err := uuid.NewRandom()
	if err != nil {
		log.Fatalf("failed to generate UUID: %v", err)
	}

	newUser := &domain.UserModel{
		UserId:      newUUID.String(),
		UserName:    user.UserName,
		Coordinates: user.Coordinates,
	}
	return s.repo.CreateUser(ctx, newUser)
}
func (s *service) UpdateUser(ctx context.Context, userName string, coordinates *types.Coordinate) (*domain.UserModel, error) {

	return s.repo.UpdateUser(ctx, userName, coordinates)
}

func (s *service) SearchUsers(ctx context.Context, location *types.Coordinate, radius float64) ([]*domain.UserModel, error) {
	users, err := s.repo.GetUsers(ctx)
	if err != nil {
		log.Fatalf("faled to get users: %v", err)
		return nil, err
	}

	var filteredUsers []*domain.UserModel
	for _, user := range users {
		distance := util.CalculateDistance(location, user.Coordinates)
		log.Println(distance, user.UserName)
		if distance <= radius {
			filteredUsers = append(filteredUsers, user)
		}
	}

	return filteredUsers, nil
}
