package testutil

import (
	"context"
	"go-clinet-locations/services/user-service/internal/domain"
	"go-clinet-locations/shared/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

// MockUserRepository is a mock implementation of UserRepository for testing
type MockUserRepository struct {
	users map[string]*domain.UserModel
}

// NewMockUserRepository creates a new mock repository
func NewMockUserRepository() *MockUserRepository {
	return &MockUserRepository{
		users: make(map[string]*domain.UserModel),
	}
}

// CreateUser mocks user creation
func (m *MockUserRepository) CreateUser(ctx context.Context, user *domain.UserModel) (*domain.UserModel, error) {
	user.ID = primitive.NewObjectID()
	m.users[user.ID.Hex()] = user
	return user, nil
}

// UpdateUser mocks user update
func (m *MockUserRepository) UpdateUser(ctx context.Context, userName string, coordinates *types.Coordinate) (*domain.UserModel, error) {
	for _, user := range m.users {
		if user.UserName == userName {
			user.Coordinates = coordinates
			return user, nil
		}
	}
	return nil, domain.ErrUserNotFound
}

// GetUsers mocks getting all users
func (m *MockUserRepository) GetUsers(ctx context.Context) ([]*domain.UserModel, error) {
	var users []*domain.UserModel
	for _, user := range m.users {
		users = append(users, user)
	}
	return users, nil
}

// SetUsers sets users in the mock repository
func (m *MockUserRepository) SetUsers(users []*domain.UserModel) {
	m.users = make(map[string]*domain.UserModel)
	for _, user := range users {
		m.users[user.ID.Hex()] = user
	}
}

// GetUserByID gets a user by ID from the mock repository
func (m *MockUserRepository) GetUserByID(id string) *domain.UserModel {
	return m.users[id]
}

// TestData contains test data for various scenarios
type TestData struct {
	ValidCoordinate   *types.Coordinate
	InvalidCoordinate *types.Coordinate
	ValidUsername     string
	InvalidUsername   string
	TestUsers         []*domain.UserModel
}

// GetTestData returns test data for various scenarios
func GetTestData() *TestData {
	return &TestData{
		ValidCoordinate: &types.Coordinate{
			Latitude:  51.11822470712269,
			Longitude: 16.990711729269563,
		},
		InvalidCoordinate: &types.Coordinate{
			Latitude:  91.0,  // Invalid latitude
			Longitude: 181.0, // Invalid longitude
		},
		ValidUsername:   "testuser123",
		InvalidUsername: "ab", // Too short
		TestUsers: []*domain.UserModel{
			{
				ID:       primitive.NewObjectID(),
				UserName: "user1",
				Coordinates: &types.Coordinate{
					Latitude:  51.11822470712269,
					Longitude: 16.990711729269563,
				},
			},
			{
				ID:       primitive.NewObjectID(),
				UserName: "user2",
				Coordinates: &types.Coordinate{
					Latitude:  52.23553956649786,
					Longitude: 20.984595191389918,
				},
			},
		},
	}
}

// CreateTestUser creates a test user with given parameters
func CreateTestUser(username string, lat, lon float64) *domain.UserModel {
	return &domain.UserModel{
		ID:       primitive.NewObjectID(),
		UserName: username,
		Coordinates: &types.Coordinate{
			Latitude:  lat,
			Longitude: lon,
		},
	}
}

// CreateTestCoordinate creates a test coordinate
func CreateTestCoordinate(lat, lon float64) *types.Coordinate {
	return &types.Coordinate{
		Latitude:  lat,
		Longitude: lon,
	}
}

// CreateTestContext creates a test context with timeout
func CreateTestContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 5*time.Second)
}
