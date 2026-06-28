package api

import (
	"context"

	pb "github.com/neo532/gofr_layout/proto/api/user/v1"
	"github.com/neo532/gokit/logger"
	"google.golang.org/protobuf/types/known/emptypb"
)

// User1Api implements pb.User1Api.
type User1Api struct {
	log logger.Logger
}

func NewUser1Api(
	log logger.Logger,
) *User1Api {
	return &User1Api{
		log: log,
	}
}
func (s *User1Api) Post(ctx context.Context, req *pb.User) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

func (s *User1Api) GetById(ctx context.Context, req *pb.GetByIdRequest) (*pb.User, error) {
	return &pb.User{Name: "Hello11"}, nil
}
