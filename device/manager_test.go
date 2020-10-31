package device_test

import (
	"testing"

	"github.com/kushsharma/coffee-dunzo/device"
	"github.com/kushsharma/coffee-dunzo/models"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestManager(t *testing.T) {
	logger := logrus.New()
	t.Run("RefillRequired", func(t *testing.T) {
		t.Run("should successfully return warning level ingridents", func(t *testing.T) {
			minStockWarning := models.ItemQuantity{
				models.Ingrident("milk"):  models.Quantiy(50),
				models.Ingrident("water"): models.Quantiy(5),
			}
			mockedStore := new(models.StoreMocked)
			mockedStore.On("Status").Return(models.ItemQuantity{
				models.Ingrident("milk"): models.Quantiy(500),
			})
			defer mockedStore.AssertExpectations(t)

			manager := device.NewManager(minStockWarning, mockedStore, logger)
			defer manager.Stop()

			almostEmpty := manager.CheckRefill()

			assert.Equal(t, 1, len(almostEmpty))
			assert.Equal(t, models.Ingrident("water"), almostEmpty[0])
		})
	})
	t.Run("FillIngrident", func(t *testing.T) {
		t.Run("should successfully fill ingridents", func(t *testing.T) {
			mockedStore := new(models.StoreMocked)
			mockedStore.On("Add", models.Ingrident("milk"), models.Quantiy(500)).Return(nil)
			defer mockedStore.AssertExpectations(t)

			manager := device.NewManager(models.ItemQuantity{}, mockedStore, logger)
			defer manager.Stop()

			manager.FillIngrident(models.Ingrident("milk"), models.Quantiy(500))
		})
	})
	t.Run("RequestItem", func(t *testing.T) {
		t.Run("should not allow to request item after machine is stopped", func(t *testing.T) {
			mockedStore := new(models.StoreMocked)
			defer mockedStore.AssertExpectations(t)

			manager := device.NewManager(models.ItemQuantity{}, mockedStore, logger)
			manager.Stop()

			err := manager.RequestItem(models.Item("milk"))
			assert.NotNil(t, err)
		})
	})
	t.Run("CreateOutlet", func(t *testing.T) {
		t.Run("should successfully create outlet that consumes from requests Queue", func(t *testing.T) {
			mockedStore := new(models.StoreMocked)
			defer mockedStore.AssertExpectations(t)

			mockedMixer := new(models.MixerMocked)
			mockedMixer.On("Run", mock.AnythingOfType("chan models.Item"), mock.AnythingOfType("chan models.Item"), mock.AnythingOfType("chan error")).Once()
			defer mockedMixer.AssertExpectations(t)

			manager := device.NewManager(models.ItemQuantity{}, mockedStore, logger)
			manager.CreateOutlet(mockedMixer)
			manager.RequestItem(models.Item("tea"))
			manager.RequestItem(models.Item("coffee"))

			manager.Stop()
		})
	})
}
