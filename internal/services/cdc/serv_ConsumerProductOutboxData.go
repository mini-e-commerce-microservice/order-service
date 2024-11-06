package cdc

import (
	"context"
	"database/sql"
	"encoding/json"
	"github.com/SyaibanAhmadRamadhan/event-bus/debezium"
	ekafka "github.com/SyaibanAhmadRamadhan/event-bus/kafka"
	"github.com/SyaibanAhmadRamadhan/go-collection"
	wsqlx "github.com/SyaibanAhmadRamadhan/sqlx-wrapper"
	"github.com/segmentio/kafka-go"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"go.opentelemetry.io/otel/trace"
	"order-service/internal/models"
	"order-service/internal/repositories/product_items"
)

func (c *cdc) ConsumerProductOutboxData(ctx context.Context) (err error) {
	output, err := c.kafkaBroker.Subscribe(ctx, ekafka.SubInput{
		Config: kafka.ReaderConfig{
			Brokers: []string{c.kafkaConf.Host},
			GroupID: c.kafkaConf.Topic.SellersvcPublicOutbox.ConsumerGroup.Ordersvc,
			Topic:   c.kafkaConf.Topic.SellersvcPublicOutbox.Name,
		},
	})
	if err != nil {
		return collection.Err(err)
	}

	for {
		data := DebeziumPayload[ProductItemData]{}
		msg, err := output.Reader.FetchMessage(ctx, &data)
		if err != nil {
			return collection.Err(err)
		}

		carrier := ekafka.NewMsgCarrier(&msg)
		ctxConsumer := c.propagators.Extract(context.Background(), carrier)

		ctxConsumer, span := otel.Tracer("").Start(ctxConsumer, string(data.Payload.Op)+" process cdc product item data from user service.",
			trace.WithAttributes(
				attribute.String("cdc.debezium.payload.op", string(data.Payload.Op)),
				attribute.Int64("cdc.debezium.payload.data.id", data.Payload.ID),
			))
		productItem := models.ProductItem{
			TraceParent: data.Payload.TraceParent,
		}
		err = json.Unmarshal([]byte(data.Payload.Payload), &productItem)
		if err != nil {
			return collection.Err(err)
		}

		switch data.Payload.Op {
		case debezium.Create, debezium.Update:
			err = c.dbTransaction.DoTxContext(ctxConsumer, &sql.TxOptions{Isolation: sql.LevelReadCommitted, ReadOnly: false},
				func(ctx context.Context, tx wsqlx.Rdbms) (err error) {
					err = c.productItemRepository.UpSert(ctx, product_items.UpSertInput{
						Tx:   tx,
						Data: productItem,
					})
					if err != nil {
						span.RecordError(collection.Err(err))
						span.SetStatus(codes.Error, err.Error())
						span.SetAttributes(semconv.ErrorTypeKey.String("failed create user"))
						return collection.Err(err)
					}

					err = output.Reader.CommitMessages(ctx, msg)
					if err != nil {
						span.RecordError(collection.Err(err))
						span.SetStatus(codes.Error, err.Error())
						span.SetAttributes(semconv.ErrorTypeKey.String("failed commit message"))
						return collection.Err(err)
					}

					span.SetStatus(codes.Ok, "cdc successfully")
					return nil
				})
			if err != nil {
				span.End()
				return collection.Err(err)
			}
		default:
			err = output.Reader.CommitMessages(ctx, msg)
			if err != nil {
				span.RecordError(collection.Err(err))
				span.SetStatus(codes.Error, err.Error())
				span.SetAttributes(semconv.ErrorTypeKey.String("failed commit message"))
				span.End()
				return err
			}
			span.SetStatus(codes.Error, "unsupported debezium operation type")
			span.SetAttributes(semconv.ErrorTypeKey.String("unsupported debezium operation type"))
		}
		span.End()
	}
}
