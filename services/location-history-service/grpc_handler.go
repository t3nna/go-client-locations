package main

import (
	"context"
	pb "go-clinet-locations/shared/proto/location"
	"go-clinet-locations/shared/types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"time"
)

type grpcHandler struct {
	service LocationsService
	pb.UnimplementedLocationServiceServer
}

func NewGrpcHandler(s *grpc.Server, service LocationsService) {
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

	locationRecords, err := h.service.RegisterLocation(ctx, req.UserId, coords, timestamp)

	locationRecordsProto := make([]*pb.LocationRecord, len(locationRecords))

	for _, r := range locationRecords {
		locationRecordsProto = append(locationRecordsProto, &pb.LocationRecord{
			Coordinate: &pb.Coordinate{
				Latitude:  r.Coordinate.Latitude,
				Longitude: r.Coordinate.Longitude,
			},
			Timestamp: r.Timestamp.String(),
		})
	}

	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to Register Location %v", err)
	}

	return &pb.RegisterLocationResponse{
		UserId:          req.GetUserId(),
		LocationRecords: locationRecordsProto,
	}, nil
}

func (h *grpcHandler) CalculateDistance(ctx context.Context, req *pb.CalculateDistanceRequest) (*pb.CalculateDistanceResponse, error) {
	var startDate string
	var endDate string

	isoLayout := time.RFC3339

	now := time.Now()
	if req.GetStartDate() != "" {
		startDate = req.GetStartDate()
		if req.GetEndDate() == "" {
			endDate = now.Format(isoLayout)
		} else {
			endDate = req.GetEndDate()
		}
	}

	if req.GetStartDate() == "" && req.GetEndDate() == "" {
		log.Println("Both empty")
		startDate = now.Add(-1 * 24 * time.Hour).Format(isoLayout)
		endDate = now.Format(isoLayout)
	}

	startDateParam, err := time.Parse(isoLayout, startDate)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to parse startDate: %v", err)
	}

	endDateParam, err := time.Parse(isoLayout, endDate)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to parse endDate: %v", err)
	}
	log.Println(startDate, endDate)

	distance, err := h.service.CalculateDistance(ctx, req.GetUserId(), startDateParam, endDateParam)
	if err != nil {

		return nil, status.Errorf(codes.Internal, "faild to calculate Distance: %v", err)
	}

	return distance.ToProto(), nil
}
