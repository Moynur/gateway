package service

import (
	"fmt"
	"log"

	luhn "github.com/Moynur/Exercism/go/luhn"
	"github.com/moynur/gateway/internal/models"
	"github.com/moynur/gateway/internal/state"
	"github.com/moynur/gateway/internal/store"
)

func (s *service) Authorize(auth models.AuthRequest) (models.AuthResponse, error) {
	pan := auth.Card.PAN
	if isSpecialPan(pan, Authorize) {
		err := handleErrorCondition(pan)
		return models.AuthResponse{}, err
	}
	if !luhn.Valid(pan) {
		log.Println("failed luhn")
		return models.AuthResponse{}, models.ErrFailedLuhn
	}

	OperationId, err := s.generator.GenerateUUID()
	if err != nil {
		return models.AuthResponse{}, err
	}
	TransactionId, err := s.generator.GenerateUUID()
	if err != nil {
		return models.AuthResponse{}, err
	}
	storeRequest := store.Transaction{
		TransactionId:   TransactionId,
		OperationId:     OperationId,
		Amount:          auth.Amount.MajorUnits,
		AmountAvailable: auth.Amount.MajorUnits,
		Currency:        auth.Amount.Currency,
		OperationType:   state.Auth,
		Pan:             pan,
	}
	log.Println("calling store")
	err = s.store.Create(&storeRequest)
	if err != nil {
		return models.AuthResponse{}, fmt.Errorf("couldn't authorize %e", err)
	}
	return models.AuthResponse{
		TransactionId: s.generator.AsString(TransactionId),
		OperationId:   s.generator.AsString(OperationId),
		Response:      models.Response{Code: models.Approved},
		AmountAvailable: models.Amount{
			MajorUnits: auth.Amount.MajorUnits,
			Currency:   auth.Amount.Currency,
		},
	}, nil
}
