package grpc

import (
	"context"
	"go-clinet-locations/services/user-service/internal/domain"
	pb "go-clinet-locations/shared/proto/user"
	"go-clinet-locations/shared/types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
)

type grpcHandler struct {
	pb.UnimplementedUserServiceServer
	service domain.UserService
}

func NewGRPCHandler(server *grpc.Server, service domain.UserService) *grpcHandler {
	handler := &grpcHandler{
		service: service,
	}

	pb.RegisterUserServiceServer(server, handler)

	return handler
}
func (h *grpcHandler) CreateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.CreateUserResponse, error) {
	reqCoordinate := req.Coordinate

	userCords := &types.Coordinate{
		Longitude: reqCoordinate.Longitude,
		Latitude:  reqCoordinate.Latitude,
	}

	newUser := &domain.UserModel{
		UserName:    req.GetUserName(),
		Coordinates: userCords,
	}

	user, err := h.service.CreateUser(ctx, newUser)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create user %v", err)
	}
	log.Printf("user created with id: %v", user.ID)
	return &pb.CreateUserResponse{
		User: &pb.User{
			ID:       user.ID.String(),
			UserName: user.UserName,
			Coordinate: &pb.Coordinate{
				Latitude:  user.Coordinates.Latitude,
				Longitude: user.Coordinates.Longitude,
			},
		},
	}, nil

}

func (h *grpcHandler) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	reqCoordinate := req.GetCoordinate()

	userCords := &types.Coordinate{
		Longitude: reqCoordinate.Longitude,
		Latitude:  reqCoordinate.Latitude,
	}

	user, err := h.service.UpdateUser(ctx, req.GetUserName(), userCords)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update user %v", err)
	}

	return &pb.UpdateUserResponse{
		User: &pb.User{
			ID:       user.ID.String(),
			UserName: user.UserName,
			Coordinate: &pb.Coordinate{
				Latitude:  user.Coordinates.Latitude,
				Longitude: user.Coordinates.Longitude,
			},
		},
	}, nil
}

func (h *grpcHandler) SearchUsers(ctx context.Context, req *pb.SearchUsersRequest) (*pb.SearchUsersResponse, error) {
	reqCoordinate := req.GetCoordinate()

	coordinate := &types.Coordinate{
		Longitude: reqCoordinate.Longitude,
		Latitude:  reqCoordinate.Latitude,
	}
	users, err := h.service.SearchUsers(ctx, coordinate, float64(req.GetRadius()))

	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to search users %v", err)
	}

	return &pb.SearchUsersResponse{Users: domain.ToUsersProto(users)}, nil

}
