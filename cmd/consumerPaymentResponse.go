package main

import (
	"context"
	"fmt"
	ekafka "github.com/SyaibanAhmadRamadhan/event-bus/kafka"
	wsqlx "github.com/SyaibanAhmadRamadhan/sqlx-wrapper"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"order-service/internal/conf"
	"order-service/internal/infra"
	"order-service/internal/repositories/orders"
	"order-service/internal/repositories/outbox_events"
	"order-service/internal/repositories/saga_states"
	"order-service/internal/services/order"
	"os/signal"
	"syscall"
)

var consumerPaymentResponse = &cobra.Command{
	Use:   "consumerPaymentResponse",
	Short: "consumerPaymentResponse",
	Run: func(cmd *cobra.Command, args []string) {
		kafkaConf := conf.LoadKafkaConf()
		otelConf := conf.LoadOtelConf()
		appConf := conf.LoadAppConf()

		closeFnOtel := infra.NewOtel(otelConf, appConf.TracerName)
		kafkaBroker := ekafka.New(ekafka.WithOtel())
		pgdb, pgdbCloseFn := infra.NewPostgresql(appConf.DatabaseDsn)
		rdbms := wsqlx.NewRdbms(pgdb)

		orderRepository := orders.New(rdbms)
		outboxEventRepository := outbox_events.New(rdbms)
		sagaStateRepository := saga_states.New(rdbms)

		orderService := order.New(order.Opt{
			DatabaseTransaction:   rdbms,
			OrderRepository:       orderRepository,
			OutboxEventRepository: outboxEventRepository,
			SagaStateRepository:   sagaStateRepository,
			KafkaConf:             kafkaConf,
			KafkaBroker:           kafkaBroker,
		})

		ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
		defer stop()

		go func() {
			err := orderService.ConsumerPaymentResponse(ctx)
			if err != nil {
				fmt.Println(err)
				stop()
			}
		}()

		<-ctx.Done()
		log.Info().Msg("Received shutdown signal, shutting down server gracefully...")

		//time.Sleep(40 * time.Second)
		closeFnOtel(context.TODO())
		pgdbCloseFn(context.TODO())
		log.Info().Msg("Shutdown complete. Exiting.")
	},
}
