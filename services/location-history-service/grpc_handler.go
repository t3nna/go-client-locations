package main

import (
	"context"
	pb "go-clinet-locations/shared/proto/location"
	"go-clinet-locations/shared/types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

type grpcHandler struct {
	service *Service
	pb.UnimplementedLocationServiceServer
}

func NewGrpcHandler(s *grpc.Server, service *Service) {
	handler := &grpcHandler{
		service: service,
	}
	pb.RegisterLocationServiceServer(s, handler)
}

func (h *grpcHandler) RegisterLocation(ctx context.Context, req *pb.RegisterLocationRequest) (*pb.RegisterLocationResponse, error) {
	coords := &types.Coordinate{
		Latitude:  req.GetCoordinate().Latitude,
		Longitude: req.GetCoordinate().Longitude,
	}
	isoLayout := time.RFC3339
	timestamp, err := time.Parse(isoLayout, req.GetTimestamp())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to parse timestamp: %v", err)
	}

	// TODO: return location records, map it to pb type
	_, err = h.service.RegisterLocation(req.UserId, coords, timestamp)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to Register Location")
	}

	return &pb.RegisterLocationResponse{
		UserId:          req.GetUserId(),
		LocationRecords: nil,
	}, nil
}
