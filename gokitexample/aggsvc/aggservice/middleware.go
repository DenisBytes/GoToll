package aggservice

import (
	"context"

	"github.com/DenisBytes/GoToll/types"
)

type Middleware func(Service) Service

type LoggingMiddleware struct {
	next Service
}

func NewLoggingMiddleware() Middleware {
	return func(next Service) Service {
		return &LoggingMiddleware{
			next: next,
		}
	}
}

func (mw *LoggingMiddleware) Aggregate(_ context.Context, dist types.Distance) error {
	return nil
}

func (mw *LoggingMiddleware) Calculate(_ context.Context, obuID int) (*types.Invoice, error) {
	return nil, nil
}

type InstrumentationMiddleware struct {
	next Service
}

func NewInstrumentationMiddleware() Middleware {
	return func(next Service) Service {
		return &InstrumentationMiddleware{
			next: next,
		}
	}
}

func (mw *InstrumentationMiddleware) Aggregate(_ context.Context, dist types.Distance) error {
	return nil
}

func (mw *InstrumentationMiddleware) Calculate(_ context.Context, dist int) (*types.Invoice, error) {
	return nil, nil
}
