package data

import (
	"context"

	"github.com/neo532/gofr_layout/internal/data/base"
	"github.com/neo532/gofr_layout/internal/data/model"
	"github.com/neo532/gofr_layout/internal/repo"
	pb "github.com/neo532/gofr_layout/proto/api/user/v1"
	"github.com/neo532/gokit/database/orm"
)

type UserRepo struct {
	db *orm.Orms
}

func NewUserRepo(
	db base.DatabaseDefault,
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
	err = r.db.Write(c).
		WithContext(c).
		Create(data).
		Error

	insID = data.ID
	return
}

func (r *UserRepo) GetById(c context.Context, id int64) (rst *pb.User, err error) {
	d := &model.User{}
	err = r.db.Read(c).
		WithContext(c).
		Select("id", "name").
		Where("id = ?", id).
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
