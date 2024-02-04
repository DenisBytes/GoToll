package main

import (
	"time"

	"github.com/DenisBytes/GoToll/types"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/sirupsen/logrus"
)

type MetricsMiddleware struct {
	errReqCounterAgg  prometheus.Counter
	errReqCounterCalc prometheus.Counter
	reqCounterAgg     prometheus.Counter
	reqCounterCalc    prometheus.Counter
	reqLatencyAgg     prometheus.Histogram
	reqLatencyCalc    prometheus.Histogram
	next              Aggregator
}

func NewMetricsMiddleware(next Aggregator) *MetricsMiddleware {

	errReqCounterAgg := promauto.NewCounter(prometheus.CounterOpts{
		Namespace: "aggregator_error_counter",
		Name:      "aggregate",
	})

	errReqCounterCalc := promauto.NewCounter(prometheus.CounterOpts{
		Namespace: "aggregator_error_counter",
		Name:      "calculate",
	})

	reqCounterAgg := promauto.NewCounter(prometheus.CounterOpts{
		Namespace: "aggregator_request_counter",
		Name:      "aggregate",
	})

	reqCounterCalc := promauto.NewCounter(prometheus.CounterOpts{
		Namespace: "aggregator_request_counter",
		Name:      "calculate",
	})

	reqLatencyAgg := promauto.NewHistogram(prometheus.HistogramOpts{
		Namespace: "aggregator_request_latency_1",
		Name:      "aggregate",
		Buckets:   []float64{0.1, 0.5, 1.00},
	})

	reqLatencyCalc := promauto.NewHistogram(prometheus.HistogramOpts{
		Namespace: "aggregator_request_latency_2",
		Name:      "calculate",
		Buckets:   []float64{0.1, 0.5, 1.00},
	})

	return &MetricsMiddleware{
		next:              next,
		errReqCounterAgg:  errReqCounterAgg,
		errReqCounterCalc: errReqCounterCalc,
		reqCounterAgg:     reqCounterAgg,
		reqCounterCalc:    reqCounterCalc,
		reqLatencyAgg:     reqLatencyAgg,
		reqLatencyCalc:    reqLatencyCalc,
	}
}

func (m *MetricsMiddleware) AggregateDistance(data types.Distance) (err error) {

	defer func(start time.Time) {
		m.reqLatencyAgg.Observe(float64(time.Since(start).Seconds()))
		m.reqCounterAgg.Inc()
		if err != nil {
			m.errReqCounterAgg.Inc()
		}
	}(time.Now())

	err = m.next.AggregateDistance(data)
	return
}

func (m *MetricsMiddleware) CalculateInvoice(obuID int) (invoice *types.Invoice, err error) {

	defer func(start time.Time) {
		m.reqLatencyCalc.Observe(float64(time.Since(start).Seconds()))
		m.reqCounterCalc.Inc()
		if err != nil {
			m.errReqCounterCalc.Inc()
		}
	}(time.Now())

	invoice, err = m.next.CalculateInvoice(obuID)
	return
}

type LogMiddleWare struct {
	next Aggregator
}

func NewLogMiddleware(next Aggregator) Aggregator {
	return &LogMiddleWare{
		next: next,
	}
}

func (l *LogMiddleWare) AggregateDistance(data types.Distance) (err error) {
	defer func(start time.Time) {
		logrus.WithFields(logrus.Fields{
			"took": time.Since(start),
			"err":  err,
		}).Info("Aggregate distance")
	}(time.Now())
	err = l.next.AggregateDistance(data)
	return
}

func (l *LogMiddleWare) CalculateInvoice(obuID int) (invoice *types.Invoice, err error) {
	defer func(start time.Time) {
		var (
			distance float64
			amount   float64
		)
		if invoice != nil {
			distance = invoice.TotalDistance
			amount = invoice.TotalAmount
		}
		logrus.WithFields(logrus.Fields{
			"took":     time.Since(start),
			"err":      err,
			"obuID":    obuID,
			"amount":   amount,
			"distance": distance,
		}).Info("Calculate Invoice")
	}(time.Now())
	invoice, err = l.next.CalculateInvoice(obuID)
	return
}
