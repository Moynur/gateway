package service

import (
	"fmt"
	"log"

	"github.com/moynur/gateway/internal/models"
	"github.com/moynur/gateway/internal/state"
	"github.com/moynur/gateway/internal/store"
)

func (s *service) Void(void models.AuthVoidRequest) (models.AuthVoidResponse, error) {
	storeRequest := store.Transaction{
		TransactionId: void.TransactionId,
	}

	resp, err := s.store.GetLatestByTransactionId(storeRequest.TransactionId)
	if err != nil {
		return models.AuthVoidResponse{}, fmt.Errorf("couldn't Fetch %e", err)
	}

	if !state.TransitionAllowed(resp.OperationType, state.Void) {
		return models.AuthVoidResponse{}, models.ErrTransitionNowAllowed
	}

	operationId, err := s.generator.GenerateUUID()
	if err != nil {
		return models.AuthVoidResponse{}, models.ErrGeneric
	}
	storeRequest = store.Transaction{
		TransactionId:   void.TransactionId,
		OperationId:     operationId,
		Amount:          resp.Amount,
		AmountAvailable: 0,
		Currency:        resp.Currency,
		OperationType:   state.Void,
	}
	log.Println("calling store")
	err = s.store.Create(&storeRequest)
	if err != nil {
		return models.AuthVoidResponse{}, fmt.Errorf("couldn't void %e", err)
	}
	log.Println("response", resp)
	return models.AuthVoidResponse{
		TransactionId: void.TransactionId,
		OperationId:   operationId,
		Response:      models.Response{Code: models.Approved},
		Amount: models.Amount{
			MajorUnits: resp.Amount,
			Currency:   resp.Currency,
		},
	}, nil
}
