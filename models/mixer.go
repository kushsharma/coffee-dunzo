package models

import "github.com/stretchr/testify/mock"

// Mixer process a device request and prepare the dish
type Mixer interface {
	Run(chan Item, chan Item, chan error)
}

type MixerMocked struct {
	mock.Mock
}

func (m *MixerMocked) Run(r chan Item, s chan Item, e chan error) {
	for _ = range r {
	}
	m.Called(r, s, e)
}
