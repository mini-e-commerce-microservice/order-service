package models

type SagaState struct {
	ID      int    `db:"id"`
	Payload any    `db:"payload"`
	Status  string `db:"status"`
	Step    any    `db:"step"`
	Type    string `db:"type"`
	Version string `db:"version"`
}

type SagaStateCreateOrderProductStep struct {
	Initiated         string `json:"initiated,omitempty"`
	ShippingCalculate string `json:"shipping_calculate,omitempty"`
	InitiatePayment   string `json:"initiate_payment,omitempty"`
	PaymentProcess    string `json:"payment_process,omitempty"`
	OrderStatus       string `json:"order_status,omitempty"`
}
