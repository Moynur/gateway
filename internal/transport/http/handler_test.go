package handler_test

import (
	"bytes"
	"encoding/json"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
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

func TestHandler_AuthorizeTransaction(t *testing.T) {
	t.Run("should return bad request when field has invalid type", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ms := service.NewMockService(ctrl)

		h, err := handler.NewHandler(ms)
		assert.NoError(t, err)

		assert.NotNil(t, h)

		handlerAuthReq := []byte(`{"name":false}`)

		capReq, err := json.Marshal(handlerAuthReq)
		assert.NoError(t, err)

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, AuthUrl, bytes.NewReader(capReq))

		h.AuthorizeTransaction(w, r)
		resp := w.Result()
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

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
	t.Run("should return bad request when field has invalid type", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ms := service.NewMockService(ctrl)

		h, err := handler.NewHandler(ms)
		assert.NoError(t, err)

		assert.NotNil(t, h)

		handlerCapReq := []byte(`{"transactionId":false}`)

		capReq, err := json.Marshal(handlerCapReq)
		assert.NoError(t, err)

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, CaptureUrl, bytes.NewReader(capReq))

		h.CaptureTransaction(w, r)
		resp := w.Result()
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("should return 422 when uuid is invalid", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ms := service.NewMockService(ctrl)

		h, err := handler.NewHandler(ms)
		assert.NoError(t, err)

		assert.NotNil(t, h)

		handlerCapReq := handler.CaptureRequest{
			MajorUnits:    1000,
			Currency:      "GBP",
			TransactionId: "TxnId which is not uuid",
		}

		capReq, err := json.Marshal(handlerCapReq)
		assert.NoError(t, err)

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, CaptureUrl, bytes.NewReader(capReq))

		h.CaptureTransaction(w, r)
		resp := w.Result()
		assert.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)

	})

	t.Run("it should do a capture successfully", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ms := service.NewMockService(ctrl)

		h, err := handler.NewHandler(ms)
		assert.NoError(t, err)

		assert.NotNil(t, h)

		handlerCapReq := handler.CaptureRequest{
			MajorUnits:    1000,
			Currency:      "GBP",
			TransactionId: uuid.NewString(),
		}

		capReq, err := json.Marshal(handlerCapReq)
		assert.NoError(t, err)

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, CaptureUrl, bytes.NewReader(capReq))

		newUUID, err := uuid.NewUUID()
		assert.NoError(t, err)

		response := models.CaptureResponse{
			TransactionId: newUUID,
			OperationId:   newUUID,
			AmountCharged: models.Amount{
				MajorUnits: 1000,
				Currency:   "GBP",
			},
			AmountAvailable: models.Amount{
				MajorUnits: 0,
				Currency:   "GBP",
			},
			Response: models.Response{
				Code: 1000,
			},
		}

		ms.EXPECT().Capture(gomock.Any()).Return(response, nil)

		expected := handler.CaptureResponse{
			MajorUnits:       1000,
			AvailableBalance: 0,
			Currency:         "GBP",
			TransactionId:    newUUID.String(),
			OperationId:      newUUID.String(),
			ResponseCode:     1000,
		}

		h.CaptureTransaction(w, r)
		resp := w.Result()

		var actual handler.CaptureResponse
		err = json.NewDecoder(resp.Body).Decode(&actual)

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, expected, actual)
	})
}

func TestHandler_RefundTransaction(t *testing.T) {
	t.Run("should return bad request when field has invalid type", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ms := service.NewMockService(ctrl)

		h, err := handler.NewHandler(ms)
		assert.NoError(t, err)

		assert.NotNil(t, h)

		handlerRefundRequest := []byte(`{"transactionId":false}`)

		capReq, err := json.Marshal(handlerRefundRequest)
		assert.NoError(t, err)

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, RefundUrl, bytes.NewReader(capReq))

		h.RefundTransaction(w, r)
		resp := w.Result()
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("should return 422 when uuid is invalid", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ms := service.NewMockService(ctrl)

		h, err := handler.NewHandler(ms)
		assert.NoError(t, err)

		assert.NotNil(t, h)

		handlerCapReq := handler.RefundRequest{
			TransactionId: "TxnId which is not uuid",
		}

		capReq, err := json.Marshal(handlerCapReq)
		assert.NoError(t, err)

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, RefundUrl, bytes.NewReader(capReq))

		h.RefundTransaction(w, r)
		resp := w.Result()
		assert.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)

	})

	t.Run("it should do a refund successfully", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ms := service.NewMockService(ctrl)

		h, err := handler.NewHandler(ms)
		assert.NoError(t, err)
		assert.NotNil(t, h)

		newUUID, err := uuid.NewUUID()
		assert.NoError(t, err)

		handlerRefundRequest := handler.RefundRequest{
			MajorUnits:    1000,
			Currency:      "GBP",
			TransactionId: newUUID.String(),
		}

		refundReq, err := json.Marshal(handlerRefundRequest)
		assert.NoError(t, err)

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, RefundUrl, bytes.NewReader(refundReq))

		response := models.RefundResponse{
			TransactionId: newUUID,
			OperationId:   newUUID,
			Amount: models.Amount{
				MajorUnits: 1000,
				Currency:   "GBP",
			},
			AmountAvailable: models.Amount{
				MajorUnits: 0,
				Currency:   "GBP",
			},
			Response: models.Response{
				Code: 1000,
			},
		}

		ms.EXPECT().Refund(gomock.Any()).Return(response, nil)

		expected := handler.RefundResponse{
			MajorUnits:       1000,
			AvailableBalance: 0,
			Currency:         "GBP",
			TransactionId:    newUUID.String(),
			OperationId:      newUUID.String(),
			ResponseCode:     1000,
		}

		h.RefundTransaction(w, r)
		resp := w.Result()

		var actual handler.RefundResponse
		err = json.NewDecoder(resp.Body).Decode(&actual)

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, expected, actual)
	})
}

func TestHandler_VoidTransaction(t *testing.T) {
	t.Run("should return bad request when field has invalid type", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ms := service.NewMockService(ctrl)

		h, err := handler.NewHandler(ms)
		assert.NoError(t, err)

		assert.NotNil(t, h)

		handlerVoidReq := []byte(`{"transactionId":false}`)

		voidReq, err := json.Marshal(handlerVoidReq)
		assert.NoError(t, err)

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, VoidUrl, bytes.NewReader(voidReq))

		h.VoidTransaction(w, r)
		resp := w.Result()
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("should return 422 when uuid is invalid", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ms := service.NewMockService(ctrl)

		h, err := handler.NewHandler(ms)
		assert.NoError(t, err)

		assert.NotNil(t, h)

		handlerVoidReq := handler.VoidRequest{TransactionId: "TxnId which is not uuid"}

		voidReq, err := json.Marshal(handlerVoidReq)
		assert.NoError(t, err)

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, VoidUrl, bytes.NewReader(voidReq))

		h.VoidTransaction(w, r)
		resp := w.Result()
		assert.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)

	})

	t.Run("it should do a void successfully", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ms := service.NewMockService(ctrl)

		h, err := handler.NewHandler(ms)
		assert.NoError(t, err)

		assert.NotNil(t, h)

		handlerCapReq := handler.VoidRequest{TransactionId: uuid.NewString()}

		capReq, err := json.Marshal(handlerCapReq)
		assert.NoError(t, err)

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, VoidUrl, bytes.NewReader(capReq))

		newUUID, err := uuid.NewUUID()
		assert.NoError(t, err)

		response := models.AuthVoidResponse{
			TransactionId: newUUID,
			OperationId:   newUUID,
			Response: models.Response{
				Code: 1000,
			},
		}

		ms.EXPECT().Void(gomock.Any()).Return(response, nil)

		expected := handler.VoidResponse{
			TransactionId: newUUID.String(),
			OperationId:   newUUID.String(),
			ResponseCode:  1000,
		}

		h.VoidTransaction(w, r)
		resp := w.Result()

		var actual handler.VoidResponse
		err = json.NewDecoder(resp.Body).Decode(&actual)

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, expected, actual)
	})
}
