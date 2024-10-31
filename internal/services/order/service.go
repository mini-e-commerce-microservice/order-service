package order

import (
	wsqlx "github.com/SyaibanAhmadRamadhan/sqlx-wrapper"
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
}

type Opt struct {
	ProductItemRepository product_items.Repository
	DatabaseTransaction   wsqlx.Tx
	OrderRepository       orders.Repository
	OrderItemRepository   order_items.Repository
	OutboxEventRepository outbox_events.Repository
	SagaStateRepository   saga_states.Repository
}

func New(opt Opt) *service {
	return &service{
		productItemRepository: opt.ProductItemRepository,
		databaseTransaction:   opt.DatabaseTransaction,
		orderRepository:       opt.OrderRepository,
		orderItemRepository:   opt.OrderItemRepository,
		outboxEventRepository: opt.OutboxEventRepository,
		sagaStateRepository:   opt.SagaStateRepository,
	}
}
