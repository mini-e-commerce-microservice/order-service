package models

import (
	"time"
)

type Order struct {
	ID                int64      `db:"id"`
	UserID            int64      `db:"user_id"`
	Status            string     `db:"status"`
	TotalAmount       float64    `db:"total_amount"`
	PaymentStatus     string     `db:"payment_status"`
	PaymentMethodCode string     `db:"payment_method_code"`
	Tax               float64    `db:"tax"`
	ShippingCost      float64    `db:"shipping_cost"`
	Discount          float64    `db:"discount"`
	OrderDate         time.Time  `db:"order_date"`
	CreatedAt         time.Time  `db:"created_at"`
	UpdatedAt         time.Time  `db:"updated_at"`
	DeletedAt         *time.Time `db:"deleted_at"`
}
