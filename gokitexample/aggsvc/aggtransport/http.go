package aggtransport

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/DenisBytes/GoToll/gokitexample/aggsvc/aggendpoint"
	"github.com/DenisBytes/GoToll/gokitexample/aggsvc/aggservice"
	"github.com/go-kit/kit/circuitbreaker"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/ratelimit"
	"github.com/go-kit/kit/transport"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/go-kit/log"
	"github.com/sony/gobreaker"
	"golang.org/x/time/rate"
)

func NewHTTPClient(instance string, logger log.Logger) (aggservice.Service, error) {
	if !strings.HasPrefix(instance, "http") {
		instance = "http://" + instance
	}
	u, err := url.Parse(instance)
	if err != nil {
		return nil, err
	}

	limiter := ratelimit.NewErroringLimiter(rate.NewLimiter(rate.Every(time.Second), 100))

	var options []httptransport.ClientOption

	var aggregateEndpoint endpoint.Endpoint
	{
		aggregateEndpoint = httptransport.NewClient(
			"POST",
			copyURL(u, "/aggregate"),
			encodeHTTPGenericRequest,
			decodeHTTPAggregateResponse,
			options...,
		).Endpoint()
		aggregateEndpoint = limiter(aggregateEndpoint)
		aggregateEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{
			Name:    "Aggregate",
			Timeout: 30 * time.Second,
		}))(aggregateEndpoint)
	}

	var calculateEnpoint endpoint.Endpoint
	{
		calculateEnpoint = httptransport.NewClient(
			"GET",
			copyURL(u, "/invoice"),
			encodeHTTPGenericRequest,
			decodeHTTPCalculateResponse,
			options...,
		).Endpoint()

		calculateEnpoint = limiter(calculateEnpoint)
		calculateEnpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{
			Name:    "Calculate",
			Timeout: 10 * time.Second,
		}))(calculateEnpoint)
	}

	return aggendpoint.Set{
		AggregateEndpoint: aggregateEndpoint,
		CalculateEndpoint: calculateEnpoint,
	}, nil
}

func copyURL(base *url.URL, path string) *url.URL {
	next := *base
	next.Path = path
	return &next
}

func errorEncoder(ctx context.Context, err error, w http.ResponseWriter) {
}

func NewHTTPHandler(endpoints aggendpoint.Set, logger log.Logger) http.Handler {

	options := []httptransport.ServerOption{
		httptransport.ServerErrorEncoder(errorEncoder),
		httptransport.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
	}

	m := http.NewServeMux()
	m.Handle("/aggregate", httptransport.NewServer(
		endpoints.AggregateEndpoint,
		decodeHTTPAggregateRequest,
		encodeHTTPGenericResponse,
		options...,
	))
	m.Handle("/invoice", httptransport.NewServer(
		endpoints.CalculateEndpoint,
		decodeHTTPCalculateRequest,
		encodeHTTPGenericResponse,
		options...,
	))
	return m
}

func decodeHTTPAggregateRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req aggendpoint.AggregateRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	return req, err
}

func decodeHTTPCalculateRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req aggendpoint.CalculateRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	return req, err
}

func decodeHTTPAggregateResponse(_ context.Context, r *http.Response) (interface{}, error) {
	if r.StatusCode != http.StatusOK {
		return nil, errors.New(r.Status)
	}
	var resp aggendpoint.AggregateResponse
	err := json.NewDecoder(r.Body).Decode(&resp)
	return resp, err
}

func decodeHTTPCalculateResponse(_ context.Context, r *http.Response) (interface{}, error) {
	if r.StatusCode != http.StatusOK {
		return nil, errors.New(r.Status)
	}
	var resp aggendpoint.CalculateResponse
	err := json.NewDecoder(r.Body).Decode(&resp)
	return resp, err
}

func encodeHTTPGenericRequest(_ context.Context, r *http.Request, request interface{}) error {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(request); err != nil {
		return err
	}
	r.Body = ioutil.NopCloser(&buf)
	return nil
}

func encodeHTTPGenericResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if f, ok := response.(endpoint.Failer); ok && f.Failed() != nil {
		errorEncoder(ctx, f.Failed(), w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

// options := []httptransport.ServerOption{
// 	httptransport.ServerErrorEncoder(errorEncoder),
// 	httptransport.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
// }

// if zipkinTracer != nil {
// 	// Zipkin HTTP Server Trace can either be instantiated per endpoint with a
// 	// provided operation name or a global tracing service can be instantiated
// 	// without an operation name and fed to each Go kit endpoint as ServerOption.
// 	// In the latter case, the operation name will be the endpoint's http method.
// 	// We demonstrate a global tracing service here.
// 	options = append(options, zipkin.HTTPServerTrace(zipkinTracer))
// }
