package primitive

type AggregateTypeOutboxEvent string

const (
	AggregateTypeOutboxEventShipped AggregateTypeOutboxEvent = "shipped"
	AggregateTypeOutboxEventPayment AggregateTypeOutboxEvent = "payment"
)
