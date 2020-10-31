package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"os/signal"
	"time"

	"github.com/kushsharma/coffee-dunzo/device"
	"github.com/kushsharma/coffee-dunzo/models"
	log "github.com/sirupsen/logrus"
)

var (
	termChan                  = make(chan os.Signal, 1)
	minStockWarningPercentage = 0.5
)

func main() {
	logger := log.New()
	logger.SetOutput(os.Stdout)
	logger.SetLevel(log.InfoLevel)
	logger.SetFormatter(&log.TextFormatter{
		DisableColors: false,
		FullTimestamp: true,
	})

	log.Info(`this application is NOT made to run standalone as requested in problem statement
	but for demonstration, here is simulation of a machine with 4 outlets that is taking some
	user input, checking for refill status, filling it once and then terminating`)

	// read specification
	specificationString, err := ioutil.ReadFile("./ingridents.json")
	if err != nil {
		panic(err)
	}
	var specifications models.Specification
	if err := json.Unmarshal(specificationString, &specifications); err != nil {
		panic(err)
	}

	// simulate warning levels
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
	machine.RequestItem(models.Item("hot_tea"))
	machine.RequestItem(models.Item("hot_coffee"))
	machine.RequestItem(models.Item("green_tea"))
	machine.RequestItem(models.Item("black_tea"))
	simulateRefill(machine, logger)
	machine.RequestItem(models.Item("hot_coffee"))

	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	signal.Notify(termChan, os.Interrupt)
	signal.Notify(termChan, os.Kill)

	// Block until we receive our signal.
	<-termChan
	logger.Info("termination request received")

	// stop taking any more request and wait for existing requests to complete
	machine.Stop()
	logger.Info("bye")
}

func simulateRefill(machine *device.Manager, logger *log.Logger) {
	refillItems := machine.CheckRefill()
	for _, item := range refillItems {
		logger.Info("refilling ", item, " with : 1000")
		machine.FillIngrident(item, models.Quantiy(1000))
	}
	time.Sleep(time.Second * 2)
}
