package product_items

import (
	"context"
	"github.com/Masterminds/squirrel"
	"github.com/SyaibanAhmadRamadhan/go-collection"
	wsqlx "github.com/SyaibanAhmadRamadhan/sqlx-wrapper"
	"order-service/internal/repositories"
)

func (r *repository) UpdateStock(ctx context.Context, input UpdateStockInput) (err error) {
	if input.Tx == nil {
		return collection.Err(repositories.ErrTxIsNil)
	}

	query := r.sq.Update("product_items").Where(squirrel.Eq{"id": input.ID}).Set("stock", input.Stock)

	_, err = input.Tx.ExecSq(ctx, query)
	if err != nil {
		return collection.Err(err)
	}

	return
}

type UpdateStockInput struct {
	Tx    wsqlx.WriterCommand
	ID    int64
	Stock int32
}

type UpdateStockInputItem struct {
}
