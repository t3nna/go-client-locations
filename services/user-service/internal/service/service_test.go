package service

import (
	"context"
	"go-clinet-locations/services/user-service/internal/domain"
	"go-clinet-locations/services/user-service/internal/testutil"
	"go-clinet-locations/shared/types"
	"testing"
)

func TestService_CreateUser(t *testing.T) {
	tests := []struct {
		name        string
		user        *domain.UserModel
		expectError bool
	}{
		{
			name: "valid user creation",
			user: &domain.UserModel{
				UserName: "testuser",
				Coordinates: &types.Coordinate{
					Latitude:  51.11822470712269,
					Longitude: 16.990711729269563,
				},
			},
			expectError: false,
		},
		{
			name: "user with empty username",
			user: &domain.UserModel{
				UserName: "",
				Coordinates: &types.Coordinate{
					Latitude:  51.11822470712269,
					Longitude: 16.990711729269563,
				},
			},
			expectError: false, // Repository should handle this
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			mockRepo := testutil.NewMockUserRepository()
			service := NewService(mockRepo)

			result, err := service.CreateUser(ctx, tt.user)

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if result == nil {
					t.Errorf("expected user but got nil")
				}
				if result.UserName != tt.user.UserName {
					t.Errorf("expected username %s, got %s", tt.user.UserName, result.UserName)
				}
			}
		})
	}
}

func TestService_UpdateUser(t *testing.T) {
	tests := []struct {
		name        string
		username    string
		coordinates *types.Coordinate
		setupUsers  []*domain.UserModel
		expectError bool
	}{
		{
			name:     "valid user update",
			username: "user1",
			coordinates: &types.Coordinate{
				Latitude:  52.0,
				Longitude: 17.0,
			},
			setupUsers: []*domain.UserModel{
				testutil.CreateTestUser("user1", 51.0, 16.0),
			},
			expectError: false,
		},
		{
			name:     "update non-existent user",
			username: "nonexistent",
			coordinates: &types.Coordinate{
				Latitude:  52.0,
				Longitude: 17.0,
			},
			setupUsers: []*domain.UserModel{
				testutil.CreateTestUser("user1", 51.0, 16.0),
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			mockRepo := testutil.NewMockUserRepository()
			mockRepo.SetUsers(tt.setupUsers)
			service := NewService(mockRepo)

			result, err := service.UpdateUser(ctx, tt.username, tt.coordinates)

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if result == nil {
					t.Errorf("expected user but got nil")
				}
				if result.UserName != tt.username {
					t.Errorf("expected username %s, got %s", tt.username, result.UserName)
				}
				if result.Coordinates.Latitude != tt.coordinates.Latitude {
					t.Errorf("expected latitude %.6f, got %.6f", tt.coordinates.Latitude, result.Coordinates.Latitude)
				}
				if result.Coordinates.Longitude != tt.coordinates.Longitude {
					t.Errorf("expected longitude %.6f, got %.6f", tt.coordinates.Longitude, result.Coordinates.Longitude)
				}
			}
		})
	}
}

func TestService_SearchUsers(t *testing.T) {
	tests := []struct {
		name          string
		location      *types.Coordinate
		radius        float64
		setupUsers    []*domain.UserModel
		expectedCount int
		expectError   bool
	}{
		{
			name: "search within radius - should find users",
			location: &types.Coordinate{
				Latitude:  51.11822470712269,
				Longitude: 16.990711729269563,
			},
			radius: 10.0, // 10 km radius
			setupUsers: []*domain.UserModel{
				testutil.CreateTestUser("user1", 51.11822470712269, 16.990711729269563), // Same location
				testutil.CreateTestUser("user2", 51.11956092410769, 17.05696305051491),  // ~4.5 km away
				testutil.CreateTestUser("user3", 52.23553956649786, 20.984595191389918), // ~300 km away
			},
			expectedCount: 2, // user1 and user2 should be found
			expectError:   false,
		},
		{
			name: "search with small radius - should find only exact match",
			location: &types.Coordinate{
				Latitude:  51.11822470712269,
				Longitude: 16.990711729269563,
			},
			radius: 0.1, // 0.1 km radius
			setupUsers: []*domain.UserModel{
				testutil.CreateTestUser("user1", 51.11822470712269, 16.990711729269563), // Same location
				testutil.CreateTestUser("user2", 51.11956092410769, 17.05696305051491),  // ~4.5 km away
			},
			expectedCount: 1, // Only user1 should be found
			expectError:   false,
		},
		{
			name: "search with large radius - should find all users",
			location: &types.Coordinate{
				Latitude:  51.11822470712269,
				Longitude: 16.990711729269563,
			},
			radius: 1000.0, // 1000 km radius
			setupUsers: []*domain.UserModel{
				testutil.CreateTestUser("user1", 51.11822470712269, 16.990711729269563), // Same location
				testutil.CreateTestUser("user2", 52.23553956649786, 20.984595191389918), // ~300 km away
				testutil.CreateTestUser("user3", 50.53401932980686, 31.178889172903055), // ~1000 km away
			},
			expectedCount: 3, // All users should be found
			expectError:   false,
		},
		{
			name: "search with no users in database",
			location: &types.Coordinate{
				Latitude:  51.11822470712269,
				Longitude: 16.990711729269563,
			},
			radius:        10.0,
			setupUsers:    []*domain.UserModel{},
			expectedCount: 0,
			expectError:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			mockRepo := testutil.NewMockUserRepository()
			mockRepo.SetUsers(tt.setupUsers)
			service := NewService(mockRepo)

			result, err := service.SearchUsers(ctx, tt.location, tt.radius)

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if len(result) != tt.expectedCount {
					t.Errorf("expected %d users, got %d", tt.expectedCount, len(result))
				}
			}
		})
	}
}
