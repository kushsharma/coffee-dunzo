package device_test

import (
	"sync"
	"testing"

	"github.com/kushsharma/coffee-dunzo/device"
	"github.com/kushsharma/coffee-dunzo/models"
	"github.com/stretchr/testify/assert"
)

func TestMixer(t *testing.T) {
	t.Run("Run", func(t *testing.T) {
		menu := map[models.Item]models.ItemQuantity{
			models.Item("tea"): {
				models.Ingrident("milk"): models.Quantiy(10),
			},
			models.Item("coffee"): {
				models.Ingrident("water"): models.Quantiy(20),
				models.Ingrident("milk"):  models.Quantiy(10),
			},
		}
		t.Run("should consume items from queue once started until finished successfully", func(t *testing.T) {
			mockedStore := new(models.StoreMocked)
			mockedStore.On("Consume", menu[models.Item("tea")]).Return(nil)
			mockedStore.On("Consume", menu[models.Item("coffee")]).Return(nil)
			defer mockedStore.AssertExpectations(t)
			wg := &sync.WaitGroup{}

			errQ := make(chan error)
			requestsQ := make(chan models.Item)
			servedQ := make(chan models.Item)

			mixer := device.NewMixer(menu, mockedStore)
			go func() {
				wg.Add(1)
				mixer.Run(requestsQ, servedQ, errQ)
				wg.Done()
			}()

			requestsQ <- models.Item("tea")
			item := <-servedQ
			assert.Equal(t, models.Item("tea"), item)

			requestsQ <- models.Item("coffee")
			item = <-servedQ
			assert.Equal(t, models.Item("coffee"), item)

			close(requestsQ)
			wg.Wait()
			close(errQ)
			close(servedQ)
		})
		t.Run("should fail to serve if store does not have enough items", func(t *testing.T) {
			mockedStore := new(models.StoreMocked)
			mockedStore.On("Consume", menu[models.Item("coffee")]).Return(device.ErrItemNotSufficient)
			defer mockedStore.AssertExpectations(t)
			wg := &sync.WaitGroup{}

			errQ := make(chan error)
			requestsQ := make(chan models.Item)
			servedQ := make(chan models.Item)

			mixer := device.NewMixer(menu, mockedStore)
			go func() {
				wg.Add(1)
				mixer.Run(requestsQ, servedQ, errQ)
				wg.Done()
			}()

			requestsQ <- models.Item("coffee")
			err := <-errQ
			assert.NotNil(t, err)

			close(requestsQ)
			wg.Wait()
			close(errQ)
			close(servedQ)
		})
		t.Run("should fail to serve if user asked to make unknown item", func(t *testing.T) {
			mockedStore := new(models.StoreMocked)
			defer mockedStore.AssertExpectations(t)
			wg := &sync.WaitGroup{}

			errQ := make(chan error)
			requestsQ := make(chan models.Item)
			servedQ := make(chan models.Item)

			mixer := device.NewMixer(menu, mockedStore)
			go func() {
				wg.Add(1)
				mixer.Run(requestsQ, servedQ, errQ)
				wg.Done()
			}()

			requestsQ <- models.Item("hot_milk")
			err := <-errQ
			assert.NotNil(t, err)

			close(requestsQ)
			wg.Wait()
			close(errQ)
			close(servedQ)
		})
	})
}
