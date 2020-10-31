package device

import (
	"errors"
	"fmt"
	"time"

	"github.com/kushsharma/coffee-dunzo/models"
)

var (
	ErrUndefinedItem = errors.New("undefined item for mixer")
)

// mixer works with the help of a menu and produces items based on
// user requests one by one
type mixer struct {
	menu  map[models.Item]models.ItemQuantity
	store models.Store
}

// Run takes jobs from requests queue and producess successful results
// on serve channel. If anything goes wrong in between, errors are pushed to
// a error queue and mixer keeps on running.
func (m *mixer) Run(requests chan models.Item, serve chan models.Item, errs chan error) {
	for item := range requests {
		requirements, ok := m.menu[item]
		if !ok {
			errs <- ErrUndefinedItem
			continue
		}
		if err := m.store.Consume(requirements); err != nil {
			errs <- fmt.Errorf("%s cannot be prepared because %s", item, err.Error())
			continue
		}
		// simulate cooking
		time.Sleep(time.Second * 2)
		serve <- item
	}
}

func NewMixer(menu map[models.Item]models.ItemQuantity, store models.Store) *mixer {
	return &mixer{
		menu:  menu,
		store: store,
	}
}
