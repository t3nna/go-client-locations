package service

import (
	"context"
	"github.com/google/uuid"
	"go-clinet-locations/services/user-service/internal/domain"
	"go-clinet-locations/shared/types"
	"log"
	"math"
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
		distance := calculateDistance(location, user.Coordinates)
		log.Println(distance, user.UserName)
		if distance <= radius {
			filteredUsers = append(filteredUsers, user)
		}
	}

	return filteredUsers, nil
}

// calculateDistance computes the distance between two coordinates using the Haversine formula.
func calculateDistance(coord1, coord2 *types.Coordinate) float64 {
	const earthRadius = 6371 // Earth's radius in kilometers

	lat1, lon1 := degreesToRadians(coord1.Latitude), degreesToRadians(coord1.Longitude)
	lat2, lon2 := degreesToRadians(coord2.Latitude), degreesToRadians(coord2.Longitude)

	deltaLat := lat2 - lat1
	deltaLon := lon2 - lon1

	a := math.Sin(deltaLat/2)*math.Sin(deltaLat/2) +
		math.Cos(lat1)*math.Cos(lat2)*math.Sin(deltaLon/2)*math.Sin(deltaLon/2)

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return earthRadius * c
}

// degreesToRadians converts degrees to radians.
func degreesToRadians(degrees float64) float64 {
	return degrees * math.Pi / 180
}
