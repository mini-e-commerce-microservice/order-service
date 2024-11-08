package order

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	"github.com/SyaibanAhmadRamadhan/go-collection"
	"github.com/SyaibanAhmadRamadhan/go-collection/generic"
	wsqlx "github.com/SyaibanAhmadRamadhan/sqlx-wrapper"
	"github.com/guregu/null/v5"
	"golang.org/x/sync/errgroup"
	"google.golang.org/protobuf/proto"
	"order-service/generated/proto/hmac_sha_256_payload"
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
		IDs:      productItemIDs,
		IsActive: null.BoolFrom(true),
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
			totalAmount := input.Courier.CourierPrice
			orderItems := make([]models.OrderItem, 0, len(input.Items))
			payloadHmacSha256ProductItem := make([]*hmac_sha_256_payload.CourierRateProductItem, 0, len(input.Items))

			for _, item := range productItemOutput.Data {
				orderItemInput := orderItemsInputMap[item.ID]
				if orderItemInput.Quantity > item.Stock {
					return collection.Err(ErrOrderQtyGTStockProduct)
				}

				if orderItemInput.Quantity < item.MinimumPurchase {
					return collection.Err(ErrYourQuantityIsLTMinimumPurchase)
				}

				orderItemTotalPrice := item.Price * float64(orderItemInput.Quantity)
				totalAmount += orderItemTotalPrice

				orderItems = append(orderItems, models.OrderItem{
					Name:              item.Name,
					Description:       item.Description,
					ProductItemID:     item.ID,
					Qty:               orderItemInput.Quantity,
					UnitPrice:         item.Price,
					TotalPrice:        orderItemTotalPrice,
					Discount:          0,
					Weight:            item.Weight,
					PackageLength:     item.PackageLength,
					PackageWidth:      item.PackageWidth,
					PackageHeight:     item.PackageHeight,
					DimensionalWeight: item.DimensionalWeight,
				})

				payloadHmacSha256ProductItem = append(payloadHmacSha256ProductItem, &hmac_sha_256_payload.CourierRateProductItem{
					Length:    int64(item.PackageLength),
					Width:     int64(item.PackageWidth),
					Height:    int64(item.PackageHeight),
					Weight:    int64(item.Weight),
					Quantity:  int64(orderItemInput.Quantity),
					Price:     item.Price,
					ProductId: item.ID,
					Name:      item.Name,
				})
			}

			err = s.validateRequest(input, payloadHmacSha256ProductItem)
			if err != nil {
				return collection.Err(err)
			}

			var eg errgroup.Group
			eg.Go(func() (err error) {
				for _, item := range productItemOutput.Data {
					orderItemInput := orderItemsInputMap[item.ID]
					err = s.productItemRepository.UpdateStock(ctx, product_items.UpdateStockInput{
						Tx:    tx,
						ID:    item.ID,
						Stock: item.Stock - orderItemInput.Quantity,
					})
					if err != nil {
						return collection.Err(err)
					}
				}
				return
			})

			timeNow := time.Now().UTC()
			orderCreateOutput, err := s.orderRepository.Create(ctx, orders.CreateInput{
				Tx: tx,
				Data: models.Order{
					UserID:            input.UserID,
					Status:            string(primitive.OrderStatusProcessed),
					TotalAmount:       totalAmount,
					PaymentStatus:     "",
					PaymentMethodCode: input.PaymentMethodCode,
					Tax:               0,
					ShippingCost:      input.Courier.CourierPrice,
					Discount:          0,
					OrderDate:         timeNow,
					CreatedAt:         timeNow,
					UpdatedAt:         timeNow,
				},
			})
			if err != nil {
				return collection.Err(err)
			}
			createOrderProductPayload := models.CreateOrderProduct{
				OrderID: orderCreateOutput.ID,
				UserID:  input.UserID,
				Courier: models.CreateOrderProductCourier{
					CourierCode:        input.Courier.CourierCode,
					CourierServiceCode: input.Courier.CourierServiceCode,
					CourierCompany:     input.Courier.Company,
					DeliveryType:       "now",
					DeliveryDate:       timeNow.Format(time.DateOnly),
					CourierType:        input.Courier.Type,
				},
				Origin: models.CreateOrderProductLocation{
					LocationID: input.Origin.LocationId,
					Latitude:   input.Origin.Latitude,
					Longitude:  input.Origin.Longitude,
					Address:    input.Origin.Address,
					PostalCode: input.Origin.PostalCode,
				},
				Destination: models.CreateOrderProductLocation{
					LocationID: input.Destination.LocationId,
					Latitude:   input.Destination.Latitude,
					Longitude:  input.Destination.Longitude,
					Address:    input.Destination.Address,
					PostalCode: input.Destination.PostalCode,
				},
				TotalAmount:       totalAmount,
				PaymentMethodCode: input.PaymentMethodCode,
				Items:             orderItems,
			}

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

			eg.Go(func() (err error) {
				err = s.outboxEventRepository.Create(ctx, outbox_events.CreateInput{
					Tx: tx,
					Data: models.OutboxEvent{
						AggregateType: string(primitive.AggregateTypeOutboxEventPayment),
						AggregateID:   fmt.Sprintf("%d", orderCreateOutput.ID),
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
						ID:      orderCreateOutput.ID,
						Payload: createOrderProductPayload,
						Status:  string(primitive.SagaStateStatusOnProcess),
						Step: models.SagaStateCreateOrderProductStep{
							Initiated: string(primitive.SagaStateStatusSuccess),
							Payment:   string(primitive.SagaStateStatusOnProcess),
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

func (s *service) validateRequest(input CreateOrderInput, payloadHmacSha256ProductItem []*hmac_sha_256_payload.CourierRateProductItem) (err error) {
	payloadSha := &hmac_sha_256_payload.CourierRate{
		ProductItem:                  payloadHmacSha256ProductItem,
		AvailableForCashOnDelivery:   input.Courier.AvailableForCashOnDelivery,
		AvailableForProofOfDelivery:  input.Courier.AvailableForProofOfDelivery,
		AvailableForInstantWaybillId: input.Courier.AvailableForInstantWaybillID,
		AvailableForInsurance:        input.Courier.AvailableForInsurance,
		Company:                      input.Courier.Company,
		CourierCode:                  input.Courier.CourierCode,
		CourierServiceCode:           input.Courier.CourierServiceCode,
		Duration:                     input.Courier.Duration,
		ShipmentDurationRange:        input.Courier.ShipmentDurationRange,
		ShipmentDurationUnit:         input.Courier.ShipmentDurationUnit,
		ServiceType:                  input.Courier.ServiceType,
		CourierPrice:                 input.Courier.CourierPrice,
		Type:                         input.Courier.Type,
		Origin: &hmac_sha_256_payload.CourierLocation{
			LocationId: input.Origin.LocationId,
			Latitude:   input.Origin.Latitude,
			Longitude:  input.Origin.Longitude,
			Address:    input.Origin.Address,
			PostalCode: input.Origin.PostalCode,
		},
		Destination: &hmac_sha_256_payload.CourierLocation{
			LocationId: input.Destination.LocationId,
			Latitude:   input.Destination.Latitude,
			Longitude:  input.Destination.Longitude,
			Address:    input.Destination.Address,
			PostalCode: input.Destination.PostalCode,
		},
	}

	payloadShaMarshal, err := proto.Marshal(payloadSha)
	if err != nil {
		return collection.Err(err)
	}

	hash := hmac.New(sha256.New, []byte(s.hmacSha256Key.ShippmentServiceCourierRate))
	hash.Write(payloadShaMarshal)

	h := hex.EncodeToString(hash.Sum(nil))

	if input.Courier.ID != h {
		return collection.Err(ErrInvalidCourier)
	}

	return
}

type CreateOrderInput struct {
	UserID            int64
	PaymentMethodCode string
	Courier           CreateOrderInputCourier
	Origin            CreateOrderInputLocation
	Destination       CreateOrderInputLocation
	Items             []CreateOrderInputItem
}

type CreateOrderInputLocation struct {
	Address    string
	Latitude   float64
	LocationId string
	Longitude  float64
	PostalCode int32
}

type CreateOrderInputCourier struct {
	ID                           string
	AvailableForCashOnDelivery   bool
	AvailableForProofOfDelivery  bool
	AvailableForInstantWaybillID bool
	AvailableForInsurance        bool
	Company                      string
	CourierCode                  string
	CourierServiceCode           string
	Duration                     string
	ShipmentDurationRange        string
	ShipmentDurationUnit         string
	ServiceType                  string
	CourierPrice                 float64
	Type                         string
}

type CreateOrderInputItem struct {
	ProductItemID int64
	Quantity      int32
}

type CreateOrderOutput struct {
	OrderID int64
}
