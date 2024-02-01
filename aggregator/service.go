package main

import (
	"fmt"

	"github.com/DenisBytes/GoToll/types"
)

const (
	basePrice = 3.15
)

type Aggregator interface {
	AggregatorDistance(types.Distance) error
	CalculateInvoice(int) (*types.Invoice, error)
}

type Storer interface {
	Insert(types.Distance) error
	Get(int) (float64, error)
}

type InvoiceAggregator struct {
	store Storer
}

func NewInvoiceAggregator (store Storer) Aggregator{
	return &InvoiceAggregator{
		store: store,
	}
}

func(i *InvoiceAggregator) AggregatorDistance(data types.Distance) error {
	fmt.Println("Processing and inserting distance in the storage:", data)
	return i.store.Insert(data)
}

func(i *InvoiceAggregator) CalculateInvoice(obuID int) (*types.Invoice, error) {
	dist, err := i.store.Get(obuID)
	if err != nil {
		return nil, err
	}

	inv := &types.Invoice{
		OBUID: obuID,
		TotalDistance: dist,
		TotalAmount: basePrice * dist,
	}

	return inv, nil
}
