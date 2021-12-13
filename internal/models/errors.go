package models

import "errors"

var (
	ErrInvalidAmount        = errors.New("invalid amount for request")
	ErrFailedLuhn           = errors.New("pan is invalid")
	ErrCantAuth             = errors.New("can't authorize pan is not allowed")
	ErrCantCapture          = errors.New("can't capture pan is not allowed")
	ErrCantRefund           = errors.New("can't refund pan is not allowed")
	ErrTransactionNotFound  = errors.New("no transaction found")
	ErrTransitionNowAllowed = errors.New("invalid request flow")
	ErrGeneric              = errors.New("internal error")
)
