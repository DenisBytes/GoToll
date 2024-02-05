package aggendpoint

import (
	"context"
	"time"

	"github.com/DenisBytes/GoToll/gokitexample/aggsvc/aggservice"
	"github.com/DenisBytes/GoToll/types"
	"github.com/go-kit/kit/circuitbreaker"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/ratelimit"
	"github.com/go-kit/log"
	"github.com/sony/gobreaker"
	"golang.org/x/time/rate"
)

type Set struct {
	AggregateEndpoint endpoint.Endpoint
	CalculateEndpoint endpoint.Endpoint
}

func New(svc aggservice.Service, logger log.Logger) Set {
	var aggregateEnpoint endpoint.Endpoint
	{
		aggregateEnpoint = MakeAggregateEndpoint(svc)

		aggregateEnpoint = ratelimit.NewErroringLimiter(rate.NewLimiter(rate.Every(time.Second), 1))(aggregateEnpoint)
		aggregateEnpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{}))(aggregateEnpoint)
		// aggregateEnpoint = LoggingMiddleware(log.With(logger, "method", "Sum"))(sumEndpoint)
		// aggregateEnpoint = InstrumentingMiddleware(duration.With("method", "Sum"))(sumEndpoint)
	}
	var calculateEndpoint endpoint.Endpoint
	{
		calculateEndpoint = MakeCalculateEndpoint(svc)
		calculateEndpoint = ratelimit.NewErroringLimiter(rate.NewLimiter(rate.Limit(1), 100))(calculateEndpoint)
		calculateEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{}))(calculateEndpoint)
		// calculateEndpoint = LoggingMiddleware(log.With(logger, "method", "Concat"))(calculateEndpoint)
		// calculateEndpoint = InstrumentingMiddleware(duration.With("method", "Concat"))(calculateEndpoint)
	}
	return Set{
		AggregateEndpoint: aggregateEnpoint,
		CalculateEndpoint: calculateEndpoint,
	}
}

type AggregateRequest struct {
	Value float64 `json:"value"`
	OBUID int     `json:"obuID"`
	Unix  int64   `json:"unix"`
}

type AggregateResponse struct {
	Err error `json:"err"`
}

func (s Set) Aggregate(ctx context.Context, dist types.Distance) error {
	_, err := s.AggregateEndpoint(ctx, AggregateRequest{
		OBUID: dist.OBUID,
		Value: dist.Value,
		Unix:  dist.Unix,
	})
	return err
}

type CalculateRequest struct {
	OBUID int `json:"obuID"`
}

type CalculateResponse struct {
	OBUID         int     `json:"obuID"`
	TotalDistance float64 `json:"totalDistance"`
	TotalAmount   float64 `json:"totalAmount"`
	Err           error   `json:"err"`
}

func (s Set) Calculate(ctx context.Context, obuID int) (*types.Invoice, error) {
	resp, err := s.CalculateEndpoint(ctx, AggregateRequest{
		OBUID: obuID,
	})
	if err != nil {
		return nil, err
	}
	result := resp.(CalculateResponse)

	return &types.Invoice{
		OBUID:         result.OBUID,
		TotalDistance: result.TotalDistance,
		TotalAmount:   result.TotalAmount,
	}, nil
}

func MakeAggregateEndpoint(s aggservice.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(AggregateRequest)
		err = s.Aggregate(ctx, types.Distance{
			OBUID: req.OBUID,
			Value: req.Value,
			Unix:  req.Unix,
		})
		return AggregateResponse{
			Err: err,
		}, nil
	}
}

func MakeCalculateEndpoint(s aggservice.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(CalculateRequest)
		inv, err := s.Calculate(ctx, req.OBUID)
		return CalculateResponse{
			OBUID:         inv.OBUID,
			TotalDistance: inv.TotalDistance,
			TotalAmount:   inv.TotalAmount,
			Err:           err,
		}, nil
	}
}
