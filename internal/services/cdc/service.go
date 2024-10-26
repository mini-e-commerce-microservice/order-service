package cdc

import (
	ekafka "github.com/SyaibanAhmadRamadhan/event-bus/kafka"
	wsqlx "github.com/SyaibanAhmadRamadhan/sqlx-wrapper"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"order-service/generated/proto/secret_proto"
	"order-service/internal/repositories/product_item"
)

type cdc struct {
	kafkaBroker           ekafka.KafkaPubSub
	kafkaConf             *secret_proto.Kafka
	propagators           propagation.TextMapPropagator
	dbTransaction         wsqlx.Tx
	productItemRepository product_item.Repository
}

func New(kafkaBroker ekafka.KafkaPubSub, kafkaConf *secret_proto.Kafka, dbTransaction wsqlx.Tx, productItemRepository product_item.Repository) *cdc {
	return &cdc{
		propagators:           otel.GetTextMapPropagator(),
		kafkaBroker:           kafkaBroker,
		dbTransaction:         dbTransaction,
		kafkaConf:             kafkaConf,
		productItemRepository: productItemRepository,
	}
}
