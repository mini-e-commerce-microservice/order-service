package primitive

type OrderStatus string

const (
	OrderStatusRejected  OrderStatus = "REJECTED"
	OrderStatusProcessed OrderStatus = "PROCESSED"
	OrderStatusCanceled  OrderStatus = "CANCELED"
	OrderStatusFailed    OrderStatus = "FAILED"
	OrderStatusShipped   OrderStatus = "SHIPPED"
)
