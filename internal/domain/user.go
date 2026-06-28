package domain

import (
	"context"

	"github.com/neo532/gofr_layout/internal/data/connect"
	"github.com/neo532/gofr_layout/internal/repo"
	pb "github.com/neo532/gofr_layout/proto/api/user/v1"
	"github.com/neo532/gokit/errorx"
)

type UserDomain struct {
	txUser  connect.TransactionUser
	user    repo.UserRepo
	pdcUser connect.ProducerUser
}

func NewUserDomain(
	txUser connect.TransactionUser,
	user repo.UserRepo,
	pdcUser connect.ProducerUser,
) *UserDomain {
	return &UserDomain{
		txUser:  txUser,
		user:    user,
		pdcUser: pdcUser,
	}
}

func (d *UserDomain) Create(c context.Context, req *pb.User) (err error) {

	if err = d.txUser(c, func(c context.Context) (err error) {

		// get
		var data *pb.User
		if data, err = d.user.GetById(c, req.Id); err != nil {
			err = errorx.Wrap(err)
			return
		}

		if data.Id > 0 {
			if err = d.user.Update(c, data.Id, req); err != nil {
				err = errorx.Wrap(err)
			}
			return
		}

		// create
		if _, err = d.user.Create(c, req); err != nil {
			err = errorx.Wrap(err)
			return
		}

		return
	}); err != nil {
		err = errorx.Wrap(err)
		return
	}

	err = d.pdcUser.Send(c, req)
	return
}

func (d *UserDomain) GetById(c context.Context, id int64) (rst *pb.User, err error) {
	return d.user.GetById(c, id)
}
