package order_items

import (
	"context"
	"github.com/SyaibanAhmadRamadhan/go-collection"
	wsqlx "github.com/SyaibanAhmadRamadhan/sqlx-wrapper"
	"order-service/internal/models"
	"order-service/internal/repositories"
)

func (r *repository) Creates(ctx context.Context, input CreatesInput) (err error) {
	if input.Tx == nil {
		return collection.Err(repositories.ErrTxIsNil)
	}

	columns := []string{
		"order_id", "product_item_id", "qty", "unit_price", "total_price", "discount",
		"name", "description", "weight", "package_length", "package_width", "package_height", "dimensional_weight",
	}
	query := r.sq.Insert("order_items").Columns(columns...)
	for _, item := range input.Items {
		query = query.Values(
			input.OrderID, item.ProductItemID, item.Qty, item.UnitPrice, item.TotalPrice, item.Discount,
			item.Name, item.Description, item.Weight, item.PackageLength, item.PackageWidth, item.PackageHeight,
			item.DimensionalWeight,
		)
	}

	_, err = input.Tx.ExecSq(ctx, query)
	if err != nil {
		return collection.Err(err)
	}
	return
}

type CreatesInput struct {
	OrderID int64
	Tx      wsqlx.WriterCommand
	Items   []models.OrderItem
}
