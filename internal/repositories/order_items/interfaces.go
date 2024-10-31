package order_items

import "context"

type Repository interface {
	Creates(ctx context.Context, input CreatesInput) (err error)
}
