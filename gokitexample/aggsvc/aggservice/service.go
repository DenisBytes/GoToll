package aggservice

import (
	"context"

	"github.com/DenisBytes/GoToll/types"
	"github.com/go-kit/log"
)

const basePrice = 3.15

type Service interface {
	Aggregate(context.Context, types.Distance) error
	Calculate(context.Context, int) (*types.Invoice, error)
}

type Storer interface {
	Insert(types.Distance) error
	Get(int) (float64, error)
}

type BasicService struct {
	store Storer
}

func NewBasicService(store Storer) Service {
	return &BasicService{
		store: store,
	}
}

func (svc *BasicService) Aggregate(_ context.Context, dist types.Distance) error {
	return svc.store.Insert(dist)
}

func (svc *BasicService) Calculate(_ context.Context, obuID int) (*types.Invoice, error) {
	dist, err := svc.store.Get(obuID)
	if err != nil {
		return nil, err
	}
	inv := &types.Invoice{
		OBUID:         obuID,
		TotalDistance: dist,
		TotalAmount:   basePrice * dist,
	}
	return inv, nil
}

// New will construct a complete microservice
func New(logger log.Logger) Service {
	var svc Service
	{
		svc = NewBasicService(NewMemoryStore())
		svc = NewLoggingMiddleware(logger)(svc)
		svc = NewInstrumentationMiddleware()(svc)
	}
	return svc
}
