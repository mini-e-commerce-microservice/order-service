package order

import "errors"

var ErrProductItemsNotTheSameInDatabase = errors.New("total product items not the same in database")
var ErrOrderQtyGTStockProduct = errors.New("order qty greater than stock product")
