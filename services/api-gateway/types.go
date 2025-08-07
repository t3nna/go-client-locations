package main

import "go-clinet-locations/shared/types"

type userLocationRequest struct {
	UserName string           `json:"userName"`
	Location types.Coordinate `json:"location"`
}
