package saga_states

import "context"

type Repository interface {
	Create(ctx context.Context, input CreateInput) (err error)
	FindOne(ctx context.Context, input FindOneInput) (output FindOneOutput, err error)
	Update(ctx context.Context, input UpdateInput) (err error)
}
