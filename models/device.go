package models

type Specification struct {
	Machine Machine
}

type Machine struct {
	Outlets   Outlets
	Stock     ItemQuantity          `json:"total_items_quantity"`
	Beverages map[Item]ItemQuantity `json:"beverages"`
}

type Outlets struct {
	Count int `json:"count_n"`
}

type ItemQuantity map[Ingrident]Quantiy

type Item string

type Ingrident string

type Quantiy int
