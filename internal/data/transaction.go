package data

import (
	"context"

	"github.com/neo532/gokit/database/orm"

	"github.com/neo532/gofr_layout/internal/data/base"
	"github.com/neo532/gofr_layout/internal/repo"
)

type TransactionDefaultRepo struct {
	db *orm.Orms
}

func NewTransactionDefaultRepo(defaultDB base.DatabaseDefault) repo.TransactionDefaultRepo {
	return &TransactionDefaultRepo{
		db: defaultDB,
	}
}

func (r *TransactionDefaultRepo) Transaction(c context.Context, fn func(ctx context.Context) error) (err error) {
	err = r.db.Transaction(c, fn)
	return
}
