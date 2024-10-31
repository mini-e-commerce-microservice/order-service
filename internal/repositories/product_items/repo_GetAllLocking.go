package product_items

import (
	"context"
	"github.com/Masterminds/squirrel"
	"github.com/SyaibanAhmadRamadhan/go-collection"
	wsqlx "github.com/SyaibanAhmadRamadhan/sqlx-wrapper"
	"github.com/jmoiron/sqlx"
	"order-service/internal/models"
	"order-service/internal/repositories"
)

func (r *repository) GetAllLocking(ctx context.Context, input GetAllLockingInput) (output GetAllLockingOutput, err error) {
	if input.Tx == nil {
		return output, collection.Err(repositories.ErrTxIsNil)
	}

	query := r.sq.Select("*").From("product_items").Suffix("FOR UPDATE")
	if input.IDs != nil {
		query = query.Where(squirrel.Eq{"product_items.id": input.IDs})
	}
	rdbms := input.Tx

	output = GetAllLockingOutput{
		Data: make([]models.ProductItem, 0),
	}

	err = rdbms.QuerySq(ctx, query, func(rows *sqlx.Rows) (err error) {
		for rows.Next() {
			var item models.ProductItem
			err = rows.StructScan(&item)
			if err != nil {
				return collection.Err(err)
			}

			output.Data = append(output.Data, item)
		}
		return
	})
	
	if err != nil {
		return output, err
	}
	return
}

type GetAllLockingInput struct {
	Tx  wsqlx.ReadQuery
	IDs []int64
}

type GetAllLockingOutput struct {
	Data []models.ProductItem
}
