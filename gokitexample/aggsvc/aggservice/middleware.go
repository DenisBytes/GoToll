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

func (mw *LoggingMiddleware) Aggregate(ctx context.Context, dist types.Distance) error {
	return mw.next.Aggregate(ctx, dist)
}

func (mw *LoggingMiddleware) Calculate(ctx context.Context, obuID int) (*types.Invoice, error) {
	return mw.next.Calculate(ctx, obuID)
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

func (mw *InstrumentationMiddleware) Aggregate(ctx context.Context, dist types.Distance) error {
	return mw.next.Aggregate(ctx, dist)
}

func (mw *InstrumentationMiddleware) Calculate(ctx context.Context, dist int) (*types.Invoice, error) {
	return mw.next.Calculate(ctx, dist)
}
