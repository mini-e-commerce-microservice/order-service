package product_items

import (
	"context"
	"github.com/Masterminds/squirrel"
	"github.com/SyaibanAhmadRamadhan/go-collection"
	wsqlx "github.com/SyaibanAhmadRamadhan/sqlx-wrapper"
)

func (r *repository) Count(ctx context.Context, input CountInput) (total int, err error) {
	query := r.sq.Select("COUNT(*)").From("product_items")

	if input.IDs != nil {
		query = query.Where(squirrel.Eq{"product_items.id": input.IDs})
	}

	err = r.rdbms.QueryRowSq(ctx, query, wsqlx.QueryRowScanTypeDefault, &total)
	if err != nil {
		return total, collection.Err(err)
	}

	return
}

type CountInput struct {
	IDs []int64
}
