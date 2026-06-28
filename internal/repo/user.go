package repo

import (
	"context"

	pb "github.com/neo532/gofr_layout/proto/api/user/v1"
)

type UserRepo interface {
	Create(c context.Context, user *pb.User) (insID int64, err error)
	Update(c context.Context, ID int64, d *pb.User) (err error)
	GetById(c context.Context, id int64) (rst *pb.User, err error)
}
