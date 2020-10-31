package device

import (
	"errors"
	"fmt"
	"sync"

	"github.com/kushsharma/coffee-dunzo/models"
)

var (
	ErrItemNotAvailable  = errors.New("is not available")
	ErrItemNotSufficient = errors.New("is not sufficient")
)

// inmemoryStore is thread safe in memory models.Store implementation
type inmemoryStore struct {
	mu    sync.Mutex
	items models.ItemQuantity
}

// Add ingrident quantity to stock
func (s *inmemoryStore) Add(ingrident models.Ingrident, quantity models.Quantiy) error {
	s.mu.Lock()
	s.items[ingrident] += quantity
	s.mu.Unlock()
	return nil
}

// checkStock verifies if we have enough items in stock to consume it
func (s *inmemoryStore) checkStock(ingridents models.ItemQuantity) error {
	for item, quantity := range ingridents {
		val, ok := s.items[item]
		if !ok {
			return fmt.Errorf("%s %s", item, ErrItemNotAvailable.Error())
		}
		if val < quantity {
			return fmt.Errorf("%s %s", item, ErrItemNotSufficient.Error())
		}
	}
	return nil
}

// Consumes provide items to the requester and update stock accordingly
func (s *inmemoryStore) Consume(ingridents models.ItemQuantity) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if err := s.checkStock(ingridents); err != nil {
		return err
	}
	for item, quantity := range ingridents {
		s.items[item] -= quantity
	}
	return nil
}

func (s *inmemoryStore) Status() models.ItemQuantity {
	s.mu.Lock()
	defer s.mu.Unlock()
	items := models.ItemQuantity{}
	for k, v := range s.items {
		items[k] = v
	}
	return items
}

func NewInmemoryStore(items models.ItemQuantity) *inmemoryStore {
	return &inmemoryStore{
		items: items,
	}
}
