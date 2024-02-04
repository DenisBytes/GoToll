package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/DenisBytes/GoToll/types"
	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}
	var (
		store          = makeStore()
		svc            = NewInvoiceAggregator(store)
		grpcListenAddr = os.Getenv("AGG_GRPC_ENDPOINT")
		httpListenAddr = os.Getenv("AGG_HTTP_ENDPOINT")
	)
	svc = NewMetricsMiddleware(svc)
	svc = NewLogMiddleware(svc)
	go func() {
		log.Fatal(makeGRPCTransport(grpcListenAddr, svc))
	}()
	makeHTTPTransport(httpListenAddr, svc)
}

func makeGRPCTransport(listenAddr string, svc Aggregator) error {
	fmt.Println("GRPC transport running on port: ", listenAddr)
	//Make a tcp listener
	ln, err := net.Listen("tcp", listenAddr)
	if err != nil {
		return err
	}
	// make a new GRPC server with (options)
	server := grpc.NewServer([]grpc.ServerOption{}...)
	// register (OUR) GRPX server implementation to the GRPC package
	types.RegisterAggregatorServer(server, NewGRPCAggregatorServer(svc))

	return server.Serve(ln)
}

func makeHTTPTransport(listenAddr string, svc Aggregator) {

	var (
		aggMetricHandler = NewHTTPMetricHandler("aggregate")
		invMetricHandler = NewHTTPMetricHandler("invoice")
		aggregateHandler = makeHTTPHandlerFunc(aggMetricHandler.instrument(handleAggregate(svc)))
		invoiceHandler   = makeHTTPHandlerFunc(invMetricHandler.instrument(handleGetInvoice(svc)))
	)
	http.HandleFunc("/aggregate", aggregateHandler)
	http.HandleFunc("/invoice", invoiceHandler)
	http.Handle("/metrics", promhttp.Handler())
	fmt.Println("HTTP transport running on port: ", listenAddr)
	log.Fatal(http.ListenAndServe(listenAddr, nil))
}

func makeStore() Storer {
	storeType := os.Getenv("AGG_STORE_TYPE")
	switch storeType {
	case "memory":
		return NewMemoryStore()
	default:
		log.Fatalf("invalid store type given %s", storeType)
		return nil
	}
}

func writeJSON(w http.ResponseWriter, status int, v any) error {
	w.WriteHeader(status)
	w.Header().Add("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(v)
}