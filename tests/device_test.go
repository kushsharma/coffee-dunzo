// functional tests
package tests

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"regexp"
	"testing"

	"github.com/kushsharma/coffee-dunzo/device"
	"github.com/kushsharma/coffee-dunzo/models"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	easy "github.com/t-tomalak/logrus-easy-formatter"
)

func TestDevice(t *testing.T) {

	testCases := []struct {
		InputFile    string
		Validate     func(string) error
		Requests     func(*device.Manager, *log.Logger)
		WarningLevel float64
	}{
		{
			InputFile: "./sample_ingridents.json",
			Validate: func(output string) error {
				err := errors.New("output mismatch")

				hotTea := regexp.MustCompile(`hot_tea is prepared`)
				if !hotTea.MatchString(output) {
					return errors.Wrap(err, `hot_tea is prepared`)
				}

				greenTea := regexp.MustCompile(`green_tea cannot be prepared because vanilla is not available`)
				if !greenTea.MatchString(output) {
					return errors.Wrap(err, `green_tea cannot be prepared because vanilla is not available`)
				}

				refilHotWater := regexp.MustCompile(`refilling hot_water with : 1000`)
				if !refilHotWater.MatchString(output) {
					return errors.Wrap(err, `refilling hot_water with : 1000`)
				}

				return nil
			},
			Requests: func(machine *device.Manager, logger *log.Logger) {
				machine.RequestItem(models.Item("hot_tea"))
				machine.RequestItem(models.Item("hot_coffee"))
				machine.RequestItem(models.Item("green_tea"))
				machine.RequestItem(models.Item("black_tea"))

				refillItems := machine.CheckRefill()
				for _, item := range refillItems {
					logger.Info("refilling ", item, " with : 1000")
					machine.FillIngrident(item, models.Quantiy(1000))
				}
				machine.RequestItem(models.Item("hot_coffee"))
			},
			WarningLevel: 0.5,
		},
	}

	t.Run("integration tests for device", func(t *testing.T) {
		for _, testCase := range testCases {
			func() {
				var buf bytes.Buffer
				logger := log.New()
				logger.SetOutput(&buf)
				logger.SetLevel(log.InfoLevel)
				logger.SetFormatter(&easy.Formatter{
					LogFormat: "%msg%\n",
				})

				// read specification
				specificationString, err := ioutil.ReadFile(testCase.InputFile)
				if err != nil {
					panic(err)
				}
				var specifications models.Specification
				if err := json.Unmarshal(specificationString, &specifications); err != nil {
					panic(err)
				}

				// simulate warning levels
				minStockWarningPercentage := testCase.WarningLevel
				minimumStockWarning := models.ItemQuantity{}
				for ingrident, quantity := range specifications.Machine.Stock {
					minimumStockWarning[ingrident] = models.Quantiy(float64(quantity) * minStockWarningPercentage)
				}

				// create a machine that knows how to make requested menu beverages
				store := device.NewInmemoryStore(models.ItemQuantity{})
				machine := device.NewManager(minimumStockWarning, store, logger)
				for i := 0; i < specifications.Machine.Outlets.Count; i++ {
					machine.CreateOutlet(device.NewMixer(specifications.Machine.Beverages, store))
				}
				logger.Infof("machine started with %d outlets", specifications.Machine.Outlets.Count)

				// simulate stock
				for ingrident, quantity := range specifications.Machine.Stock {
					machine.FillIngrident(ingrident, quantity)
				}

				// simulate users
				testCase.Requests(machine, logger)

				machine.Stop()

				finalOutput := buf.String()
				assert.Nil(t, testCase.Validate(finalOutput))
				//t.Error(finalOutput)
			}()
		}
	})

}
