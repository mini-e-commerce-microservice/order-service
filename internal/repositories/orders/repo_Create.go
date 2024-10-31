package orders

import (
	"context"
	"github.com/SyaibanAhmadRamadhan/go-collection"
	wsqlx "github.com/SyaibanAhmadRamadhan/sqlx-wrapper"
	"order-service/internal/models"
	"order-service/internal/repositories"
)

func (r *repository) Create(ctx context.Context, input CreateInput) (output CreateOutput, err error) {
	if input.Tx == nil {
		return output, collection.Err(repositories.ErrTxIsNil)
	}

	columns, values := collection.GetTagsWithValues(input.Data, "db", "deleted_at", "id")
	query := r.sq.Insert("orders").Columns(columns...).Values(values...).Suffix("RETURNING id")

	err = input.Tx.QueryRowSq(ctx, query, wsqlx.QueryRowScanTypeDefault, &output.ID)
	if err != nil {
		return output, collection.Err(err)
	}
	return
}

type CreateInput struct {
	Tx   wsqlx.ReadQuery
	Data models.Order
}

type CreateOutput struct {
	ID int64
}
