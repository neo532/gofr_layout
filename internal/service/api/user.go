package api

import (
	"context"

	"github.com/neo532/gofr_layout/internal/domain"
	pb "github.com/neo532/gofr_layout/proto/api/user/v1"
	"github.com/neo532/gokit/logger"
	"google.golang.org/protobuf/types/known/emptypb"
)

// UserApi implements pb.UserApi.
type UserApi struct {
	dn  *domain.UserDomain
	log logger.Logger
}

func NewUserApi(
	dn *domain.UserDomain,
	log logger.Logger,
) *UserApi {
	return &UserApi{
		dn:  dn,
		log: log,
	}
}

func (s *UserApi) Post(c context.Context, req *pb.User) (r *emptypb.Empty, err error) {
	err = s.dn.Create(c, req)
	return
}

func (s *UserApi) GetById(c context.Context, req *pb.GetByIdRequest) (*pb.User, error) {
	return s.dn.GetById(c, req.Id)
}
