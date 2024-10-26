package product_item

import (
	"github.com/Masterminds/squirrel"
	wsqlx "github.com/SyaibanAhmadRamadhan/sqlx-wrapper"
)

type AggregateType string

const AggregateTypeProduct AggregateType = "aggregate_product"

type repository struct {
	rdbms wsqlx.Rdbms
	sq    squirrel.StatementBuilderType
}

func New(rdbms wsqlx.Rdbms) *repository {
	return &repository{
		rdbms: rdbms,
		sq:    squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}
