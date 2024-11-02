package product_items

import "context"

type Repository interface {
	UpSert(ctx context.Context, input UpSertInput) (err error)
	Count(ctx context.Context, input CountInput) (total int, err error)
	GetAllLocking(ctx context.Context, input GetAllLockingInput) (output GetAllLockingOutput, err error)
	UpdateStock(ctx context.Context, input UpdateStockInput) (err error)
}
