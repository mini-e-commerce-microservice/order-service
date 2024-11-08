package order

import (
	ekafka "github.com/SyaibanAhmadRamadhan/event-bus/kafka"
	wsqlx "github.com/SyaibanAhmadRamadhan/sqlx-wrapper"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"order-service/generated/proto/secret_proto"
	"order-service/internal/repositories/order_items"
	"order-service/internal/repositories/orders"
	"order-service/internal/repositories/outbox_events"
	"order-service/internal/repositories/product_items"
	"order-service/internal/repositories/saga_states"
)

type service struct {
	productItemRepository product_items.Repository
	databaseTransaction   wsqlx.Tx
	orderRepository       orders.Repository
	orderItemRepository   order_items.Repository
	outboxEventRepository outbox_events.Repository
	sagaStateRepository   saga_states.Repository
	hmacSha256Key         *secret_proto.HmacSha256Key

	propagators propagation.TextMapPropagator
	kafkaConf   *secret_proto.Kafka
	kafkaBroker ekafka.KafkaPubSub
}

type Opt struct {
	ProductItemRepository product_items.Repository
	DatabaseTransaction   wsqlx.Tx
	OrderRepository       orders.Repository
	OrderItemRepository   order_items.Repository
	OutboxEventRepository outbox_events.Repository
	SagaStateRepository   saga_states.Repository
	HmacSha256Key         *secret_proto.HmacSha256Key

	KafkaConf   *secret_proto.Kafka
	KafkaBroker ekafka.KafkaPubSub
}

func New(opt Opt) *service {
	return &service{
		productItemRepository: opt.ProductItemRepository,
		databaseTransaction:   opt.DatabaseTransaction,
		orderRepository:       opt.OrderRepository,
		orderItemRepository:   opt.OrderItemRepository,
		outboxEventRepository: opt.OutboxEventRepository,
		sagaStateRepository:   opt.SagaStateRepository,
		hmacSha256Key:         opt.HmacSha256Key,

		propagators: otel.GetTextMapPropagator(),
		kafkaConf:   opt.KafkaConf,
		kafkaBroker: opt.KafkaBroker,
	}
}
