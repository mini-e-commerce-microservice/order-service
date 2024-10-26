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
	"order-service/internal/repositories/product_item"
	"order-service/internal/services/cdc"
	"os/signal"
	"syscall"
)

var consumerProductOutbox = &cobra.Command{
	Use:   "consumerProductOutbox",
	Short: "consumerProductOutbox",
	Run: func(cmd *cobra.Command, args []string) {
		kafkaConf := conf.LoadKafkaConf()
		otelConf := conf.LoadOtelConf()
		appConf := conf.LoadAppConf()

		closeFnOtel := infra.NewOtel(otelConf, appConf.TracerName)
		kafkaBroker := ekafka.New(ekafka.WithOtel())
		pgdb, pgdbCloseFn := infra.NewPostgresql(appConf.DatabaseDsn)
		rdbms := wsqlx.NewRdbms(pgdb)

		productItemRepository := product_item.New(rdbms)
		cdcService := cdc.New(kafkaBroker, kafkaConf, rdbms, productItemRepository)

		ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
		defer stop()

		go func() {
			err := cdcService.ConsumerProductOutboxData(ctx)
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
