package main

import (
	pb "go-clinet-locations/shared/proto/user"
	"go-clinet-locations/shared/types"
)

type userLocationRequest struct {
	UserName   string           `json:"userName"`
	Coordinate types.Coordinate `json:"coordinate"`
}

func (userLocation *userLocationRequest) toProto() *pb.UpdateUserRequest {
	return &pb.UpdateUserRequest{
		UserName: userLocation.UserName,
		Coordinate: &pb.Coordinate{
			Latitude:  userLocation.Coordinate.Latitude,
			Longitude: userLocation.Coordinate.Longitude,
		},
	}

}

type calculateDistanceRequest struct {
	UserId    string `json:"userId"`
	DateRange string `json:"dateRange"`
}
