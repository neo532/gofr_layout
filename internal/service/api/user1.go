package api

import (
	"context"

	pb "github.com/neo532/gofr_layout/proto/api/user/v1"
)

// UserApi1Service implements pb.UserApi1Service.
type UserApi1Service struct{}

func NewUserApi1Service() *UserApi1Service {
	return &UserApi1Service{}
}
func (s *UserApi1Service) Post(ctx context.Context, req *pb.User) (*pb.User, error) {
	return req, nil
}

func (s *UserApi1Service) GetById(ctx context.Context, req *pb.GetByIdRequest) (*pb.User, error) {
	return &pb.User{Name: "Hello11"}, nil
}
