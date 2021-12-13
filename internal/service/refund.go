package service

import (
	"fmt"
	"log"

	"github.com/moynur/gateway/internal/models"
	"github.com/moynur/gateway/internal/state"
	"github.com/moynur/gateway/internal/store"
)

func (s *service) Refund(refund models.RefundRequest) (models.RefundResponse, error) {
	log.Println("refund request", refund)
	amountToRefund := refund.Amount.MajorUnits
	storeRequest := store.Transaction{
		TransactionId: refund.TransactionId,
		Amount:        amountToRefund,
		Currency:      refund.Amount.Currency,
	}

	latestTransaction, err := s.store.GetLatestByTransactionId(storeRequest.TransactionId)
	if err != nil {
		return models.RefundResponse{}, models.ErrTransactionNotFound
	}

	if isSpecialPan(latestTransaction.Pan, Refund) {
		err := handleErrorCondition(latestTransaction.Pan)
		return models.RefundResponse{}, err
	}

	if !state.TransitionAllowed(latestTransaction.OperationType, state.Refund) {
		return models.RefundResponse{}, models.ErrTransitionNowAllowed
	}

	if !refundIsPossible(latestTransaction, storeRequest) || !refund.IsValid() {
		return models.RefundResponse{}, models.ErrInvalidAmount

	}
	newAvailableAmount := amountToRefund + latestTransaction.AmountAvailable

	operationId, err := s.generator.GenerateUUID()
	if err != nil {
		return models.RefundResponse{}, err
	}
	storeRequest = store.Transaction{
		TransactionId:   refund.TransactionId,
		OperationId:     operationId,
		Amount:          latestTransaction.Amount,
		AmountAvailable: newAvailableAmount,
		Currency:        refund.Amount.Currency,
		OperationType:   state.Refund,
		Pan:             latestTransaction.Pan,
	}
	log.Println("calling store", storeRequest)
	err = s.store.Create(&storeRequest)
	if err != nil {
		return models.RefundResponse{}, fmt.Errorf("couldn't capture %e", err)
	}
	resp := models.RefundResponse{
		TransactionId: refund.TransactionId,
		OperationId:   operationId,
		AmountAvailable: models.Amount{
			MajorUnits: newAvailableAmount,
			Currency:   refund.Amount.Currency,
		},
		Amount: models.Amount{
			MajorUnits: amountToRefund,
			Currency:   refund.Amount.Currency,
		},
		Response: models.Response{Code: models.Approved},
	}
	log.Println("response", resp)
	return resp, nil
}

func refundIsPossible(currentTransaction store.Transaction, request store.Transaction) bool {
	return (currentTransaction.Amount-currentTransaction.AmountAvailable) >= request.Amount &&
		currentTransaction.Currency == request.Currency
}
