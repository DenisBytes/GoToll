package main

import (
	"fmt"

	"github.com/DenisBytes/GoToll/types"
)

type Aggregator interface {
	AggregatorDistance(types.Distance) error
}

type Storer interface {
	Insert(types.Distance) error
}

type InvoiceAggregator struct {
	store Storer
}

func NewInvoiceAggregator (store Storer) Aggregator{
	return &InvoiceAggregator{
		store: store,
	}
}

func (i *InvoiceAggregator) AggregatorDistance(data types.Distance) error {
	fmt.Println("Processing and inserting distance in the storage:", data)
	return i.store.Insert(data)
}