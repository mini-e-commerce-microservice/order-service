package order

import (
	"github.com/guregu/null/v5"
	"time"
)

type PayloadAggregatePaymentResponse struct {
	PaymentData PayloadAggregatePayment `json:"payment_data"`
	ErrorReason *string                 `json:"error_reason,omitempty"`
}

type PayloadAggregatePayment struct {
	ID                int64                           `db:"id" json:"id"`
	PaymentSourceCode string                          `db:"payment_source_code" json:"payment_source_code"`
	PaymentMethodCode string                          `db:"payment_method_code" json:"payment_method_code"`
	Status            string                          `db:"status" json:"status"`
	FraudStatus       string                          `db:"fraud_status" json:"fraud_status"`
	ExpiredAt         time.Time                       `db:"expired_at" json:"expired_at"`
	TransactionAt     null.Time                       `db:"transaction_at" json:"transaction_at"`
	OrderID           int64                           `db:"order_id" json:"order_id"`
	StatusMessage     string                          `db:"status_message" json:"status_message"`
	PaymentType       string                          `db:"payment_type" json:"payment_type"`
	GrossAmount       float64                         `db:"gross_amount" json:"gross_amount"`
	SignatureKey      string                          `db:"signature_key" json:"signature_key"`
	Actions           []PayloadAggregatePaymentAction `db:"actions" json:"actions"`
	CreatedAt         time.Time                       `db:"created_at" json:"created_at"`
	UpdatedAt         time.Time                       `db:"updated_at" json:"updated_at"`
}

type PayloadAggregatePaymentAction struct {
	Name   string `json:"name"`
	URL    string `json:"url"`
	Method string `json:"method"`
}
