package order

import "errors"

var ErrProductItemsNotTheSameInDatabase = errors.New("total product items not the same in database")
var ErrOrderQtyGTStockProduct = errors.New("order qty greater than stock product")
var ErrInvalidCourier = errors.New("invalid courier")
var ErrYourQuantityIsLTMinimumPurchase = errors.New("quantity is less than the minimum purchase requirement")
