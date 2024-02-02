package main

import (
	"log"

	"github.com/DenisBytes/GoToll/aggregator/client"
)

const (
	kafkaTopic = "obudata"
	aggregatorEndpoint = "http://localhost:3000/aggregate"
)

func main() {
	var (
		err error
		svc CalculatorServicer
	)
	svc = NewCalculatorService()
	svc = NewLogMiddleware(svc)
	kafkaConsumer, err := NewKafkaConsumer(kafkaTopic, svc, client.NewHTTPClient(aggregatorEndpoint))
	if err != nil {
		log.Fatal(err)
	}
	kafkaConsumer.Start()
}
