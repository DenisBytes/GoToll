package main

import (
	"time"

	"github.com/DenisBytes/GoToll/types"
	"github.com/sirupsen/logrus"
)

type LogMiddleWare struct {
	next Aggregator
}

func NewLogMiddleware(next Aggregator) Aggregator {
	return &LogMiddleWare{
		next: next,
	}
}

func (l *LogMiddleWare) AggregatorDistance(data types.Distance) (err error) {
	defer func (start time.Time){
		logrus.WithFields(logrus.Fields{
			"took": time.Since(start),
			"err": err,
		}).Info("Aggregate distance")
	}(time.Now())
	err = l.next.AggregatorDistance(data)
	return 
}

func (l *LogMiddleWare) CalculateInvoice(obuID int) (invoice *types.Invoice, err error) {
	defer func (start time.Time){
		var (
			distance float64
			amount float64
		)
		if invoice !=  nil {
			distance = invoice.TotalDistance
			amount = invoice.TotalAmount
		}
		logrus.WithFields(logrus.Fields{
			"took": time.Since(start),
			"err": err,
			"obuID": obuID,
			"amount": amount,
			"distance": distance,
		}).Info("Calculate Invoice")
	}(time.Now())
	invoice, err = l.next.CalculateInvoice(obuID)
	return 
}