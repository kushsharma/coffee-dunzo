package device_test

import (
	"testing"

	"github.com/kushsharma/coffee-dunzo/device"
	"github.com/kushsharma/coffee-dunzo/models"
	"github.com/stretchr/testify/assert"
)

func TestInmemoryStore(t *testing.T) {
	t.Run("Add", func(t *testing.T) {
		t.Run("should add/subtract item quantity to map successfully", func(t *testing.T) {
			store := device.NewInmemoryStore(make(models.ItemQuantity))
			store.Add(models.Ingrident("water"), models.Quantiy(100))
			assert.Equal(t, models.Quantiy(100), store.Status()["water"])

			store.Add(models.Ingrident("water"), models.Quantiy(-50))
			assert.Equal(t, models.Quantiy(50), store.Status()["water"])
		})
	})
	t.Run("Consume", func(t *testing.T) {
		t.Run("should successfully consume items if we have enough resources", func(t *testing.T) {
			store := device.NewInmemoryStore(models.ItemQuantity{
				models.Ingrident("water"): models.Quantiy(100),
				models.Ingrident("milk"):  models.Quantiy(200),
			})
			err := store.Consume(models.ItemQuantity{
				models.Ingrident("water"): models.Quantiy(10),
				models.Ingrident("milk"):  models.Quantiy(100),
			})
			assert.Nil(t, err)
			status := store.Status()
			assert.Equal(t, models.Quantiy(90), status["water"])
			assert.Equal(t, models.Quantiy(100), status["milk"])
		})
		t.Run("should fail to consume if we do not have enough resources", func(t *testing.T) {
			store := device.NewInmemoryStore(models.ItemQuantity{
				models.Ingrident("water"): models.Quantiy(100),
				models.Ingrident("milk"):  models.Quantiy(200),
			})
			err := store.Consume(models.ItemQuantity{
				models.Ingrident("water"): models.Quantiy(50),
				models.Ingrident("milk"):  models.Quantiy(300),
			})
			assert.Error(t, device.ErrItemNotSufficient, err)
			status := store.Status()
			assert.Equal(t, models.Quantiy(100), status["water"])
			assert.Equal(t, models.Quantiy(200), status["milk"])
		})
		t.Run("should fail to consume if we do not have some specific resource", func(t *testing.T) {
			store := device.NewInmemoryStore(models.ItemQuantity{
				models.Ingrident("water"): models.Quantiy(100),
				models.Ingrident("milk"):  models.Quantiy(200),
			})
			err := store.Consume(models.ItemQuantity{
				models.Ingrident("water"): models.Quantiy(50),
				models.Ingrident("tea"):   models.Quantiy(200),
			})
			assert.Error(t, device.ErrItemNotAvailable, err)
			status := store.Status()
			assert.Equal(t, models.Quantiy(100), status["water"])
			assert.Equal(t, models.Quantiy(200), status["milk"])
		})
	})
	t.Run("Status", func(t *testing.T) {
		t.Run("should provide current status of store", func(t *testing.T) {
			store := device.NewInmemoryStore(models.ItemQuantity{
				models.Ingrident("water"): models.Quantiy(100),
				models.Ingrident("milk"):  models.Quantiy(200),
			})
			err := store.Consume(models.ItemQuantity{
				models.Ingrident("water"): models.Quantiy(50),
			})
			assert.Nil(t, err)
			status := store.Status()
			assert.Equal(t, models.Quantiy(50), status["water"])
			assert.Equal(t, models.Quantiy(200), status["milk"])
		})
	})
}
