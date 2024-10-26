package product_item

import (
	"context"
	"github.com/SyaibanAhmadRamadhan/go-collection"
	wsqlx "github.com/SyaibanAhmadRamadhan/sqlx-wrapper"
	"order-service/internal/models"
)

func (r *repository) UpSert(ctx context.Context, input UpSertInput) (err error) {
	variant1Marshal, err := input.Data.Variant1.MarshalJSON()
	if err != nil {
		return collection.Err(err)
	}
	variant2Marshal, err := input.Data.Variant1.MarshalJSON()
	if err != nil {
		return collection.Err(err)
	}
	tags, values := collection.GetTagsWithValues(input.Data, "db", "variant_1", "variant_2")
	tags = append(tags, "variant_1", "variant_2")
	values = append(values, variant1Marshal, variant2Marshal)
	query := r.sq.Insert("product_items").
		Columns(tags...).Values(values...).
		Suffix(
			`ON CONFLICT(id) DO UPDATE
					SET variant_1 = EXCLUDED.variant_1,
						variant_2 = EXCLUDED.variant_2,
						sub_category_item_name = EXCLUDED.sub_category_item_name,
						name = EXCLUDED.name,
						description = EXCLUDED.description,
						price = EXCLUDED.price,
						stock = EXCLUDED.stock,
						sku = EXCLUDED.sku,
						weight = EXCLUDED.weight,
						package_length = EXCLUDED.package_length,
						package_width = EXCLUDED.package_width,
						package_height = EXCLUDED.package_height,
						dimensional_weight = EXCLUDED.dimensional_weight,
						is_active = EXCLUDED.is_active,
						product_condition = EXCLUDED.product_condition,
						minimum_purchase = EXCLUDED.minimum_purchase,
						size_guide_image = EXCLUDED.size_guide_image,
						created_at = EXCLUDED.created_at,
						updated_at = EXCLUDED.updated_at,
						deleted_at = EXCLUDED.deleted_at
				`,
		)

	rdbms := input.Tx
	if input.Tx == nil {
		rdbms = r.rdbms
	}

	_, err = rdbms.ExecSq(ctx, query)
	if err != nil {
		return err
	}

	return
}

type UpSertInput struct {
	Tx   wsqlx.WriterCommand
	Data models.ProductItem
}
