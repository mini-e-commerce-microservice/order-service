package cdc

import "context"

type Service interface {
	ConsumerProductOutboxData(ctx context.Context) (err error)
}
