package saga_states

import (
	"context"
	"encoding/json"
	"github.com/Masterminds/squirrel"
	"github.com/SyaibanAhmadRamadhan/go-collection"
	wsqlx "github.com/SyaibanAhmadRamadhan/sqlx-wrapper"
	"order-service/internal/models"
	"order-service/internal/repositories"
)

func (r *repository) Update(ctx context.Context, input UpdateInput) (err error) {
	if input.Tx == nil {
		return collection.Err(repositories.ErrTxIsNil)
	}

	dataMarshal, err := json.Marshal(input.Data)
	if err != nil {
		return collection.Err(err)
	}

	query := r.sq.Update("saga_states").Set("step", dataMarshal).Where(squirrel.Eq{"id": input.ID})
	_, err = input.Tx.ExecSq(ctx, query)
	if err != nil {
		return collection.Err(err)
	}
	return
}

type UpdateInput struct {
	Tx   wsqlx.WriterCommand
	ID   int64
	Data models.SagaStateCreateOrderProductStep
}
