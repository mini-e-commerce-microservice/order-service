package order

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	ekafka "github.com/SyaibanAhmadRamadhan/event-bus/kafka"
	"github.com/SyaibanAhmadRamadhan/go-collection"
	"github.com/SyaibanAhmadRamadhan/go-collection/generic"
	wsqlx "github.com/SyaibanAhmadRamadhan/sqlx-wrapper"
	"github.com/guregu/null/v5"
	"github.com/segmentio/kafka-go"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"order-service/internal/models"
	"order-service/internal/repositories"
	"order-service/internal/repositories/orders"
	"order-service/internal/repositories/outbox_events"
	"order-service/internal/repositories/saga_states"
	"order-service/internal/util"
	"order-service/internal/util/primitive"
)

func (s *service) ConsumerPaymentResponse(ctx context.Context) (err error) {
	output, err := s.kafkaBroker.Subscribe(ctx, ekafka.SubInput{
		Config: kafka.ReaderConfig{
			Brokers: []string{s.kafkaConf.Host},
			GroupID: s.kafkaConf.Topic.OrderSaga.AggregatePaymentResponse.ConsumerGroup.Ordersvc,
			Topic:   s.kafkaConf.Topic.OrderSaga.AggregatePaymentResponse.Name,
		},
	})
	if err != nil {
		return collection.Err(err)
	}

	for {
		data := PayloadAggregatePaymentResponse{}
		msg, err := output.Reader.FetchMessage(ctx, &data)
		if err != nil {
			return collection.Err(err)
		}

		carrier := ekafka.NewMsgCarrier(&msg)
		ctxConsumer := s.propagators.Extract(context.Background(), carrier)

		ctxConsumer, span := otel.Tracer("").Start(ctxConsumer, "process order saga, payment response.",
			trace.WithAttributes(
				attribute.String("payment.status_message", data.PaymentData.StatusMessage),
				attribute.Int64("payment.order_id", data.PaymentData.OrderID),
				attribute.String("payment.status", data.PaymentData.Status),
			))
		if data.ErrorReason != nil {
			span.SetAttributes(attribute.String("payment.error_reason", *data.ErrorReason))
		}

		outboxEventData, err := s.outboxEventRepository.FindOne(ctxConsumer, outbox_events.FindOneInput{
			ID: data.PaymentData.OrderID,
		})
		if err != nil {
			if !errors.Is(err, repositories.ErrNoRecordFound) {
				return collection.Err(err)
			}

			util.SpanRecordErrorWIthEnd(span, err, "outbox event not found")
			err = output.Reader.CommitMessages(ctxConsumer, msg)
			if err != nil {
				util.SpanRecordErrorWIthEnd(span, err, "failed commit message")
				return collection.Err(err)
			}
			continue
		}

		err = s.databaseTransaction.DoTxContext(ctxConsumer, &sql.TxOptions{Isolation: sql.LevelReadCommitted},
			func(ctx context.Context, tx wsqlx.Rdbms) (err error) {
				paymentStatus := primitive.SagaStateStatusSuccess
				if data.ErrorReason != nil || data.PaymentData.Status == "FAILED" {
					paymentStatus = primitive.SagaStateStatusFailed
				} else if data.PaymentData.Status == "PENDING" {
					paymentStatus = primitive.SagaStateStatusOnProcess
				}

				err = s.sagaStateRepository.Update(ctx, saga_states.UpdateInput{
					Tx: tx,
					ID: data.PaymentData.OrderID,
					Data: models.SagaStateCreateOrderProductStep{
						Initiated: string(primitive.SagaStateStatusSuccess),
						Payment:   string(paymentStatus),
						Shipping: generic.Ternary(
							data.PaymentData.Status == "SETTLED",
							string(primitive.SagaStateStatusOnProcess),
							"",
						),
					},
				})
				if err != nil {
					return collection.Err(err)
				}

				if data.ErrorReason != nil {
					err = s.orderRepository.Update(ctx, orders.UpdateInput{
						Tx:     tx,
						ID:     data.PaymentData.OrderID,
						Status: null.StringFrom(string(primitive.OrderStatusFailed)),
					})
					if err != nil {
						return collection.Err(err)
					}

					return
				}

				if data.PaymentData.Status == "SETTLED" {
					err = s.outboxEventRepository.Create(ctx, outbox_events.CreateInput{
						Tx: tx,
						Data: models.OutboxEvent{
							AggregateType: string(primitive.AggregateTypeOutboxEventShipped),
							AggregateID:   fmt.Sprintf("%d", data.PaymentData.OrderID),
							Type:          "created-order",
							Payload:       outboxEventData.Data.Payload,
							TraceParent:   null.StringFrom(carrier.Get("traceparent")).Ptr(),
						},
					})
					if err != nil {
						return collection.Err(err)
					}
				}

				return
			},
		)
		if err != nil {
			util.SpanRecordErrorWIthEnd(span, err, "failed transaction database")
			return collection.Err(err)
		}

		err = output.Reader.CommitMessages(ctx, msg)
		if err != nil {
			util.SpanRecordErrorWIthEnd(span, err, "failed commit message")
			return collection.Err(err)
		}
		span.End()
	}

}
