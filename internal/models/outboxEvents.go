package models

type OutboxEvent struct {
	ID            int64   `db:"id"`
	AggregateType string  `db:"aggregatetype"`
	AggregateID   string  `db:"aggregateid"`
	Type          string  `db:"type"`
	Payload       any     `db:"payload"`
	TraceParent   *string `db:"trace_parent"`
}

type CreateOrderProduct struct {
	OrderID              int64       `json:"order_id"`
	UserID               int64       `json:"user_id"`
	DestinationAddressID int64       `json:"destination_address_id"`
	OriginAddressID      int64       `json:"origin_address_id"`
	CourierCode          string      `json:"courier_code"`
	CourierServiceCode   string      `json:"courier_service_code"`
	TotalAmount          float64     `json:"total_amount"`
	PaymentMethodCode    string      `json:"payment_method_code"`
	Items                []OrderItem `json:"items"`
}
