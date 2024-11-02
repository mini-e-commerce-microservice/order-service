package main

import (
	"context"
	wsqlx "github.com/SyaibanAhmadRamadhan/sqlx-wrapper"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"order-service/internal/conf"
	"order-service/internal/infra"
	"order-service/internal/presentations"
	"order-service/internal/repositories/order_items"
	"order-service/internal/repositories/orders"
	"order-service/internal/repositories/outbox_events"
	"order-service/internal/repositories/product_items"
	"order-service/internal/repositories/saga_states"
	"order-service/internal/services/order"
	"os/signal"
	"syscall"
)

var restApi = &cobra.Command{
	Use:   "restApi",
	Short: "use restApi",
	Run: func(cmd *cobra.Command, args []string) {
		otelConf := conf.LoadOtelConf()
		appConf := conf.LoadAppConf()
		jwtConf := conf.LoadJwtConf()

		closeFnOtel := infra.NewOtel(otelConf, appConf.TracerName)
		pgdb, pgdbCloseFn := infra.NewPostgresql(appConf.DatabaseDsn)
		rdbms := wsqlx.NewRdbms(pgdb)

		productItemRepository := product_items.New(rdbms)
		orderRepository := orders.New(rdbms)
		orderItemRepository := order_items.New(rdbms)
		outboxEventRepository := outbox_events.New(rdbms)
		sagaStateRepository := saga_states.New(rdbms)

		orderService := order.New(order.Opt{
			ProductItemRepository: productItemRepository,
			DatabaseTransaction:   rdbms,
			OrderRepository:       orderRepository,
			OrderItemRepository:   orderItemRepository,
			OutboxEventRepository: outboxEventRepository,
			SagaStateRepository:   sagaStateRepository,
		})

		server := presentations.New(&presentations.Presenter{
			Port:               int(appConf.AppPort),
			JwtAccessTokenConf: jwtConf.AccessToken,
			OrderService:       orderService,
		})
		ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
		defer stop()

		go func() {
			if err := server.ListenAndServe(); err != nil {
				log.Err(err).Msg("failed start serve")
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
