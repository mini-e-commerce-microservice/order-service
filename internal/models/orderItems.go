package models

type OrderItem struct {
	ID            int64   `db:"id"`
	OrderID       int64   `db:"order_id"`
	ProductItemID int64   `db:"product_item_id"`
	Qty           int32   `db:"qty"`
	UnitPrice     float64 `db:"unit_price"`
	TotalPrice    float64 `db:"total_price"`
	Discount      float64 `db:"discount"`
}
