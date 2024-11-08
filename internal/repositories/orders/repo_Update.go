package orders

import (
	"context"
	"github.com/Masterminds/squirrel"
	"github.com/SyaibanAhmadRamadhan/go-collection"
	wsqlx "github.com/SyaibanAhmadRamadhan/sqlx-wrapper"
	"github.com/guregu/null/v5"
	"order-service/internal/repositories"
)

func (r *repository) Update(ctx context.Context, input UpdateInput) (err error) {
	if input.Tx == nil {
		return collection.Err(repositories.ErrTxIsNil)
	}

	query := r.sq.Update("orders").Where(squirrel.Eq{"id": input.ID})
	if input.Status.Valid {
		query = query.Set("status", input.Status.String)
	}

	_, err = input.Tx.ExecSq(ctx, query)
	if err != nil {
		return collection.Err(err)
	}
	return
}

type UpdateInput struct {
	Tx     wsqlx.WriterCommand
	ID     int64
	Status null.String
}
