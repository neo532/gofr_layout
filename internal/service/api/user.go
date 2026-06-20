package api

import (
	"context"

	pb "github.com/neo532/gofr_layout/proto/api/user/v1"
)

// UserApiService implements pb.UserApiService.
type UserApiService struct{}

func NewUserApiService() *UserApiService {
	return &UserApiService{}
}

func (s *UserApiService) Post(ctx context.Context, req *pb.User) (*pb.User, error) {
	return req, nil
}

func (s *UserApiService) GetById(ctx context.Context, req *pb.GetByIdRequest) (*pb.User, error) {
	return &pb.User{Name: "Hello"}, nil
}
