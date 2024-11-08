package saga_states

import (
	"context"
	"database/sql"
	"errors"
	"github.com/Masterminds/squirrel"
	"github.com/SyaibanAhmadRamadhan/go-collection"
	wsqlx "github.com/SyaibanAhmadRamadhan/sqlx-wrapper"
	"order-service/internal/models"
	"order-service/internal/repositories"
)

func (r *repository) FindOne(ctx context.Context, input FindOneInput) (output FindOneOutput, err error) {
	query := r.sq.Select("*").From("saga_states").Where(squirrel.Eq{"id": input.ID})

	err = r.rdbms.QueryRowSq(ctx, query, wsqlx.QueryRowScanTypeStruct, &output.Data)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = repositories.ErrNoRecordFound
		}

		return output, collection.Err(err)
	}

	return
}

type FindOneInput struct {
	ID int64
}

type FindOneOutput struct {
	Data models.SagaState
}
