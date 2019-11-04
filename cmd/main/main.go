package main

import (
	"log"
	"os"

	"github.com/giefferre/carrierpricing"
	"github.com/giefferre/carrierpricing/internal/httpserver"
)

var (
	logger *log.Logger
)

func init() {
	// logger is initialized to print any information on standard out
	logger = log.New(os.Stdout, "carrierpricing ", log.LstdFlags)
}

func main() {
	carrierPricingService := carrierpricing.NewService(logger)
	httpServer := httpserver.NewHTTPServer(logger, carrierPricingService)

	httpServer.Start()
}
