package service_test

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/moynur/gateway/internal/helpers"
	"github.com/moynur/gateway/internal/models"
	"github.com/moynur/gateway/internal/service"
	"github.com/moynur/gateway/internal/state"
	"github.com/moynur/gateway/internal/store"
	"github.com/stretchr/testify/assert"
)

func TestService_Void(t *testing.T) {
	t.Run("Should Fail void as amount is too large", func(t *testing.T) {
		id := uuid.MustParse("123e4567-e89b-12d3-a456-426655440000")
		id2 := uuid.MustParse("123e4567-0000-12d3-0000-426655440000")

		req := models.RefundRequest{
			TransactionId: id,
			Amount: models.Amount{
				MajorUnits: 1001,
				Currency:   "GBP",
			},
		}
		latestCapture := store.Transaction{
			TransactionId:   id,
			OperationId:     id2,
			Amount:          1000,
			AmountAvailable: 0,
			Currency:        "GBP",
			OperationType:   state.Capture,
			Pan:             "059",
		}
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mg := helpers.NewMockGenerator(ctrl)
		ms := store.NewMockStorer(ctrl)
		service := service.NewService(ms, mg)
		ms.EXPECT().GetLatestByTransactionId(gomock.Any()).Return(latestCapture, nil)
		_, err := service.Refund(req)
		assert.Error(t, err)
		assert.Equal(t, models.ErrInvalidAmount, err)
	})

	t.Run("Refund attempted on a capture that doesn't exist should fail", func(t *testing.T) {
		id := uuid.MustParse("123e4567-e89b-12d3-a456-426655440000")

		req := models.RefundRequest{
			TransactionId: id,
			Amount: models.Amount{
				MajorUnits: 1000,
				Currency:   "GBP",
			},
		}
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mg := helpers.NewMockGenerator(ctrl)
		ms := store.NewMockStorer(ctrl)
		service := service.NewService(ms, mg)
		ms.EXPECT().GetLatestByTransactionId(gomock.Any()).Return(store.Transaction{}, errors.New("some error"))
		_, err := service.Refund(req)
		assert.Error(t, err)
		assert.Equal(t, models.ErrTransactionNotFound, err)
	})

	t.Run("Should do Refund successfully and carry balance", func(t *testing.T) {

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		ms := store.NewMockStorer(ctrl)
		mg := helpers.NewMockGenerator(ctrl)
		service := service.NewService(ms, mg)

		id := uuid.MustParse("123e4567-e89b-12d3-a456-426655440000")
		id2 := uuid.MustParse("123e4567-0000-12d3-0000-426655440000")
		mg.EXPECT().GenerateUUID().Return(id, nil).Return(id2, nil)

		assert.NotEqual(t, id, id2)
		req := models.RefundRequest{
			TransactionId: id,
			Amount: models.Amount{
				MajorUnits: 500,
				Currency:   "GBP",
			},
		}
		initialAuth := store.Transaction{
			TransactionId:   id,
			OperationId:     id2,
			Amount:          1000,
			AmountAvailable: 0, // Fully Captured
			Currency:        "GBP",
			OperationType:   state.Capture,
			Pan:             "059",
		}
		expectedStoreReq := store.Transaction{
			TransactionId: id,
			OperationId:   id2,
			Amount:        1000,
			// This increases again and can reach 1000 at which point the txn is fully refunded
			AmountAvailable: 500,
			Currency:        "GBP",
			OperationType:   state.Refund,
			Pan:             "059",
		}
		expectedResp := models.RefundResponse{
			TransactionId: id,
			OperationId:   id2,
			Response:      models.Response{Code: models.Approved},
			AmountAvailable: models.Amount{
				MajorUnits: 500,
				Currency:   "GBP",
			},
			Amount: models.Amount{
				MajorUnits: 500,
				Currency:   "GBP",
			},
		}
		ms.EXPECT().GetLatestByTransactionId(gomock.Any()).Return(initialAuth, nil)
		ms.EXPECT().Create(&expectedStoreReq).Return(nil)
		resp, err := service.Refund(req)
		assert.NoError(t, err)
		assert.Equal(t, expectedResp, resp)
	})
}
