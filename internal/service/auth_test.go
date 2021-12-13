package service_test

import (
	"github.com/google/uuid"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/moynur/gateway/internal/helpers"
	"github.com/moynur/gateway/internal/models"

	"github.com/moynur/gateway/internal/service"
	"github.com/moynur/gateway/internal/store"
	"github.com/stretchr/testify/assert"
)

func TestService_Authorize(t *testing.T) {
	t.Run("Should Fail Auth if Pan is special type which fails on Auth", func(t *testing.T) {
		req := models.AuthRequest{
			Card: models.Card{PAN: "4000000000000119"},
		}
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mg := helpers.NewMockGenerator(ctrl)
		ms := store.NewMockStorer(ctrl)
		newService := service.NewService(ms, mg)
		_, err := newService.Authorize(req)
		assert.Error(t, err)
	})

	t.Run("Should Fail Auth if pan is invalid", func(t *testing.T) {
		req := models.AuthRequest{
			Card: models.Card{PAN: "0"},
		}
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mg := helpers.NewMockGenerator(ctrl)
		ms := store.NewMockStorer(ctrl)
		newService := service.NewService(ms, mg)
		_, err := newService.Authorize(req)
		assert.Error(t, err)
	})

	t.Run("Should Fail Auth if pan is invalid", func(t *testing.T) {
		req := models.AuthRequest{
			Card: models.Card{PAN: "0"},
		}
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mg := helpers.NewMockGenerator(ctrl)
		ms := store.NewMockStorer(ctrl)
		newService := service.NewService(ms, mg)
		_, err := newService.Authorize(req)
		assert.Error(t, err)
	})

	t.Run("Should throw unknown error when doing auth", func(t *testing.T) {
		req := models.AuthRequest{
			Card: models.Card{PAN: "059"},
		}
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mg := helpers.NewMockGenerator(ctrl)
		ms := store.NewMockStorer(ctrl)
		newService := service.NewService(ms, mg)
		mg.EXPECT().GenerateUUID().Return(uuid.NewUUID()).Times(2)
		ms.EXPECT().Create(gomock.Any()).Return(models.ErrInvalidAmount)
		_, err := newService.Authorize(req)
		assert.Error(t, err)
	})

	t.Run("Should do auth successfully", func(t *testing.T) {
		req := models.AuthRequest{
			Card: models.Card{
				Name:     "some name",
				Postcode: "some postcode",
				Expiry: models.Expiry{
					Month: "08",
					Year:  "2021",
				},
				PAN: "059",
				CVV: 123,
			},
			Amount: models.Amount{
				MajorUnits: 1000,
				Currency:   "GBP",
			},
		}
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mg := helpers.NewMockGenerator(ctrl)
		ms := store.NewMockStorer(ctrl)
		newService := service.NewService(ms, mg)
		mg.EXPECT().GenerateUUID().Return(uuid.NewUUID()).Times(2)
		mg.EXPECT().AsString(gomock.Any()).Return("a uuid as a string").Times(2)
		ms.EXPECT().Create(gomock.Any()).Return(nil)
		resp, err := newService.Authorize(req)
		assert.NoError(t, err)
		assert.Equal(t, models.Approved, resp.Response.Code)
	})
}
