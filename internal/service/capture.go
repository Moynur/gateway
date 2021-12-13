package service

import (
	"fmt"
	"log"

	"github.com/moynur/gateway/internal/models"
	"github.com/moynur/gateway/internal/state"
	"github.com/moynur/gateway/internal/store"
)

func (s *service) Capture(capture models.CaptureRequest) (models.CaptureResponse, error) {
	amountToCapture := capture.Amount.MajorUnits
	storeRequest := store.Transaction{
		TransactionId: capture.TransactionId,
		Amount:        amountToCapture,
		Currency:      capture.Amount.Currency,
	}

	latestTransaction, err := s.store.GetLatestByTransactionId(storeRequest.TransactionId)
	if err != nil {
		return models.CaptureResponse{}, models.ErrTransactionNotFound
	}

	if isSpecialPan(latestTransaction.Pan, Capture) {
		err := handleErrorCondition(latestTransaction.Pan)
		return models.CaptureResponse{}, err
	}

	if !state.TransitionAllowed(latestTransaction.OperationType, state.Capture) {
		return models.CaptureResponse{}, models.ErrTransitionNowAllowed
	}

	if !captureIsPossible(latestTransaction, storeRequest) {
		return models.CaptureResponse{}, models.ErrInvalidAmount

	}
	balanceLeft := latestTransaction.AmountAvailable - amountToCapture

	operationId, err := s.generator.GenerateUUID()
	if err != nil {
		return models.CaptureResponse{}, err
	}
	storeRequest = store.Transaction{
		TransactionId:   capture.TransactionId,
		OperationId:     operationId,
		Amount:          latestTransaction.Amount,
		AmountAvailable: balanceLeft,
		Currency:        capture.Amount.Currency,
		OperationType:   state.Capture,
		Pan:             latestTransaction.Pan,
	}
	log.Println("balance left", balanceLeft)
	err = s.store.Create(&storeRequest)
	if err != nil {
		return models.CaptureResponse{}, fmt.Errorf("couldn't capture %e", err)
	}
	log.Println("response", latestTransaction)
	return models.CaptureResponse{
		TransactionId: capture.TransactionId,
		OperationId:   operationId,
		AmountCharged: models.Amount{
			MajorUnits: amountToCapture,
			Currency:   capture.Amount.Currency,
		},
		AmountAvailable: models.Amount{
			MajorUnits: balanceLeft,
			Currency:   capture.Amount.Currency,
		},
		Response: models.Response{Code: models.Approved},
	}, nil
}

func captureIsPossible(currentState store.Transaction, request store.Transaction) bool {
	log.Println(currentState.AmountAvailable >= request.Amount)
	log.Println(fmt.Sprintf("currencies acc current %s actual %s", currentState.Currency, request.Currency), currentState.Currency == request.Currency)
	return currentState.AmountAvailable >= request.Amount && currentState.Currency == request.Currency
}
