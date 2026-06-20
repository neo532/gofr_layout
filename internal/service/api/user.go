package api

import (
	"context"

	"github.com/neo532/gofr_layout/internal/domain"
	pb "github.com/neo532/gofr_layout/proto/api/user/v1"
)

// UserApiService implements pb.UserApiService.
type UserApiService struct {
	dn *domain.UserDomain
}

func NewUserApiService(
	dn *domain.UserDomain,
) *UserApiService {
	return &UserApiService{
		dn: dn,
	}
}

func (s *UserApiService) Post(ctx context.Context, req *pb.User) (*pb.User, error) {
	return req, nil
}

func (s *UserApiService) GetById(ctx context.Context, req *pb.GetByIdRequest) (*pb.User, error) {
	return s.dn.GetById(ctx, req.Id)
}
