package order

import "context"

type Service interface {
	CreateOrder(ctx context.Context, input CreateOrderInput) (output CreateOrderOutput, err error)
}
