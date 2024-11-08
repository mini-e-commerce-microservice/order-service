package order

import "context"

type Service interface {
	CreateOrder(ctx context.Context, input CreateOrderInput) (output CreateOrderOutput, err error)
	ConsumerPaymentResponse(ctx context.Context) (err error)
}
