package handler_test

import (
	"bytes"
	"encoding/json"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/moynur/gateway/internal/models"
	"github.com/moynur/gateway/internal/service"
	handler "github.com/moynur/gateway/internal/transport/http"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

const (
	AuthUrl    = "/authorize"
	CaptureUrl = "/capture"
	RefundUrl  = "/refund"
	VoidUrl    = "/void"
)

func TestHandler_NewHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ms := service.NewMockService(ctrl)

	h, err := handler.NewHandler(ms)
	assert.NoError(t, err)

	assert.NotNil(t, h)

}

func TestHandler_ApplyRoutes(t *testing.T) {
}

func TestHandler_AuthorizeTransaction(t *testing.T) {

	t.Run("should return approved", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ms := service.NewMockService(ctrl)

		h, err := handler.NewHandler(ms)
		assert.NoError(t, err)

		assert.NotNil(t, h)

		handlerAuthReq := handler.AuthRequest{
			ExpiryMonth: "month",
			ExpiryYear:  "year",
			Name:        "name",
			Postcode:    "postcode",
			CVV:         123,
			PAN:         "059",
			MajorUnits:  1000,
			Currency:    "GBP",
		}

		authReqMarshalled, err := json.Marshal(handlerAuthReq)
		assert.NoError(t, err)

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, AuthUrl, bytes.NewReader(authReqMarshalled))
		r = mux.SetURLVars(r, map[string]string{
			"bla": "bla",
		})

		expectedServerResp := models.AuthResponse{
			TransactionId: "TransactionUUID",
			OperationId:   "OperationUUID",
			Response:      models.Response{Code: models.Approved},
			AmountAvailable: models.Amount{
				MajorUnits: 1000,
				Currency:   "GBP",
			},
		}

		ms.EXPECT().Authorize(gomock.Any()).Return(expectedServerResp, nil)

		h.AuthorizeTransaction(w, r)
		resp := w.Result()
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		err = resp.Body.Close()
		assert.NoError(t, err)

		expected := handler.AuthResponse{
			MajorUnits:    1000,
			Currency:      expectedServerResp.AmountAvailable.Currency,
			TransactionId: expectedServerResp.TransactionId,
			OperationId:   expectedServerResp.OperationId,
			ResponseCode:  expectedServerResp.Response.AsInt(),
		}

		var out handler.AuthResponse
		assert.NoError(t, json.NewDecoder(resp.Body).Decode(&out))
		assert.Equal(t, expected, out)
	})
}

func TestHandler_CaptureTransaction(t *testing.T) {

	t.Run("should return approved", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ms := service.NewMockService(ctrl)

		h, err := handler.NewHandler(ms)
		assert.NoError(t, err)

		assert.NotNil(t, h)

		handlerAuthReq := handler.CaptureRequest{
			MajorUnits:    1000,
			Currency:      "GBP",
			TransactionId: "TxnId",
		}

		authReqMarshalled, err := json.Marshal(handlerAuthReq)
		assert.NoError(t, err)

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, AuthUrl, bytes.NewReader(authReqMarshalled))
		r = mux.SetURLVars(r, map[string]string{
			"bla": "bla",
		})

		expectedServerResp := models.AuthResponse{
			TransactionId: "TransactionUUID",
			OperationId:   "OperationUUID",
			Response:      models.Response{Code: models.Approved},
			AmountAvailable: models.Amount{
				MajorUnits: 1000,
				Currency:   "GBP",
			},
		}

		ms.EXPECT().Authorize(gomock.Any()).Return(expectedServerResp, nil)

		h.AuthorizeTransaction(w, r)
		resp := w.Result()
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		err = resp.Body.Close()
		assert.NoError(t, err)

		expected := handler.AuthResponse{
			MajorUnits:    1000,
			Currency:      expectedServerResp.AmountAvailable.Currency,
			TransactionId: expectedServerResp.TransactionId,
			OperationId:   expectedServerResp.OperationId,
			ResponseCode:  expectedServerResp.Response.AsInt(),
		}

		var out handler.AuthResponse
		assert.NoError(t, json.NewDecoder(resp.Body).Decode(&out))
		assert.Equal(t, expected, out)
	})
}
