package outbox_events

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

	payloadMarshal, err := json.Marshal(input.Data.Payload)
	if err != nil {
		return collection.Err(err)
	}

	query := r.sq.Insert("outbox_events").Columns("aggregatetype", "aggregateid", "type", "payload", "trace_parent").
		Values(input.Data.AggregateType, input.Data.AggregateID, input.Data.Type, string(payloadMarshal), input.Data.TraceParent)

	_, err = input.Tx.ExecSq(ctx, query)
	if err != nil {
		return collection.Err(err)
	}

	return
}

type CreateInput struct {
	Tx   wsqlx.WriterCommand
	Data models.OutboxEvent
}
