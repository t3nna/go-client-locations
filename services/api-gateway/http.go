package main

import (
	"encoding/json"
	"fmt"
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

	if reqBody.Location.Longitude == 0 || reqBody.Location.Latitude == 0 {
		http.Error(w, "invalid location data", http.StatusBadRequest)
		return
	}

	fmt.Printf("User location handler")
	res := map[string]string{"Kinda": "success"}

	writeJSON(w, http.StatusOK, res)
}
