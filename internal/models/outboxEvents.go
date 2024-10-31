package models

type OutboxEvent struct {
	ID            int64   `db:"id"`
	Aggregatetype string  `db:"aggregatetype"`
	Aggregateid   string  `db:"aggregateid"`
	Type          string  `db:"type"`
	Payload       any     `db:"payload"`
	TraceParent   *string `db:"trace_parent"`
}

type CreateOrderProduct struct {
	OrderID           int64   `json:"order_id"`
	UserID            int64   `db:"user_id"`
	ShippingAddressID int64   `db:"shipping_address_id"`
	TotalAmount       float64 `db:"total_amount"`
	PaymentMethodCode string  `db:"payment_method_code"`
}
