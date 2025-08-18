package grpc

import (
	"context"
	"go-clinet-locations/services/trip-service/internal/domain"
	pb "go-clinet-locations/shared/proto/user"
	"go-clinet-locations/shared/types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
func (h *grpcHandler) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
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
	return &pb.CreateUserResponse{
		UserId:   user.UserId,
		UserName: user.UserName,
		Coordinate: &pb.Coordinate{
			Latitude:  user.Coordinates.Latitude,
			Longitude: user.Coordinates.Longitude,
		},
	}, nil

}
