package data

import (
	"context"

	"github.com/neo532/gofr_layout/internal/connect"
	"github.com/neo532/gofr_layout/internal/data/model"
	"github.com/neo532/gofr_layout/internal/repo"
	pb "github.com/neo532/gofr_layout/proto/api/user/v1"
	"github.com/neo532/gokit/database/orm"
)

type UserRepo struct {
	db *orm.Orms
}

func NewUserRepo(
	db connect.DatabaseUser,
) repo.UserRepo {
	return &UserRepo{
		db: db,
	}
}

func (r *UserRepo) Create(c context.Context, d *pb.User) (insID int64, err error) {

	data := &model.User{
		ID:   d.Id,
		Name: d.Name,
	}
	err = r.db.Master(c).
		Create(data).
		Error

	insID = data.ID
	return
}

func (r *UserRepo) Update(c context.Context, ID int64, d *pb.User) (err error) {

	err = r.db.Master(c).
		Where("id=?", ID).
		Updates(&model.User{
			Name: d.Name,
		}).
		Error

	return
}

func (r *UserRepo) GetById(c context.Context, ID int64) (rst *pb.User, err error) {
	d := &model.User{}
	err = r.db.Slave(c).
		Select("id", "name").
		Where("id = ?", ID).
		Take(d).
		Error
	if err != nil {
		return
	}

	rst = &pb.User{
		Id:   d.ID,
		Name: d.Name,
	}
	return
}
