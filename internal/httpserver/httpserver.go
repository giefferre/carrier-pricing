package httpserver

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/giefferre/carrierpricing"
)

const (
	errMessageInternalServerError = "an internal server error occurred, try again later"
)

// HTTPServer implements an HTTP REST API server
type HTTPServer struct {
	logger  *log.Logger
	service carrierpricing.ServiceInterface
}

// NewHTTPServer returns a new HTTPServer object with the given parameters:
// - logger, which is a standard Go *log.Logger object, used to log errors and other info
// - service, which actually implements the carrierpricing Service
func NewHTTPServer(logger *log.Logger, service carrierpricing.ServiceInterface) *HTTPServer {
	return &HTTPServer{
		logger:  logger,
		service: service,
	}
}

// Start setups the HTTP router and starts the HTTP server on port 80
func (s *HTTPServer) Start() {
	http.HandleFunc("/quotes/byvehicle", s.getQuotesByVehicleHandler)
	http.HandleFunc("/quotes/bycarrier", s.getQuotesByCarrierHandler)
	http.HandleFunc("/quotes/basic", s.getBasicQuotesHandler)
	http.HandleFunc("/quotes", s.getBasicQuotesHandler)

	s.logger.Println("Starting HTTP server...")
	err := http.ListenAndServe(":80", nil)
	if err != nil {
		s.logger.Fatal(err)
	}
}

func (s *HTTPServer) getBasicQuotesHandler(w http.ResponseWriter, r *http.Request) {
	// decode the request into a GetBasicQuoteArgs object
	requestObject := &carrierpricing.GetBasicQuoteArgs{}

	err := decodeRequestBodyAsRequestObject(r.Body, &requestObject)
	if err != nil {
		s.logger.Println(err)
		http.Error(w, errMessageInternalServerError, http.StatusInternalServerError)
		return
	}

	responseObject, err := s.service.GetBasicQuote(*requestObject)
	if err != nil {
		s.logger.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	writeResponse(w, responseObject)
}

func (s *HTTPServer) getQuotesByVehicleHandler(w http.ResponseWriter, r *http.Request) {
	// decode the request into a GetQuotesByVehicleArgs object
	requestObject := &carrierpricing.GetQuotesByVehicleArgs{}

	err := decodeRequestBodyAsRequestObject(r.Body, &requestObject)
	if err != nil {
		s.logger.Println(err)
		http.Error(w, errMessageInternalServerError, http.StatusInternalServerError)
		return
	}

	responseObject, err := s.service.GetQuotesByVehicle(*requestObject)
	if err != nil {
		s.logger.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	writeResponse(w, responseObject)
}

func (s *HTTPServer) getQuotesByCarrierHandler(w http.ResponseWriter, r *http.Request) {
}

// decodeRequestBodyAsRequestObject is a utility method which abstracts the way an
// HTTP request body is decoded into the given requestObject passed as argument
func decodeRequestBodyAsRequestObject(requestBody io.ReadCloser, requestObject interface{}) error {
	// always remember to close the requestBody
	defer requestBody.Close()

	// try to read request body as bytes
	requestBodyAsBytes, err := ioutil.ReadAll(requestBody)
	if err != nil {
		return err
	}

	// try to unmarshal the request bytes into the provided request object
	return json.Unmarshal(requestBodyAsBytes, requestObject)
}

// writeResponse is a utility method which encodes the responseObject into a
// JSON object via a http.ResponseWriter object.
func writeResponse(w http.ResponseWriter, responseObject interface{}) {
	responseDataAsBytes, err := json.Marshal(responseObject)
	if err != nil {
		http.Error(w, errMessageInternalServerError, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, "%s", string(responseDataAsBytes))
}
