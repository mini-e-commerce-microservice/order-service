package models

type OutboxEvent struct {
	ID            string  `db:"id"`
	AggregateType string  `db:"aggregatetype"`
	AggregateID   string  `db:"aggregateid"`
	Type          string  `db:"type"`
	Payload       any     `db:"payload"`
	TraceParent   *string `db:"trace_parent"`
}

type CreateOrderProduct struct {
	OrderID           int64                      `json:"order_id"`
	UserID            int64                      `json:"user_id"`
	Courier           CreateOrderProductCourier  `json:"courier"`
	Origin            CreateOrderProductLocation `json:"origin"`
	Destination       CreateOrderProductLocation `json:"destination"`
	TotalAmount       float64                    `json:"total_amount"`
	PaymentMethodCode string                     `json:"payment_method_code"`
	Items             []OrderItem                `json:"items"`
}

type CreateOrderProductLocation struct {
	LocationID string  `json:"location_id"`
	Latitude   float64 `json:"latitude"`
	Longitude  float64 `json:"longitude"`
	PostalCode int32   `json:"postal_code"`
	Address    string  `json:"address"`
}

type CreateOrderProductCourier struct {
	CourierCode        string `json:"courier_code"`
	CourierServiceCode string `json:"courier_service_code"`
	CourierCompany     string `json:"courier_company"`
	CourierType        string `json:"courier_type"`
	DeliveryType       string `json:"delivery_type"`
	DeliveryDate       string `json:"delivery_date"`
}
