package aggservice

import (
	"context"
	"time"

	"github.com/DenisBytes/GoToll/types"
	"github.com/go-kit/log"
)

type Middleware func(Service) Service

type LoggingMiddleware struct {
	log  log.Logger
	next Service
}

func NewLoggingMiddleware(logger log.Logger) Middleware {
	return func(next Service) Service {
		return &LoggingMiddleware{
			log:  logger,
			next: next,
		}
	}
}

func (mw *LoggingMiddleware) Aggregate(ctx context.Context, dist types.Distance) (err error) {
	defer func(start time.Time) {
		mw.log.Log("method", "Aggregate", "took", time.Since(start), "obu", dist.OBUID, "dist", dist.Value, "err", err)
	}(time.Now())
	err = mw.next.Aggregate(ctx, dist)
	return mw.next.Aggregate(ctx, dist)
}

func (mw *LoggingMiddleware) Calculate(ctx context.Context, obuID int) (inv *types.Invoice, err error) {
	defer func(start time.Time) {
		mw.log.Log("method", "Calculate", "took", time.Since(start), "obuID", obuID, "inv", inv, "err", err)
	}(time.Now())
	inv, err = mw.next.Calculate(ctx, obuID)
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
