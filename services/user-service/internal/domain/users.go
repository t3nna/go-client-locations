package domain

import (
	"context"
	"go-clinet-locations/shared/types"
)

type UserModel struct {
	UserId      string
	UserName    string
	Coordinates *types.Coordinate
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
