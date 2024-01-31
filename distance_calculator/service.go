package main

import (
	"math"

	"github.com/DenisBytes/GoToll/types"
)

type CalculatorServicer interface {
	CalculateDistance(types.OBUData) (float64, error) 
}

type CalculatorService struct {
	prevPoint [] float64
}

func NewCalculatorService() CalculatorServicer{
	return &CalculatorService{}
}

// calculating dist between prevPoints to lastPoints
func (s *CalculatorService) CalculateDistance(data types.OBUData) (float64, error) {
	distance := 0.0
	if len(s.prevPoint) > 0{
		distance = calcDist(s.prevPoint[0], s.prevPoint[1], data.Lat, data.Long)
	}
	s.prevPoint = []float64{data.Lat, data.Long}
	return distance, nil
}

func calcDist(x1, x2, y1, y2 float64) float64{
	return math.Sqrt(math.Pow(x2-x1, 2) + math.Pow(y2-y1, 2))
}