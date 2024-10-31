package primitive

type OrderStatus string

const (
	OrderStatusPending   OrderStatus = "PENDING"
	OrderStatusRejected  OrderStatus = "REJECTED"
	OrderStatusProcessed OrderStatus = "PROCESSED"
	OrderStatusCanceled  OrderStatus = "CANCELED"
	OrderStatusFailed    OrderStatus = "FAILED"
	OrderStatusShipped   OrderStatus = "SHIPPED"
)
