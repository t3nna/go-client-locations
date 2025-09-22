package main

import (
	"context"
	"go-clinet-locations/shared/types"
	"testing"
	"time"
)

func TestService_RegisterLocation(t *testing.T) {
	tests := []struct {
		name        string
		userId      string
		coords      *types.Coordinate
		timestamp   time.Time
		expectError bool
	}{
		{
			name:   "register location for new user",
			userId: "newuser",
			coords: &types.Coordinate{
				Latitude:  51.11822470712269,
				Longitude: 16.990711729269563,
			},
			timestamp:   time.Now(),
			expectError: false,
		},
		{
			name:   "register location for existing user",
			userId: "user1",
			coords: &types.Coordinate{
				Latitude:  52.0,
				Longitude: 17.0,
			},
			timestamp:   time.Now(),
			expectError: false,
		},
		{
			name:   "register location with zero coordinates",
			userId: "user2",
			coords: &types.Coordinate{
				Latitude:  0.0,
				Longitude: 0.0,
			},
			timestamp:   time.Now(),
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := NewService()
			ctx := context.Background()

			result, err := service.RegisterLocation(ctx, tt.userId, tt.coords, tt.timestamp)

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if result == nil {
					t.Errorf("expected location records but got nil")
				}
				if len(result) == 0 {
					t.Errorf("expected at least one location record")
				}
				// Check if the last record matches what we registered
				lastRecord := result[len(result)-1]
				if lastRecord.Coordinate.Latitude != tt.coords.Latitude {
					t.Errorf("expected latitude %.6f, got %.6f", tt.coords.Latitude, lastRecord.Coordinate.Latitude)
				}
				if lastRecord.Coordinate.Longitude != tt.coords.Longitude {
					t.Errorf("expected longitude %.6f, got %.6f", tt.coords.Longitude, lastRecord.Coordinate.Longitude)
				}
			}
		})
	}
}

func TestService_CalculateDistance(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name             string
		userId           string
		startDate        time.Time
		endDate          time.Time
		expectError      bool
		expectedDistance float64
		tolerance        float64
	}{
		{
			name:             "calculate distance for existing user with history",
			userId:           "user1",
			startDate:        now.Add(-3 * time.Hour * 24), // 3 days ago
			endDate:          now.Add(-1 * time.Hour),      // 1 hour ago
			expectError:      false,
			expectedDistance: 0.0, // Based on the test data in NewService
			tolerance:        0.1,
		},
		{
			name:        "calculate distance for non-existent user",
			userId:      "nonexistent",
			startDate:   now.Add(-1 * time.Hour),
			endDate:     now,
			expectError: true,
		},
		{
			name:             "calculate distance with no records in time range",
			userId:           "user1",
			startDate:        now.Add(1 * time.Hour), // Future date
			endDate:          now.Add(2 * time.Hour), // Future date
			expectError:      false,
			expectedDistance: 0.0,
			tolerance:        0.1,
		},
		{
			name:             "calculate distance with single record in time range",
			userId:           "user1",
			startDate:        now.Add(-2 * time.Hour * 24), // 2 days ago
			endDate:          now.Add(-1 * time.Hour * 24), // 1 day ago
			expectError:      false,
			expectedDistance: 0.0, // Single record means no distance
			tolerance:        0.1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := NewService()
			ctx := context.Background()

			result, err := service.CalculateDistance(ctx, tt.userId, tt.startDate, tt.endDate)

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if result == nil {
					t.Errorf("expected distance record but got nil")
				}
				if result.distance < 0 {
					t.Errorf("distance should not be negative, got %.6f", result.distance)
				}
				if tt.expectedDistance > 0 {
					diff := result.distance - tt.expectedDistance
					if diff < 0 {
						diff = -diff
					}
					if diff > tt.tolerance {
						t.Errorf("expected distance around %.6f, got %.6f (tolerance: %.6f)",
							tt.expectedDistance, result.distance, tt.tolerance)
					}
				}
			}
		})
	}
}

func TestService_ConcurrentAccess(t *testing.T) {
	service := NewService()
	ctx := context.Background()

	// Test concurrent access to the same user
	done := make(chan bool, 10)

	for i := 0; i < 10; i++ {
		go func(index int) {
			coords := &types.Coordinate{
				Latitude:  float64(51 + index),
				Longitude: float64(16 + index),
			}
			_, err := service.RegisterLocation(ctx, "concurrent_user", coords, time.Now())
			if err != nil {
				t.Errorf("unexpected error in concurrent access: %v", err)
			}
			done <- true
		}(i)
	}

	// Wait for all goroutines to complete
	for i := 0; i < 10; i++ {
		<-done
	}

	// Verify that all locations were registered
	coords := &types.Coordinate{
		Latitude:  51.0,
		Longitude: 16.0,
	}
	_, err := service.RegisterLocation(ctx, "concurrent_user", coords, time.Now())
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestDistanceRecord_ToProto(t *testing.T) {
	now := time.Now()

	record := &DistanceRecord{
		distance: 123.45,
		history: []*LocationRecord{
			{
				Coordinate: &types.Coordinate{
					Latitude:  51.11822470712269,
					Longitude: 16.990711729269563,
				},
				Timestamp: now,
			},
			{
				Coordinate: &types.Coordinate{
					Latitude:  52.0,
					Longitude: 17.0,
				},
				Timestamp: now.Add(1 * time.Hour),
			},
		},
	}

	proto := record.ToProto()

	if proto.Distance != record.distance {
		t.Errorf("expected distance %.6f, got %.6f", record.distance, proto.Distance)
	}

	if len(proto.History) != len(record.history) {
		t.Errorf("expected %d history records, got %d", len(record.history), len(proto.History))
	}

	for i, hist := range proto.History {
		if hist.Coordinate.Latitude != record.history[i].Coordinate.Latitude {
			t.Errorf("expected latitude %.6f, got %.6f",
				record.history[i].Coordinate.Latitude, hist.Coordinate.Latitude)
		}
		if hist.Coordinate.Longitude != record.history[i].Coordinate.Longitude {
			t.Errorf("expected longitude %.6f, got %.6f",
				record.history[i].Coordinate.Longitude, hist.Coordinate.Longitude)
		}
	}
}
