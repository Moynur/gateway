package models

import "github.com/google/uuid"

type AuthRequest struct {
	Card   Card
	Amount Amount
}

type AuthResponse struct {
	TransactionId   string
	OperationId     string
	Response        Response
	AmountAvailable Amount
}

type CaptureRequest struct {
	TransactionId uuid.UUID
	Amount        Amount
}

type CaptureResponse struct {
	TransactionId   uuid.UUID
	OperationId     uuid.UUID
	AmountCharged   Amount
	AmountAvailable Amount
	Response        Response
}

type AuthVoidRequest struct {
	TransactionId uuid.UUID
}

type AuthVoidResponse struct {
	TransactionId uuid.UUID
	OperationId   uuid.UUID
	Response      Response
	Amount        Amount
}

type RefundRequest struct {
	TransactionId uuid.UUID
	Amount        Amount
}

type RefundResponse struct {
	TransactionId   uuid.UUID
	OperationId     uuid.UUID
	Response        Response
	Amount          Amount
	AmountAvailable Amount
}
