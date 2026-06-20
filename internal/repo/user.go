package repo

import (
	"context"

	pb "github.com/neo532/gofr_layout/proto/api/user/v1"
)

type TransactionUserRepo interface {
	Transaction(c context.Context, fn func(c context.Context) (err error)) (err error)
}

type UserRepo interface {
	Create(c context.Context, user *pb.User) (insID int64, err error)
	GetById(c context.Context, id int64) (rst *pb.User, err error)
}
