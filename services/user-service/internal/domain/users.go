package domain

import (
	"context"
	pb "go-clinet-locations/shared/proto/user"
	"go-clinet-locations/shared/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserModel struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	UserId      string             `bson:"userId"`
	UserName    string             `bson:"userName"`
	Coordinates *types.Coordinate  `bson:"coordinates"`
}

type UserRepository interface {
	CreateUser(ctx context.Context, user *UserModel) (*UserModel, error)
	UpdateUser(ctx context.Context, userName string, coordinates *types.Coordinate) (*UserModel, error)
	GetUsers(ctx context.Context) ([]*UserModel, error)
}

type UserService interface {
	CreateUser(ctx context.Context, user *UserModel) (*UserModel, error)
	UpdateUser(ctx context.Context, userName string, coordinates *types.Coordinate) (*UserModel, error)
	SearchUsers(ctx context.Context, location *types.Coordinate, radius float64) ([]*UserModel, error)
}

func (u *UserModel) ToProto() *pb.User {
	return &pb.User{
		UserId:   u.UserId,
		UserName: u.UserName,
		Coordinate: &pb.Coordinate{
			Latitude:  u.Coordinates.Latitude,
			Longitude: u.Coordinates.Longitude,
		},
	}
}

func ToUsersProto(users []*UserModel) []*pb.User {
	var protoUsers []*pb.User
	for _, u := range users {
		protoUsers = append(protoUsers, u.ToProto())
	}
	return protoUsers
}
