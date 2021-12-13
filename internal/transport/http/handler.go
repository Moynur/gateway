//go:generate mockgen -package=handler -source=handler.go -destination=./handler_mock.go Handler
package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"

	"github.com/moynur/gateway/internal/models"
	"github.com/moynur/gateway/internal/service"
)

type AuthRequest struct {
	ExpiryMonth string `json:"month"`
	ExpiryYear  string `json:"year"`
	Name        string `json:"name"`
	Postcode    string `json:"postcode"`
	CVV         int    `json:"cvv"`
	PAN         string `json:"pan"`
	MajorUnits  int    `json:"value"`
	Currency    string `json:"currency"`
}

type AuthResponse struct {
	MajorUnits    int    `json:"value"`
	Currency      string `json:"currency"`
	TransactionId string `json:"transactionId"`
	OperationId   string `json:"operationId"`
	ResponseCode  int    `json:"responseCode"`
}

type CaptureRequest struct {
	MajorUnits    int    `json:"value"`
	Currency      string `json:"currency"`
	TransactionId string `json:"transactionId"`
}

type CaptureResponse struct {
	MajorUnits       int    `json:"value"`
	AvailableBalance int    `json:"availableBalance"`
	Currency         string `json:"currency"`
	TransactionId    string `json:"transactionId"`
	OperationId      string `json:"operationId"`
	ResponseCode     int    `json:"responseCode"`
}

type RefundRequest struct {
	MajorUnits    int    `json:"value"`
	Currency      string `json:"currency"`
	TransactionId string `json:"transactionId"`
}

type RefundResponse struct {
	TransactionId    string `json:"transactionId"`
	OperationId      string `json:"operationId"`
	ResponseCode     int    `json:"responseCode"`
	AvailableBalance int    `json:"availableBalance"`
	MajorUnits       int    `json:"value"`
}

type VoidRequest struct {
	TransactionId string `json:"transactionId"`
}

type VoidResponse struct {
	TransactionId string `json:"transactionId"`
	OperationId   string `json:"operationIdd"`
	ResponseCode  int    `json:"responseCode"`
}

type Handler struct {
	service service.Service
}

func NewHandler(svc service.Service) (*Handler, error) {
	return &Handler{
		service: svc,
	}, nil
}

func (h *Handler) ApplyRoutes(r *mux.Router) {
	r.HandleFunc("/authorize", h.AuthorizeTransaction).Methods(http.MethodPost)
	r.HandleFunc("/capture", h.CaptureTransaction).Methods(http.MethodPost)
	r.HandleFunc("/refund", h.RefundTransaction).Methods(http.MethodPost)
	r.HandleFunc("/void", h.VoidTransaction).Methods(http.MethodPost)

}

func (h *Handler) AuthorizeTransaction(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	NewDecoder := json.NewDecoder(r.Body)
	NewDecoder.DisallowUnknownFields()

	var newAuthRequest AuthRequest
	err := NewDecoder.Decode(&newAuthRequest)
	if err != nil {
		log.Println("\n Unable to decode req", err)
		return
	}
	cleanedReq, err := transformAuth(newAuthRequest)
	if err != nil {
		log.Println("Unable to parse req to internal fields", err)
		return
	}

	resp, err := h.service.Authorize(cleanedReq)
	if err != nil {
		log.Println("cant auth", err)
		errorUnknownFailure(w, "failed to auth")
		return
	}
	handlerResp := AuthResponse{
		MajorUnits:    resp.AmountAvailable.MajorUnits,
		Currency:      resp.AmountAvailable.Currency,
		TransactionId: resp.TransactionId,
		OperationId:   resp.OperationId,
		ResponseCode:  resp.Response.AsInt(),
	}
	err = json.NewEncoder(w).Encode(handlerResp)
	if err != nil {
		log.Println("failure to write resp")
		errorUnknownFailure(w, "unknown failure")
		return
	}
}

func (h *Handler) CaptureTransaction(w http.ResponseWriter, r *http.Request) {
	NewDecoder := json.NewDecoder(r.Body)
	NewDecoder.DisallowUnknownFields()

	var newCaptureRequest CaptureRequest
	err := NewDecoder.Decode(&newCaptureRequest)
	if err != nil {
		log.Println("\n Unable to decode req", err)
		errorBadRequest(w, "unable to decode request")
		return
	}
	cleanedReq, err := transformCapture(newCaptureRequest)
	if err != nil {
		log.Println("Unable to parse req to internal fields", err)
		errorUnprocessable(w, "unprocessable entity")
		return
	}

	resp, err := h.service.Capture(cleanedReq)
	if err != nil {
		errorUnknownFailure(w, "unable to capture")
		return
	}

	handlerResp := CaptureResponse{
		MajorUnits:       resp.AmountCharged.MajorUnits,
		AvailableBalance: resp.AmountAvailable.MajorUnits,
		Currency:         resp.AmountCharged.Currency,
		TransactionId:    resp.TransactionId.String(),
		OperationId:      resp.OperationId.String(),
		ResponseCode:     resp.Response.AsInt(),
	}

	err = json.NewEncoder(w).Encode(handlerResp)
	if err != nil {
		errorUnknownFailure(w, "internal server error")
		return
	}
}

func (h *Handler) RefundTransaction(w http.ResponseWriter, r *http.Request) {
	NewDecoder := json.NewDecoder(r.Body)
	NewDecoder.DisallowUnknownFields()

	var newRefundRequest RefundRequest
	err := NewDecoder.Decode(&newRefundRequest)
	log.Println(newRefundRequest, "refund request")
	if err != nil {
		log.Println("\n Unable to decode req", err)
		errorBadRequest(w, "unable to decode request")
		return
	}
	cleanedReq, err := transformRefund(newRefundRequest)
	if err != nil {
		log.Println("Unable to parse req to internal fields", err)
		errorUnprocessable(w, "unprocessable entity")
		return
	}

	resp, err := h.service.Refund(cleanedReq)
	if err != nil {
		log.Println("refund failed", err)
		errorUnknownFailure(w, "refund failed")
		return
	}
	log.Println("refund raw resp", resp)

	handlerResp := RefundResponse{
		TransactionId:    resp.TransactionId.String(),
		OperationId:      resp.OperationId.String(),
		ResponseCode:     resp.Response.AsInt(),
		MajorUnits:       resp.Amount.MajorUnits,
		AvailableBalance: resp.AmountAvailable.MajorUnits,
	}

	log.Println("refund response", handlerResp)
	err = json.NewEncoder(w).Encode(handlerResp)
	if err != nil {
		errorUnknownFailure(w, "internal server error")
		return
	}
}

func (h *Handler) VoidTransaction(w http.ResponseWriter, r *http.Request) {
	NewDecoder := json.NewDecoder(r.Body)
	NewDecoder.DisallowUnknownFields()

	var newVoidRequest VoidRequest
	err := NewDecoder.Decode(&newVoidRequest)
	if err != nil {
		log.Println("\n Unable to decode req", err)
		errorBadRequest(w, "unable to decode request")
		return
	}
	cleanedReq, err := transformVoid(newVoidRequest)
	if err != nil {
		log.Println("Unable to parse req to internal fields", err)
		errorUnprocessable(w, "unprocessable entity")
		return
	}

	resp, err := h.service.Void(cleanedReq)
	if err != nil {
		errorUnknownFailure(w, "void failed")
		return
	}

	handlerResp := VoidResponse{
		TransactionId: resp.TransactionId.String(),
		OperationId:   resp.OperationId.String(),
		ResponseCode:  resp.Response.AsInt(),
	}

	err = json.NewEncoder(w).Encode(handlerResp)
	if err != nil {
		errorUnknownFailure(w, "internal server error")
		return
	}
}

func transformAuth(req AuthRequest) (models.AuthRequest, error) {
	return models.AuthRequest{
		CardInformation: models.CardData{
			Name:     req.Name,
			Postcode: req.Postcode,
		},
		Expiry: models.Expiry{
			Month: req.ExpiryMonth,
			Year:  req.ExpiryYear,
		},
		Amount: models.Amount{
			MajorUnits: req.MajorUnits,
			Currency:   req.Currency,
		},
		PAN: req.PAN,
		CVV: req.CVV,
	}, nil
}

func transformCapture(req CaptureRequest) (models.CaptureRequest, error) {
	u, err := uuid.Parse(req.TransactionId)
	if err != nil {
		return models.CaptureRequest{}, errors.New("invalid transactionId")
	}
	return models.CaptureRequest{
		TransactionId: u,
		Amount: models.Amount{
			MajorUnits: req.MajorUnits,
			Currency:   req.Currency,
		},
	}, nil
}

func transformRefund(req RefundRequest) (models.RefundRequest, error) {
	u, err := uuid.Parse(req.TransactionId)
	if err != nil {
		return models.RefundRequest{}, errors.New("invalid transactionId")
	}
	return models.RefundRequest{
		TransactionId: u,
		Amount: models.Amount{
			MajorUnits: req.MajorUnits,
			Currency:   req.Currency,
		},
	}, nil
}

func transformVoid(req VoidRequest) (models.AuthVoidRequest, error) {
	u, err := uuid.Parse(req.TransactionId)
	if err != nil {
		return models.AuthVoidRequest{}, errors.New("invalid transactionId")
	}
	return models.AuthVoidRequest{
		TransactionId: u,
	}, nil
}

func errorUnknownFailure(w http.ResponseWriter, message string) {
	writeError(w, message, http.StatusInternalServerError)
}

func errorBadRequest(w http.ResponseWriter, message string) {
	writeError(w, message, http.StatusBadRequest)
}

func errorUnprocessable(w http.ResponseWriter, message string) {
	writeError(w, message, http.StatusUnprocessableEntity)
}

func writeError(w http.ResponseWriter, message string, code int) {
	var err error
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	enc := json.NewEncoder(w)

	if encoderErr := enc.Encode(code); err != nil {
		err = encoderErr
	}
}
