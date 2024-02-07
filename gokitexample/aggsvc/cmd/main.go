package main

import (
	"net"
	"net/http"
	"os"

	"github.com/DenisBytes/GoToll/gokitexample/aggsvc/aggendpoint"
	"github.com/DenisBytes/GoToll/gokitexample/aggsvc/aggservice"
	"github.com/DenisBytes/GoToll/gokitexample/aggsvc/aggtransport"
	"github.com/go-kit/log"
)

func main() {

	logger := log.NewLogfmtLogger(os.Stderr)
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)
	logger = log.With(logger, "caller", log.DefaultCaller)

	service := aggservice.New(logger)
	endpoints := aggendpoint.New(service, logger)
	httpHandler := aggtransport.NewHTTPHandler(endpoints, logger)

	// The HTTP listener mounts the Go kit HTTP handler we created.
	httpListener, err := net.Listen("tcp", ":3005")
	if err != nil {
		logger.Log("transport", "HTTP", "during", "Listen", "err", err)
		os.Exit(1)
	}
	logger.Log("transport", "HTTP", "addr", ":3005")
	err = http.Serve(httpListener, httpHandler)
	if err != nil {
		panic(err)
	}
}
