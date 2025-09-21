package main

import (
	"context"
	"fmt"
	"go-clinet-locations/shared/db"
	"go-clinet-locations/shared/types"
	"go-clinet-locations/shared/util"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"sync"
	"time"
)

type mongoService struct {
	db *mongo.Database
	mu sync.RWMutex
}

func NewMongoService(db *mongo.Database) *mongoService {
	return &mongoService{db: db}
}

func (m *mongoService) RegisterLocation(ctx context.Context, userId string, coords *types.Coordinate, timestamp time.Time) ([]*LocationRecord, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	collection := m.db.Collection(db.LocationCollection)

	// Convert the string userId to a primitive.ObjectID
	objID, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID format: %v", err)
	}

	locationRecord := &LocationRecord{
		Coordinate: coords,
		Timestamp:  timestamp,
	}

	filter := bson.M{"_id": objID} // Use the converted ObjectID
	update := bson.M{"$push": bson.M{"history": locationRecord}}

	result, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, fmt.Errorf("failed to register location: %v", err)
	}

	// Check if a document was actually updated
	if result.ModifiedCount == 0 {
		log.Println("User not found, attempting to insert new document.")

		// If no document was updated, it means the user doesn't exist.
		// In this case, we'll create a new document.
		newDoc := bson.M{
			"_id":     objID,
			"history": []*LocationRecord{locationRecord},
		}

		_, err := collection.InsertOne(ctx, newDoc)
		if err != nil {
			return nil, fmt.Errorf("failed to insert new user document: %v", err)
		}
	}

	// Retrieve updated history
	var user struct {
		History []*LocationRecord `bson:"history"`
	}
	err = collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve updated history: %v", err)
	}

	return user.History, nil
}

func (m *mongoService) CalculateDistance(ctx context.Context, userId string, startDate time.Time, endDate time.Time) (*DistanceRecord, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	collection := m.db.Collection(db.LocationCollection)

	// 1. Convert the string userId to a primitive.ObjectID
	objID, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID format: %v", err)
	}

	// 2. Define the filter to retrieve the user's history
	// We only need the history field, so we use a projection.
	filter := bson.M{"_id": objID}
	projection := bson.M{"_id": 0, "history": 1}

	var userDoc struct {
		History []*LocationRecord `bson:"history"`
	}

	err = collection.FindOne(ctx, filter, options.FindOne().SetProjection(projection)).Decode(&userDoc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return &DistanceRecord{distance: 0.0, history: nil}, nil
		}
		return nil, fmt.Errorf("failed to retrieve user history: %v", err)
	}

	// 3. Filter the retrieved history in-memory to match the date range
	var filteredHistory []*LocationRecord
	for _, record := range userDoc.History {
		if record.Timestamp.After(startDate) && record.Timestamp.Before(endDate) {
			filteredHistory = append(filteredHistory, record)
		}
	}

	// Handle case where no records are in the time range
	if len(filteredHistory) < 2 {
		return &DistanceRecord{distance: 0.0, history: filteredHistory}, nil
	}

	// 4. Calculate the distance between consecutive points in the filtered history
	var totalDistance float64
	for i := 0; i < len(filteredHistory)-1; i++ {
		currentRecord := filteredHistory[i]
		nextRecord := filteredHistory[i+1]

		// This is where you would use your util.CalculateDistance function
		totalDistance += util.CalculateDistance(currentRecord.Coordinate, nextRecord.Coordinate)
	}

	return &DistanceRecord{
		distance: totalDistance,
		history:  filteredHistory,
	}, nil
}
