package main

import (
	"encoding/json"
	"go-clinet-locations/services/api-gateway/grpc_clients"
	"go-clinet-locations/shared/contracts"
	"log"
	"net/http"
)

func HandleUpdateUserLocation(w http.ResponseWriter, r *http.Request) {
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
		log.Printf("Failed to preveiw a trip: %v", err)
		http.Error(w, "Failed to Preview trip", http.StatusInternalServerError)
		return

	}

	response := contracts.APIResponse{Data: newUser}

	writeJSON(w, http.StatusOK, response)
}
