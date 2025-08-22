package main

import (
	"go-clinet-locations/shared/types"
	"sync"
	"time"
)

type LocationRecord struct {
	Coordinate *types.Coordinate
	Timestamp  time.Time
}
type Service struct {
	history map[string][]*LocationRecord
	mu      sync.RWMutex
}

func NewService() *Service {
	now := time.Now()
	return &Service{
		history: map[string][]*LocationRecord{
			"user1": []*LocationRecord{
				{
					Coordinate: &types.Coordinate{
						Latitude:  51.11822470712269,
						Longitude: 16.990711729269563,
					},
					Timestamp: now.Add(-2 * time.Hour * 24),
				},
				{
					Coordinate: &types.Coordinate{
						Latitude:  51.11822470712269,
						Longitude: 16.990711729269563,
					},
					Timestamp: now.Add(-1 * time.Hour),
				},
				{
					Coordinate: &types.Coordinate{
						Latitude:  51.11822470712269,
						Longitude: 16.990711729269563,
					},
					Timestamp: now.Add(-30 * time.Minute),
				},
			},
		},
	}
}

func (s *Service) RegisterLocation(userId string, coords *types.Coordinate, timestamp time.Time) ([]*LocationRecord, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	user, ok := s.history[userId]
	if !ok {
		s.history[userId] = []*LocationRecord{
			{
				Coordinate: coords,
				Timestamp:  timestamp,
			},
		}
		return s.history[userId], nil
	}

	user = append(user, &LocationRecord{
		Coordinate: coords,
		Timestamp:  timestamp,
	},
	)

	return user, nil
}
