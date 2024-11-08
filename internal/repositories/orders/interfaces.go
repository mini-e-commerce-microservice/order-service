package orders

import "context"

type Repository interface {
	Create(ctx context.Context, input CreateInput) (output CreateOutput, err error)
	Update(ctx context.Context, input UpdateInput) (err error)
}
