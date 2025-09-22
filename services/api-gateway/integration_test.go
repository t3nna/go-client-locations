//go:build integration

package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"go-clinet-locations/shared/contracts"
	"go-clinet-locations/shared/types"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// Integration tests require the full system to be running
// Run with: go test -tags=integration

func TestIntegration_CreateUserFlow(t *testing.T) {
	// This test requires the full system to be running
	// including user-service and location-history-service

	tests := []struct {
		name           string
		requestBody    userLocationRequest
		expectedStatus int
	}{
		{
			name: "create user with valid data",
			requestBody: userLocationRequest{
				UserName: "integrationuser",
				Coordinate: types.Coordinate{
					Latitude:  51.11822470712269,
					Longitude: 16.990711729269563,
				},
			},
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var body bytes.Buffer
			json.NewEncoder(&body).Encode(tt.requestBody)

			req := httptest.NewRequest("POST", "/user/create", &body)
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			HandleCreateUser(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if w.Code == http.StatusOK {
				var response contracts.APIResponse
				if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
					t.Errorf("failed to decode response: %v", err)
				}
				if response.Data == nil {
					t.Errorf("expected response data but got nil")
				}
			}
		})
	}
}

func TestIntegration_UpdateUserFlow(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    userLocationRequest
		expectedStatus int
	}{
		{
			name: "update user with valid data",
			requestBody: userLocationRequest{
				UserName: "integrationuser",
				Coordinate: types.Coordinate{
					Latitude:  52.0,
					Longitude: 17.0,
				},
			},
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var body bytes.Buffer
			json.NewEncoder(&body).Encode(tt.requestBody)

			req := httptest.NewRequest("PATCH", "/user/update", &body)
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			HandleUpdateUser(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}

func TestIntegration_SearchUserFlow(t *testing.T) {
	tests := []struct {
		name           string
		queryParams    map[string]string
		expectedStatus int
	}{
		{
			name: "search users with valid coordinates",
			queryParams: map[string]string{
				"lat": "51.11822470712269",
				"lon": "16.990711729269563",
				"r":   "10.0",
			},
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/user/search", nil)
			q := req.URL.Query()
			for key, value := range tt.queryParams {
				q.Add(key, value)
			}
			req.URL.RawQuery = q.Encode()
			w := httptest.NewRecorder()

			HandleSearchUser(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if w.Code == http.StatusOK {
				var response contracts.APIResponse
				if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
					t.Errorf("failed to decode response: %v", err)
				}
			}
		})
	}
}

func TestIntegration_CalculateDistanceFlow(t *testing.T) {
	tests := []struct {
		name           string
		queryParams    map[string]string
		expectedStatus int
	}{
		{
			name: "calculate distance with valid userId",
			queryParams: map[string]string{
				"userId": "integrationuser",
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "calculate distance with time range",
			queryParams: map[string]string{
				"userId":    "integrationuser",
				"startTime": "2023-01-01T00:00:00Z",
				"endTime":   "2023-12-31T23:59:59Z",
			},
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/user/distance", nil)
			q := req.URL.Query()
			for key, value := range tt.queryParams {
				q.Add(key, value)
			}
			req.URL.RawQuery = q.Encode()
			w := httptest.NewRecorder()

			HandleCalculateDistance(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if w.Code == http.StatusOK {
				var response contracts.APIResponse
				if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
					t.Errorf("failed to decode response: %v", err)
				}
			}
		})
	}
}

func TestIntegration_EndToEndFlow(t *testing.T) {
	// This test performs a complete end-to-end flow:
	// 1. Create a user
	// 2. Update the user's location
	// 3. Search for users near the location
	// 4. Calculate distance traveled

	ctx := context.Background()
	timeout := 30 * time.Second
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	// Step 1: Create a user
	t.Run("create_user", func(t *testing.T) {
		requestBody := userLocationRequest{
			UserName: "e2euser",
			Coordinate: types.Coordinate{
				Latitude:  51.11822470712269,
				Longitude: 16.990711729269563,
			},
		}

		var body bytes.Buffer
		json.NewEncoder(&body).Encode(requestBody)

		req := httptest.NewRequest("POST", "/user/create", &body)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		HandleCreateUser(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("create user failed with status %d", w.Code)
		}
	})

	// Step 2: Update user location
	t.Run("update_user", func(t *testing.T) {
		requestBody := userLocationRequest{
			UserName: "e2euser",
			Coordinate: types.Coordinate{
				Latitude:  52.0,
				Longitude: 17.0,
			},
		}

		var body bytes.Buffer
		json.NewEncoder(&body).Encode(requestBody)

		req := httptest.NewRequest("PATCH", "/user/update", &body)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		HandleUpdateUser(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("update user failed with status %d", w.Code)
		}
	})

	// Step 3: Search for users
	t.Run("search_users", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/user/search", nil)
		q := req.URL.Query()
		q.Add("lat", "51.5")
		q.Add("lon", "16.5")
		q.Add("r", "50.0")
		req.URL.RawQuery = q.Encode()
		w := httptest.NewRecorder()

		HandleSearchUser(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("search users failed with status %d", w.Code)
		}
	})

	// Step 4: Calculate distance
	t.Run("calculate_distance", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/user/distance", nil)
		q := req.URL.Query()
		q.Add("userId", "e2euser")
		req.URL.RawQuery = q.Encode()
		w := httptest.NewRecorder()

		HandleCalculateDistance(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("calculate distance failed with status %d", w.Code)
		}
	})
}

func TestIntegration_ConcurrentRequests(t *testing.T) {
	// Test concurrent requests to ensure thread safety
	concurrency := 10
	done := make(chan bool, concurrency)

	for i := 0; i < concurrency; i++ {
		go func(index int) {
			requestBody := userLocationRequest{
				UserName: fmt.Sprintf("concurrentuser%d", index),
				Coordinate: types.Coordinate{
					Latitude:  float64(51 + index),
					Longitude: float64(16 + index),
				},
			}

			var body bytes.Buffer
			json.NewEncoder(&body).Encode(requestBody)

			req := httptest.NewRequest("POST", "/user/create", &body)
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			HandleCreateUser(w, req)

			if w.Code != http.StatusOK {
				t.Errorf("concurrent request %d failed with status %d", index, w.Code)
			}

			done <- true
		}(i)
	}

	// Wait for all goroutines to complete
	for i := 0; i < concurrency; i++ {
		<-done
	}
}
