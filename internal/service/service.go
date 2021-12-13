//go:generate mockgen -package=service -source=service.go -destination=./service_mock.go Service
package service

import (
	"github.com/moynur/gateway/internal/helpers"
	"github.com/moynur/gateway/internal/models"
	"github.com/moynur/gateway/internal/store"
)

type Service interface {
	Authorize(models.AuthRequest) (models.AuthResponse, error)
	Capture(models.CaptureRequest) (models.CaptureResponse, error)
	Refund(models.RefundRequest) (models.RefundResponse, error)
	Void(models.AuthVoidRequest) (models.AuthVoidResponse, error)
}

var (
	panAuthFailure    = "4000000000000119"
	panCaptureFailure = "4000000000000259"
	panRefundFailure  = "4000000000003238"
)

type service struct {
	store     store.Storer
	generator helpers.Generator
}

const (
	Authorize = "Authorize"
	Capture   = "Capture"
	Refund    = "Refund"
)

func NewService(db store.Storer, generator helpers.Generator) *service {
	return &service{
		store:     db,
		generator: generator,
	}
}

func isSpecialPan(pan string, method string) bool {
	return (pan == panAuthFailure && method == Authorize) ||
		(pan == panCaptureFailure && method == Capture) ||
		(pan == panRefundFailure && method == Refund)
}

func handleErrorCondition(pan string) error {
	switch pan {
	case panAuthFailure:
		return models.ErrCantAuth
	case panCaptureFailure:
		return models.ErrCantCapture
	case panRefundFailure:
		return models.ErrCantRefund
	default:
		return nil
	}
}
