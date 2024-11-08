package handler

import (
	"errors"
	"net/http"
	"order-service/generated/api"
	"order-service/internal/services/order"
)

func (h *handler) OrderPost(w http.ResponseWriter, r *http.Request) {
	req := api.V1OrderPostJSONRequestBody{}
	if !h.httpOtel.BindBodyRequest(w, r, &req) {
		return
	}

	userData, ok := h.getUserFromBearerAuth(w, r, false)
	if !ok {
		return
	}

	createOrderInputItem := make([]order.CreateOrderInputItem, 0, len(req.Items))
	for _, item := range req.Items {
		createOrderInputItem = append(createOrderInputItem, order.CreateOrderInputItem{
			ProductItemID: item.ProductItemId,
			Quantity:      item.Qty,
		})
	}

	createOrderOutput, err := h.serv.orderService.CreateOrder(r.Context(), order.CreateOrderInput{
		UserID:            userData.UserId,
		PaymentMethodCode: req.PaymentMethodCode,
		Courier: order.CreateOrderInputCourier{
			ID:                           req.Courier.Id,
			AvailableForCashOnDelivery:   req.Courier.AvailableForCashOnDelivery,
			AvailableForProofOfDelivery:  req.Courier.AvailableForProofOfDelivery,
			AvailableForInstantWaybillID: req.Courier.AvailableForInstantWaybillId,
			AvailableForInsurance:        req.Courier.AvailableForInsurance,
			Company:                      req.Courier.Company,
			CourierCode:                  req.Courier.CourierCode,
			CourierServiceCode:           req.Courier.CourierServiceCode,
			Duration:                     req.Courier.Duration,
			ShipmentDurationRange:        req.Courier.ShipmentDurationRange,
			ShipmentDurationUnit:         req.Courier.ShipmentDurationUnit,
			ServiceType:                  req.Courier.ServiceType,
			CourierPrice:                 req.Courier.Price,
			Type:                         req.Courier.Type,
		},
		Origin: order.CreateOrderInputLocation{
			Address:    req.Origin.Address,
			Latitude:   req.Origin.Latitude,
			LocationId: req.Origin.LocationId,
			Longitude:  req.Origin.Longitude,
			PostalCode: req.Origin.PostalCode,
		},
		Destination: order.CreateOrderInputLocation{
			Address:    req.Destination.Address,
			Latitude:   req.Destination.Latitude,
			LocationId: req.Destination.LocationId,
			Longitude:  req.Destination.Longitude,
			PostalCode: req.Destination.PostalCode,
		},
		Items: createOrderInputItem,
	})
	if err != nil {
		switch {
		case errors.Is(err, order.ErrProductItemsNotTheSameInDatabase):
			h.httpOtel.Err(w, r, http.StatusBadRequest, err, "invalid product order")
		case errors.Is(err, order.ErrOrderQtyGTStockProduct):
			h.httpOtel.Err(w, r, http.StatusBadRequest, err, "stock product is not available")
		case errors.Is(err, order.ErrInvalidCourier):
			h.httpOtel.Err(w, r, http.StatusBadRequest, err, "invalid courier")
		case errors.Is(err, order.ErrYourQuantityIsLTMinimumPurchase):
			h.httpOtel.Err(w, r, http.StatusBadRequest, err, order.ErrYourQuantityIsLTMinimumPurchase.Error())
		default:
			h.httpOtel.Err(w, r, http.StatusInternalServerError, err)
		}
		return
	}

	resp := api.V1OrderResponse200{
		OrderId: createOrderOutput.OrderID,
	}

	h.httpOtel.WriteJson(w, r, http.StatusCreated, resp)
}
