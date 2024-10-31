package repositories

import "errors"

var ErrNoRecordFound = errors.New("record not found")
var ErrTxIsNil = errors.New("tx is nil")
