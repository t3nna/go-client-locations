package main

import (
	"context"
	"fmt"
	pb "go-clinet-locations/shared/proto/location"
	"go-clinet-locations/shared/types"
	"go-clinet-locations/shared/util"
	"log"
	"sync"
	"time"
)

type Service struct {
	history map[string][]*LocationRecord
	mu      sync.RWMutex
}

type DistanceRecord struct {
	distance float64
	history  []*LocationRecord
}

func NewService() *Service {
	now := time.Now()
	return &Service{
		history: map[string][]*LocationRecord{
			"user1": []*LocationRecord{
				{
					Coordinate: &types.Coordinate{
						Latitude:  51.11822470712269,
						Longitude: 1.990711729269563,
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
						Longitude: 16.390711729269563,
					},
					Timestamp: now.Add(-30 * time.Minute),
				},
			},
		},
	}
}

func (s *Service) RegisterLocation(ctx context.Context, userId string, coords *types.Coordinate, timestamp time.Time) ([]*LocationRecord, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	log.Println("Registering location...")

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
	s.history[userId] = user

	return user, nil
}

func (s *Service) CalculateDistance(ctx context.Context, userId string, startDate time.Time, endDate time.Time) (*DistanceRecord, error) {
	user, ok := s.history[userId]
	if !ok {
		return &DistanceRecord{
			distance: 0.0,
			history:  nil,
		}, fmt.Errorf("there is no user with such id: %v", userId)
	}

	var totalDistance float64
	var lastValidRecord *LocationRecord

	// Find the first location record within the time range
	for i := 0; i < len(user); i++ {
		if user[i].Timestamp.After(startDate) {
			lastValidRecord = user[i]
			break
		}
	}

	if lastValidRecord == nil {
		return &DistanceRecord{
			distance: 0.0,
			history:  user,
		}, nil
	}

	for i := 1; i < len(user); i++ {
		record := user[i]
		log.Println("Record Time: ", record.Timestamp)
		log.Println("StartDate: ", startDate)
		log.Println("EndTime: ", endDate)

		// Check if the current record is within the time range
		if record.Timestamp.After(startDate) && record.Timestamp.Before(endDate) {
			totalDistance += util.CalculateDistance(lastValidRecord.Coordinate, record.Coordinate)
			lastValidRecord = record
		}
	}

	log.Println("Total distance:", totalDistance)
	return &DistanceRecord{
		distance: totalDistance,
		history:  user,
	}, nil
}

func (d *DistanceRecord) ToProto() *pb.CalculateDistanceResponse {
	var protoLocationRecords []*pb.LocationRecord
	for _, u := range d.history {
		protoLocationRecords = append(protoLocationRecords, &pb.LocationRecord{
			Coordinate: &pb.Coordinate{
				Latitude:  u.Coordinate.Latitude,
				Longitude: u.Coordinate.Longitude,
			},
			Timestamp: u.Timestamp.String(),
		},
		)
	}
	return &pb.CalculateDistanceResponse{
		Distance: d.distance,
		History:  protoLocationRecords,
	}
}
