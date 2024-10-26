package product_item

import "context"

type Repository interface {
	UpSert(ctx context.Context, input UpSertInput) (err error)
}
