package domain

import (
	"context"

	// "github.com/neo532/gofr/tool"

	"github.com/neo532/gofr_layout/internal/repo"
	pb "github.com/neo532/gofr_layout/proto/api/user/v1"
)

type UserDomain struct {
	tx   repo.TransactionUserRepo
	user repo.UserRepo
}

func NewUserDomain(
	tx repo.TransactionUserRepo,
	user repo.UserRepo,
) *UserDomain {
	return &UserDomain{
		tx:   tx,
		user: user,
	}
}

func (d *UserDomain) Create(c context.Context, req *pb.User) (err error) {

	err = d.tx.Transaction(c, func(ctx context.Context) (err error) {

		// get
		if _, err = d.user.GetById(ctx, req.Id); err != nil {
			return
		}

		// create
		if _, err = d.user.Create(ctx, req); err != nil {
			return
		}

		return
	})

	return
}

func (d *UserDomain) GetById(c context.Context, id int64) (rst *pb.User, err error) {
	return d.user.GetById(c, id)
}
