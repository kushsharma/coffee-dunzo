package models

import "github.com/stretchr/testify/mock"

// Store handles stock of a device
// should be able to add or consume item
// should provice current stock status
type Store interface {
	Add(Ingrident, Quantiy) error
	Consume(ItemQuantity) error
	Status() ItemQuantity
}

type StoreMocked struct {
	mock.Mock
}

func (m *StoreMocked) Add(i Ingrident, q Quantiy) error {
	return m.Called(i, q).Error(0)
}

func (m *StoreMocked) Consume(iq ItemQuantity) error {
	return m.Called(iq).Error(0)
}

func (m *StoreMocked) Status() ItemQuantity {
	return m.Called().Get(0).(ItemQuantity)
}
