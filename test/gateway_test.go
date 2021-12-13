//go:build integration

package gateway_test

import (
	"encoding/json"
	"github.com/go-resty/resty/v2"
	handler "github.com/moynur/gateway/internal/transport/http"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

const (
	BaseUrl         = "http://localhost:8080"
	AuthEndpoint    = "/authorize"
	CaptureEndpoint = "/capture"
	RefundEndpoint  = "/refund"
	VoidEndpoint    = "/void"
)

func TestAuth(t *testing.T) {
	t.Run("can do an auth", func(t *testing.T) {

		log.Println("starting test")
		client := resty.New()
		resp, err := client.R().
			SetBody(`{
    "month": "12",
    "year": "2020",
    "name": "my name is",
    "postcode": "my postcode is",
    "cvv": 123,
    "pan": "059",
    "currency": "GBP",
    "value": 123
}`).Post(BaseUrl + AuthEndpoint)
		assert.NoError(t, err)

		assert.Equal(t, 200, resp.StatusCode())
	})

	t.Run("Should fail auth due to PAN", func(t *testing.T) {

		log.Println("starting test")
		client := resty.New()
		resp, err := client.R().
			SetBody(`{
    "month": "12",
    "year": "2020",
    "name": "name",
    "postcode": "postcode",
    "cvv": 123,
    "pan": "4000000000000119",
    "currency": "GBP",
    "value": 123
}`).Post(BaseUrl + AuthEndpoint)
		assert.NoError(t, err)

		assert.Equal(t, 500, resp.StatusCode())
	})
}

func TestVoid(t *testing.T) {
	t.Run("Should be able to void a payment", func(t *testing.T) {

		log.Println("starting test")
		client := resty.New()
		resp, err := client.R().
			SetBody(`{
    "month": "12",
    "year": "2020",
    "name": "my name is",
    "postcode": "my postcode is",
    "cvv": 123,
    "pan": "059",
    "currency": "GBP",
    "value": 1000
}`).Post(BaseUrl + AuthEndpoint)

		var newAuthResponse = new(handler.AuthResponse)
		err = json.Unmarshal(resp.Body(), &newAuthResponse)
		assert.NoError(t, err)

		assert.Equal(t, 200, resp.StatusCode())
		assert.Equal(t, 1000, newAuthResponse.ResponseCode)

		voidRequest := handler.VoidRequest{TransactionId: newAuthResponse.TransactionId}

		voidResp, err := client.R().SetBody(voidRequest).Post(BaseUrl + VoidEndpoint)

		var voidResponse = new(handler.VoidResponse)
		err = json.Unmarshal(voidResp.Body(), &voidResponse)
		assert.NoError(t, err)
		assert.Equal(t, 200, resp.StatusCode())
		assert.Equal(t, 1000, voidResponse.ResponseCode)
	})
}

func TestRefund(t *testing.T) {
	t.Run("Should do an Auth and Capture in 2 parts then refund the whole amount", func(t *testing.T) {

		log.Println("starting test")
		client := resty.New()
		resp, err := client.R().
			SetBody(`{
    "month": "12",
    "year": "2020",
    "name": "my name is",
    "postcode": "my postcode is",
    "cvv": 123,
    "pan": "059",
    "currency": "GBP",
    "value": 100
}`).Post(BaseUrl + AuthEndpoint)

		var newAuthResponse handler.AuthResponse
		err = json.Unmarshal(resp.Body(), &newAuthResponse)
		assert.NoError(t, err)

		assert.Equal(t, 200, resp.StatusCode())
		assert.Equal(t, 1000, newAuthResponse.ResponseCode)

		captureRequest := handler.CaptureRequest{
			MajorUnits:    newAuthResponse.MajorUnits / 2, // should allow 2 captures for half amount
			Currency:      newAuthResponse.Currency,
			TransactionId: newAuthResponse.TransactionId,
		}

		captureResp, err := client.R().SetBody(captureRequest).Post(BaseUrl + CaptureEndpoint)

		var firstCaptureResponse handler.CaptureResponse
		err = json.Unmarshal(captureResp.Body(), &firstCaptureResponse)
		assert.NoError(t, err)
		assert.Equal(t, 200, resp.StatusCode())
		assert.Equal(t, 1000, firstCaptureResponse.ResponseCode)
		log.Println(firstCaptureResponse)
		assert.Equal(t, 50, firstCaptureResponse.MajorUnits)

		captureResp, err = client.R().SetBody(captureRequest).Post(BaseUrl + CaptureEndpoint)

		var secondCaptureResponse handler.CaptureResponse
		err = json.Unmarshal(captureResp.Body(), &secondCaptureResponse)
		assert.NoError(t, err)
		assert.Equal(t, 200, resp.StatusCode())
		assert.Equal(t, 1000, secondCaptureResponse.ResponseCode)
		assert.Equal(t, 50, secondCaptureResponse.MajorUnits)

		RefundRequest := handler.RefundRequest{
			MajorUnits: newAuthResponse.MajorUnits,
			// could be that we don't want to allow this refund as it's larger than a single presentment
			Currency:      newAuthResponse.Currency,
			TransactionId: newAuthResponse.TransactionId,
		}

		refundResp, err := client.R().SetBody(RefundRequest).Post(BaseUrl + RefundEndpoint)
		assert.Equal(t, 200, resp.StatusCode())

		log.Println("refund response", refundResp.RawBody())
		var newRefundResponse handler.RefundResponse
		err = json.Unmarshal(refundResp.Body(), &newRefundResponse)
		assert.NoError(t, err)
		assert.Equal(t, 1000, newRefundResponse.ResponseCode)
		assert.Equal(t, 100, newRefundResponse.AvailableBalance)
	})
}
