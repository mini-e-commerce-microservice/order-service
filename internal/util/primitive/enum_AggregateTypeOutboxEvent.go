package primitive

type AggregateTypeOutboxEvent string

const (
	AggregateTypeOutboxEventCourierRate AggregateTypeOutboxEvent = "courier-rate"
	AggregateTypeOutboxEventShipped     AggregateTypeOutboxEvent = "shipped"
	AggregateTypeOutboxEventPayment     AggregateTypeOutboxEvent = "payment"
)
