package saga_states

import (
	"context"
	"encoding/json"
	"github.com/SyaibanAhmadRamadhan/go-collection"
	wsqlx "github.com/SyaibanAhmadRamadhan/sqlx-wrapper"
	"order-service/internal/models"
	"order-service/internal/repositories"
)

func (r *repository) Create(ctx context.Context, input CreateInput) (err error) {
	if input.Tx == nil {
		return collection.Err(repositories.ErrTxIsNil)
	}

	payloaddMarshal, err := json.Marshal(input.Data.Payload)
	if err != nil {
		return collection.Err(err)
	}

	stepMarshal, err := json.Marshal(input.Data.Step)
	if err != nil {
		return collection.Err(err)
	}

	query := r.sq.Insert("saga_states").Columns("id", "payload", "status", "step", "type", "version").
		Values(input.Data.ID, string(payloaddMarshal), input.Data.Status, string(stepMarshal), input.Data.Type, input.Data.Version)

	_, err = input.Tx.ExecSq(ctx, query)
	if err != nil {
		return collection.Err(err)
	}

	return
}

type CreateInput struct {
	Tx   wsqlx.WriterCommand
	Data models.SagaState
}
