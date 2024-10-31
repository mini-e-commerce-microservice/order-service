package order

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/SyaibanAhmadRamadhan/go-collection"
	"github.com/SyaibanAhmadRamadhan/go-collection/generic"
	wsqlx "github.com/SyaibanAhmadRamadhan/sqlx-wrapper"
	"golang.org/x/sync/errgroup"
	"order-service/internal/models"
	"order-service/internal/repositories/order_items"
	"order-service/internal/repositories/orders"
	"order-service/internal/repositories/outbox_events"
	"order-service/internal/repositories/product_items"
	"order-service/internal/repositories/saga_states"
	"order-service/internal/util"
	"order-service/internal/util/primitive"
	"time"
)

func (s *service) CreateOrder(ctx context.Context, input CreateOrderInput) (output CreateOrderOutput, err error) {
	productItemIDs := generic.Appends(input.Items, func(t CreateOrderInputItem) int64 {
		return t.ProductItemID
	}, generic.WithUnique(true))

	totalProductItem, err := s.productItemRepository.Count(ctx, product_items.CountInput{
		IDs: productItemIDs,
	})
	if err != nil {
		return output, collection.Err(err)
	}
	if totalProductItem != len(input.Items) {
		return output, collection.Err(ErrProductItemsNotTheSameInDatabase)
	}

	err = s.databaseTransaction.DoTxContext(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted},
		func(ctx context.Context, tx wsqlx.Rdbms) (err error) {
			productItemOutput, err := s.productItemRepository.GetAllLocking(ctx, product_items.GetAllLockingInput{
				Tx:  tx,
				IDs: productItemIDs,
			})
			if err != nil {
				return collection.Err(err)
			}

			orderItemsInputMap := generic.ConvertToMap(input.Items, func(t CreateOrderInputItem) int64 {
				return t.ProductItemID
			})
			totalAmount := float64(0)
			orderItems := make([]models.OrderItem, 0, len(input.Items))

			for _, item := range productItemOutput.Data {
				orderItemInput := orderItemsInputMap[item.ID]
				if orderItemInput.Quantity > item.Stock {
					return collection.Err(ErrOrderQtyGTStockProduct)
				}

				orderItemTotalPrice := item.Price * float64(orderItemInput.Quantity)
				totalAmount += orderItemTotalPrice

				orderItems = append(orderItems, models.OrderItem{
					ProductItemID: item.ID,
					Qty:           orderItemInput.Quantity,
					UnitPrice:     item.Price,
					TotalPrice:    orderItemTotalPrice,
					Discount:      0,
				})
			}

			timeNow := time.Now().UTC()
			orderCreateOutput, err := s.orderRepository.Create(ctx, orders.CreateInput{
				Data: models.Order{
					UserID:            input.UserID,
					ShippingAddressID: input.ShippingAddressID,
					Status:            string(primitive.OrderStatusPending),
					TotalAmount:       totalAmount,
					PaymentStatus:     "",
					PaymentMethodCode: input.PaymentMethodCode,
					Tax:               0,
					ShippingCost:      0,
					Discount:          0,
					OrderDate:         timeNow,
					CreatedAt:         timeNow,
					UpdatedAt:         timeNow,
				},
			})
			if err != nil {
				return collection.Err(err)
			}

			var eg errgroup.Group
			eg.Go(func() (err error) {
				err = s.orderItemRepository.Creates(ctx, order_items.CreatesInput{
					OrderID: orderCreateOutput.ID,
					Tx:      tx,
					Items:   orderItems,
				})
				if err != nil {
					return collection.Err(err)
				}

				return
			})

			createOrderProductPayload := models.CreateOrderProduct{
				OrderID:           orderCreateOutput.ID,
				UserID:            input.UserID,
				ShippingAddressID: input.ShippingAddressID,
				TotalAmount:       totalAmount,
				PaymentMethodCode: input.PaymentMethodCode,
			}
			eg.Go(func() (err error) {
				err = s.outboxEventRepository.Create(ctx, outbox_events.CreateInput{
					Tx: tx,
					Data: models.OutboxEvent{
						Aggregatetype: string(primitive.AggregateTypeOutboxEventCourierRate),
						Aggregateid:   fmt.Sprintf("%d", orderCreateOutput.ID),
						Type:          "created-order",
						Payload:       createOrderProductPayload,
						TraceParent:   util.GetTraceParent(ctx),
					},
				})
				if err != nil {
					return collection.Err(err)
				}
				return
			})

			eg.Go(func() (err error) {
				err = s.sagaStateRepository.Create(ctx, saga_states.CreateInput{
					Tx: tx,
					Data: models.SagaState{
						Payload: createOrderProductPayload,
						Status:  string(primitive.SagaStateStatusOnProcess),
						Step: models.SagaStateCreateOrderProductStep{
							Initiated:         "success",
							ShippingCalculate: "on process",
						},
						Type:    "order placement",
						Version: "1",
					},
				})
				if err != nil {
					return collection.Err(err)
				}

				return
			})

			if err = eg.Wait(); err != nil {
				return collection.Err(err)
			}

			output.OrderID = orderCreateOutput.ID
			return
		},
	)
	if err != nil {
		return output, collection.Err(err)
	}

	return
}

type CreateOrderInput struct {
	UserID            int64
	ShippingAddressID int64
	CourierCode       string
	PaymentMethodCode string
	Items             []CreateOrderInputItem
}

type CreateOrderInputItem struct {
	ProductItemID int64
	Quantity      int32
}

type CreateOrderOutput struct {
	OrderID int64
}
