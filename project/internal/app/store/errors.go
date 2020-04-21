package store

import "errors"

// запись своих ошибок
var (
	ErrRecordNotFound   = errors.New("record not found")
	NumberSenderAndDest = errors.New("departure and destination addresses cannot be the same")
	GreaterAmount       = errors.New("the amount cannot be greater than the account balance")
)
