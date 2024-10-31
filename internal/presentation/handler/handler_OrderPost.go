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
		ShippingAddressID: req.ShippingAddressId,
		CourierCode:       req.CourierCode,
		PaymentMethodCode: req.PaymentMethodCode,
		Items:             createOrderInputItem,
	})
	if err != nil {
		switch {
		case errors.Is(err, order.ErrProductItemsNotTheSameInDatabase):
			h.httpOtel.Err(w, r, http.StatusBadRequest, err, "invalid product order")
		case errors.Is(err, order.ErrOrderQtyGTStockProduct):
			h.httpOtel.Err(w, r, http.StatusBadRequest, err, "stock product is not available")
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
