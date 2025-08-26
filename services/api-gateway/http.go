package main

import (
	"context"
	"encoding/json"
	"go-clinet-locations/services/api-gateway/grpc_clients"
	"go-clinet-locations/shared/contracts"
	pb_loction "go-clinet-locations/shared/proto/location"
	pb_user "go-clinet-locations/shared/proto/user"
	"go-clinet-locations/shared/util"
	"log"
	"net/http"
	"strconv"
	"time"
)

func HandleCreateUser(w http.ResponseWriter, r *http.Request) {
	var reqBody userLocationRequest

	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		log.Println(err)
		http.Error(w, "failed to parse JSON data", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// validation

	if reqBody.UserName == "" {
		http.Error(w, "failed to parse JSON data", http.StatusBadRequest)
		return
	}

	if reqBody.Coordinate.Longitude == 0 || reqBody.Coordinate.Latitude == 0 {
		http.Error(w, "invalid location data", http.StatusBadRequest)
		return
	}

	userService, err := grpc_clients.NewUserServiceClient()

	if err != nil {
		log.Fatal(err)
	}

	defer userService.Close()

	newUser, err := userService.Client.CreateUser(r.Context(), reqBody.toProto())
	if err != nil {
		log.Printf("Failed to create a user: %v", err)
		http.Error(w, "Failed to create a user", http.StatusInternalServerError)
		return

	}

	now := time.Now()
	isoNow := now.Format(time.RFC3339)

	// Calling location history service
	err = handleRegisterLocation(r.Context(), &pb_loction.RegisterLocationRequest{
		UserId: newUser.User.UserId,
		Coordinate: &pb_loction.Coordinate{
			Latitude:  newUser.User.Coordinate.Latitude,
			Longitude: newUser.User.Coordinate.Longitude,
		},
		Timestamp: isoNow,
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return

	}

	response := contracts.APIResponse{Data: newUser}

	writeJSON(w, http.StatusOK, response)
}
func HandleUpdateUser(w http.ResponseWriter, r *http.Request) {
	var reqBody userLocationRequest

	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		log.Println(err)
		http.Error(w, "failed to parse JSON data", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// validation

	if reqBody.UserName == "" {
		http.Error(w, "failed to parse JSON data", http.StatusBadRequest)
		return
	}

	if reqBody.Coordinate.Longitude == 0 || reqBody.Coordinate.Latitude == 0 {
		http.Error(w, "invalid location data", http.StatusBadRequest)
		return
	}
	userService, err := grpc_clients.NewUserServiceClient()

	if err != nil {
		log.Fatal(err)
	}

	defer userService.Close()

	newUser, err := userService.Client.UpdateUser(r.Context(), reqBody.toProto())

	if err != nil {
		log.Printf("Failed to update a user: %v", err)
		http.Error(w, "Failed to update a user", http.StatusInternalServerError)
		return

	}

	// TODO: move timestamp to location service
	now := time.Now()
	isoNow := now.Format(time.RFC3339)

	// Calling location history service
	err = handleRegisterLocation(r.Context(), &pb_loction.RegisterLocationRequest{
		UserId: newUser.User.UserId,
		Coordinate: &pb_loction.Coordinate{
			Latitude:  newUser.User.Coordinate.Latitude,
			Longitude: newUser.User.Coordinate.Longitude,
		},
		Timestamp: isoNow,
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return

	}

	response := contracts.APIResponse{Data: newUser}

	writeJSON(w, http.StatusOK, response)

}

func HandleSearchUser(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	log.Println("Handle Search User ")

	// Coordinates validation

	if len(q["lat"]) != 1 {
		http.Error(w, "failed to retrieve coordinates", http.StatusBadRequest)
		return
	}
	if len(q["lon"]) != 1 {
		http.Error(w, "failed to retrieve coordinates", http.StatusBadRequest)
		return
	}
	lat := q["lat"][0]
	lon := q["lon"][0]

	latitude, err := strconv.ParseFloat(lat, 64)
	if err != nil {
		http.Error(w, "failed to parse latitude", http.StatusBadRequest)
		return
	}

	longitude, err := strconv.ParseFloat(lon, 64)

	if err != nil {
		http.Error(w, "failed to parse longitude", http.StatusBadRequest)
		return
	}

	err = util.ValidateCords(latitude, longitude)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	radius := 5.0
	if len(q["r"]) > 0 {
		if rad := q["r"][0]; rad != "" {
			r, err := strconv.ParseFloat(rad, 64)
			if err != nil {
				http.Error(w, "failed to parse radius", http.StatusBadRequest)
			}
			radius = r
		}
	}

	userService, err := grpc_clients.NewUserServiceClient()

	if err != nil {
		log.Fatal(err)
	}

	defer userService.Close()

	filteredUsers, err := userService.Client.SearchUsers(r.Context(), &pb_user.SearchUsersRequest{
		Coordinate: &pb_user.Coordinate{
			Latitude:  latitude,
			Longitude: longitude,
		},
		Radius: float32(radius),
	})

	if err != nil {
		log.Printf("Failed to search users: %v", err)
		http.Error(w, "Failed to search users", http.StatusInternalServerError)
		return
	}

	res := contracts.APIResponse{Data: filteredUsers}

	writeJSON(w, http.StatusOK, res)
}

func HandleCalculateDistance(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	userId := q["userId"]
	startTime := q["startTime"]
	endTime := q["endTime"]

	// validation

	if len(userId) != 1 {
		http.Error(w, "something wrong with userId param", http.StatusBadRequest)
		return
	}

	if userId[0] == "" {
		http.Error(w, "userId is missing", http.StatusBadRequest)
	}

	if len(startTime) != 1 {
		http.Error(w, "something wrong with startTime param", http.StatusBadRequest)
		return
	}

	if startTime[0] == "" {
		http.Error(w, "startTime is missing", http.StatusBadRequest)
		return
	}

	if len(endTime) != 1 {
		http.Error(w, "something wrong with endTime param", http.StatusBadRequest)
		return
	}

	if endTime[0] == "" {
		http.Error(w, "endTime is missing", http.StatusBadRequest)
		return
	}

	// TODO make start and end time optional

	userService, err := grpc_clients.NewLocationServiceClient()

	if err != nil {
		log.Fatal(err)
	}

	defer userService.Close()

	distance, err := userService.Client.CalculateDistance(r.Context(), &pb_loction.CalculateDistanceRequest{
		UserId:    userId[0],
		StartDate: startTime[0],
		EndDate:   endTime[0],
	})
	if err != nil {
		log.Printf("Failed to calculate distance: %v", err)
		http.Error(w, "Failed to calculate distance", http.StatusInternalServerError)
		return

	}

	res := contracts.APIResponse{Data: distance}

	writeJSON(w, http.StatusOK, res)

}

func handleRegisterLocation(ctx context.Context, req *pb_loction.RegisterLocationRequest) error {
	locationService, err := grpc_clients.NewLocationServiceClient()
	if err != nil {
		log.Fatal(err)
	}
	defer locationService.Close()
	newLocationRecord, err := locationService.Client.RegisterLocation(ctx, req)

	if err != nil {
		log.Printf("Failed to register user location: %v", err)
		return err
	}

	log.Println(newLocationRecord)

	return nil
}
