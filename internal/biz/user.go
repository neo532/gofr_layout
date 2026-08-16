package biz

import (
	"context"

	"github.com/neo532/gofr_layout/internal/connect"
	"github.com/neo532/gofr_layout/internal/repo"
	pb "github.com/neo532/gofr_layout/proto/api/user/v1"
	"github.com/neo532/gokit/errorx"
)

type UserBiz struct {
	txUser  connect.TransactionUser
	pdcUser connect.ProducerUser
	rUser   repo.UserRepo
}

func NewUserBiz(
	txUser connect.TransactionUser,
	pdcUser connect.ProducerUser,
	rUser repo.UserRepo,
) *UserBiz {
	return &UserBiz{
		txUser:  txUser,
		pdcUser: pdcUser,
		rUser:   rUser,
	}
}

func (d *UserBiz) Create(c context.Context, req *pb.User) (err error) {

	if err = d.txUser(c, func(c context.Context) (err error) {

		// get
		var data *pb.User
		if data, err = d.rUser.GetById(c, req.Id); err != nil {
			err = errorx.Wrap(err)
			return
		}

		if data.Id > 0 {
			if err = d.rUser.Update(c, data.Id, req); err != nil {
				err = errorx.Wrap(err)
			}
			return
		}

		// create
		if _, err = d.rUser.Create(c, req); err != nil {
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

func (d *UserBiz) GetById(c context.Context, id int64) (rst *pb.User, err error) {
	return d.rUser.GetById(c, id)
}
