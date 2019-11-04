package main

import (
	"log"
	"os"

	"github.com/giefferre/carrierpricing"
	"github.com/giefferre/carrierpricing/carrierservicefinders"
	"github.com/giefferre/carrierpricing/internal/httpserver"
)

var (
	logger               *log.Logger
	carrierServiceFinder carrierpricing.CarrierServiceFinder
)

func init() {
	// logger is initialized to print any information on standard out
	logger = log.New(os.Stdout, "carrierpricing ", log.LstdFlags)

	// here we configure the carrierservicefinder;
	// in this case we want to use a CSFFromJSONFile object, passing the file path
	// of the source data via CSF_JSON_FILE environment variable.
	var err error
	jsonFilePath := os.Getenv("CSF_JSON_FILE")

	logger.Printf("Trying to use CSFFromJSONFile with file: %s", jsonFilePath)
	carrierServiceFinder, err = carrierservicefinders.NewCSFFromJSONFile(jsonFilePath)
	if err != nil {
		logger.Fatalf("NewCSFFromJSONFile method returned error %v", err)
	}

	// want to use a simple carrierServiceFinder?
	// comment lines 21:31 and uncomment the following one
	// carrierServiceFinder = carrierservicefinders.NewCSFFromStaticData()
}

func main() {
	carrierPricingService := carrierpricing.NewService(logger, carrierServiceFinder)
	httpServer := httpserver.NewHTTPServer(logger, carrierPricingService)

	httpServer.Start()
}
