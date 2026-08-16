package api

import (
	"context"

	"github.com/neo532/gofr_layout/internal/biz"
	pb "github.com/neo532/gofr_layout/proto/api/user/v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

// UserApi implements pb.UserApi.
type UserApi struct {
	bUser *biz.UserBiz
}

func NewUserApi(
	bUser *biz.UserBiz,
) *UserApi {
	return &UserApi{
		bUser: bUser,
	}
}

func (s *UserApi) Post(c context.Context, req *pb.User) (r *emptypb.Empty, err error) {
	err = s.bUser.Create(c, req)
	return
}

func (s *UserApi) GetById(c context.Context, req *pb.GetByIdRequest) (*pb.User, error) {
	return s.bUser.GetById(c, req.Id)
}
