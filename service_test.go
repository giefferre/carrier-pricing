package carrierpricing

import (
	"bytes"
	"errors"
	"log"
	"os"
	"reflect"
	"testing"
)

// TESTS

func TestNewService(t *testing.T) {
	// tests that NewService method returns a valid Service object
	logger := log.New(os.Stdout, "", log.LstdFlags)
	csf := &mockCarrierServiceFinder{}

	expectedService := &Service{
		carrierServiceFinder: csf,
		logger:               logger,
	}

	service := NewService(logger, csf)

	if !reflect.DeepEqual(expectedService, service) {
		t.Fatal("NewService method didn't return the correct Service object")
	}
}

func TestGetBasicQuoteLogs(t *testing.T) {
	// tests that the GetBasicQuote logs the execution
	logDestination := bytes.NewBufferString("")

	logger := log.New(logDestination, "", 0)
	csf := &mockCarrierServiceFinder{}

	service := NewService(logger, csf)

	service.GetBasicQuote(GetBasicQuoteArgs{
		"FROM",
		"TO",
	})

	expectedLogString := "executing GetBasicQuote with args: {FROM TO}\n"
	actualLogString := logDestination.String()

	if expectedLogString != logDestination.String() {
		t.Fatalf("expected log message was not received, had %v", actualLogString)
	}
}

func TestGetBasicQuote(t *testing.T) {
	tests := []struct {
		Arguments      GetBasicQuoteArgs
		ExpectedResult *GetBasicQuoteResponse
		ExpectedError  error
	}{
		// case #1 invalid PickupPostcode argument
		{
			Arguments: GetBasicQuoteArgs{
				PickupPostcode:   "_",
				DeliveryPostcode: "EC2A3LT",
			},
			ExpectedResult: nil,
			ExpectedError:  errors.New("strconv.ParseInt: parsing \"_\": invalid syntax"),
		},
		// case #2 invalid DeliveryPostcode argument
		{
			Arguments: GetBasicQuoteArgs{
				PickupPostcode:   "SW1A1AA",
				DeliveryPostcode: "",
			},
			ExpectedResult: nil,
			ExpectedError:  errors.New("strconv.ParseInt: parsing \"\": invalid syntax"),
		},
		// case #3 valid request, expected result
		{
			Arguments: GetBasicQuoteArgs{
				PickupPostcode:   "SW1A1AA",
				DeliveryPostcode: "EC2A3LT",
			},
			ExpectedResult: &GetBasicQuoteResponse{
				PickupPostcode:   "SW1A1AA",
				DeliveryPostcode: "EC2A3LT",
				Price:            316,
			},
			ExpectedError: nil,
		},
	}

	logger := log.New(os.Stdout, "", log.LstdFlags)
	csf := &mockCarrierServiceFinder{}

	service := NewService(logger, csf)

	for _, tc := range tests {
		result, err := service.GetBasicQuote(tc.Arguments)
		if (tc.ExpectedError != nil && err == nil) ||
			(tc.ExpectedError == nil && err != nil) ||
			(tc.ExpectedError != nil && err != nil && tc.ExpectedError.Error() != err.Error()) {
			t.Fatalf(
				"expected error '%v', received: '%v'\n",
				tc.ExpectedError,
				err,
			)
		}
		if !reflect.DeepEqual(tc.ExpectedResult, result) {
			t.Fatalf(
				"expected result '%v', received: '%v'\n",
				tc.ExpectedResult,
				result,
			)
		}
	}
}

func TestGetQuotesByVehicleLogs(t *testing.T) {
	// tests that the GetQuotesByVehicle logs the execution
	logDestination := bytes.NewBufferString("")

	logger := log.New(logDestination, "", 0)
	csf := &mockCarrierServiceFinder{}

	service := NewService(logger, csf)

	service.GetQuotesByVehicle(GetQuotesByVehicleArgs{
		"FROM",
		"TO",
		"bicycle",
	})

	expectedLogString := "executing GetQuotesByVehicle with args: {FROM TO bicycle}\n"
	actualLogString := logDestination.String()

	if expectedLogString != logDestination.String() {
		t.Fatalf("expected log message was not received, had %v", actualLogString)
	}
}

func TestGetQuotesByVehicle(t *testing.T) {
	tests := []struct {
		Arguments      GetQuotesByVehicleArgs
		ExpectedResult *GetQuotesByVehicleResponse
		ExpectedError  error
	}{
		// case #1 invalid PickupPostcode argument
		{
			Arguments: GetQuotesByVehicleArgs{
				PickupPostcode:   "_",
				DeliveryPostcode: "EC2A3LT",
				Vehicle:          "bicycle",
			},
			ExpectedResult: nil,
			ExpectedError:  errors.New("strconv.ParseInt: parsing \"_\": invalid syntax"),
		},
		// case #2 invalid DeliveryPostcode argument
		{
			Arguments: GetQuotesByVehicleArgs{
				PickupPostcode:   "SW1A1AA",
				DeliveryPostcode: "",
				Vehicle:          "bicycle",
			},
			ExpectedResult: nil,
			ExpectedError:  errors.New("strconv.ParseInt: parsing \"\": invalid syntax"),
		},
		// case #3 invalid Vehicle argument
		{
			Arguments: GetQuotesByVehicleArgs{
				PickupPostcode:   "SW1A1AA",
				DeliveryPostcode: "EC2A3LT",
				Vehicle:          "scooter",
			},
			ExpectedResult: nil,
			ExpectedError:  errors.New("invalid vehicle provided"),
		},
		// case #4 valid request, expected result
		{
			Arguments: GetQuotesByVehicleArgs{
				PickupPostcode:   "SW1A1AA",
				DeliveryPostcode: "EC2A3LT",
				Vehicle:          "bicycle",
			},
			ExpectedResult: &GetQuotesByVehicleResponse{
				PickupPostcode:   "SW1A1AA",
				DeliveryPostcode: "EC2A3LT",
				Vehicle:          "bicycle",
				Price:            348,
			},
			ExpectedError: nil,
		},
	}

	logger := log.New(os.Stdout, "", log.LstdFlags)
	csf := &mockCarrierServiceFinder{}

	service := NewService(logger, csf)

	for _, tc := range tests {
		result, err := service.GetQuotesByVehicle(tc.Arguments)
		if (tc.ExpectedError != nil && err == nil) ||
			(tc.ExpectedError == nil && err != nil) ||
			(tc.ExpectedError != nil && err != nil && tc.ExpectedError.Error() != err.Error()) {
			t.Fatalf(
				"expected error '%v', received: '%v'\n",
				tc.ExpectedError,
				err,
			)
		}
		if !reflect.DeepEqual(tc.ExpectedResult, result) {
			t.Fatalf(
				"expected result '%v', received: '%v'\n",
				tc.ExpectedResult,
				result,
			)
		}
	}
}

func TestGetQuotesByCarrierLogs(t *testing.T) {
	// tests that the GetQuotesByCarrier logs the execution
	logDestination := bytes.NewBufferString("")

	logger := log.New(logDestination, "", 0)
	csf := &mockCarrierServiceFinder{}

	service := NewService(logger, csf)

	service.GetQuotesByCarrier(GetQuotesByCarrierArgs{
		"FROM",
		"TO",
		"small_van",
	})

	expectedLogString := "executing GetQuotesByCarrier with args: {FROM TO small_van}\n"
	actualLogString := logDestination.String()

	if expectedLogString != logDestination.String() {
		t.Fatalf("expected log message was not received, had %v", actualLogString)
	}
}

func TestGetQuotesByCarrier(t *testing.T) {
	tests := []struct {
		Arguments      GetQuotesByCarrierArgs
		ExpectedResult *GetQuotesByCarrierResponse
		ExpectedError  error
	}{
		// case #1 invalid PickupPostcode argument
		{
			Arguments: GetQuotesByCarrierArgs{
				PickupPostcode:   "_",
				DeliveryPostcode: "EC2A3LT",
				Vehicle:          "bicycle",
			},
			ExpectedResult: nil,
			ExpectedError:  errors.New("strconv.ParseInt: parsing \"_\": invalid syntax"),
		},
		// case #2 invalid DeliveryPostcode argument
		{
			Arguments: GetQuotesByCarrierArgs{
				PickupPostcode:   "SW1A1AA",
				DeliveryPostcode: "",
				Vehicle:          "bicycle",
			},
			ExpectedResult: nil,
			ExpectedError:  errors.New("strconv.ParseInt: parsing \"\": invalid syntax"),
		},
		// case #3 invalid Vehicle argument
		{
			Arguments: GetQuotesByCarrierArgs{
				PickupPostcode:   "SW1A1AA",
				DeliveryPostcode: "EC2A3LT",
				Vehicle:          "scooter",
			},
			ExpectedResult: nil,
			ExpectedError:  errors.New("invalid vehicle provided"),
		},
		// case #4 carrier service finder gives 0 results
		{
			Arguments: GetQuotesByCarrierArgs{
				PickupPostcode:   "SW1A1AA",
				DeliveryPostcode: "EC2A3LT",
				Vehicle:          "bicycle", // vehicle is valid but the provider has no carrier services for it
			},
			ExpectedResult: nil,
			ExpectedError:  errors.New("no available carrier services for the given vehicle"),
		},
		// case #5 valid request, expected result
		{
			Arguments: GetQuotesByCarrierArgs{
				PickupPostcode:   "SW1A1AA",
				DeliveryPostcode: "EC2A3LT",
				Vehicle:          "small_van",
			},
			ExpectedResult: &GetQuotesByCarrierResponse{
				PickupPostcode:   "SW1A1AA",
				DeliveryPostcode: "EC2A3LT",
				Vehicle:          "small_van",
				PriceList: PriceByCarrierList{
					PriceByCarrier{
						CarrierName:  "MockService2",
						Amount:       421,
						DeliveryTime: 5,
					},
					PriceByCarrier{
						CarrierName:  "MockService1",
						Amount:       431,
						DeliveryTime: 1,
					},
				},
			},
			ExpectedError: nil,
		},
	}

	logger := log.New(os.Stdout, "", log.LstdFlags)
	csf := &mockCarrierServiceFinder{}

	service := NewService(logger, csf)

	for _, tc := range tests {
		result, err := service.GetQuotesByCarrier(tc.Arguments)
		if (tc.ExpectedError != nil && err == nil) ||
			(tc.ExpectedError == nil && err != nil) ||
			(tc.ExpectedError != nil && err != nil && tc.ExpectedError.Error() != err.Error()) {
			t.Fatalf(
				"expected error '%v', received: '%v'\n",
				tc.ExpectedError,
				err,
			)
		}
		if !reflect.DeepEqual(tc.ExpectedResult, result) {
			t.Fatalf(
				"expected result '%v', received: '%v'\n",
				tc.ExpectedResult,
				result,
			)
		}
	}
}

// UTILS

type mockCarrierServiceFinder struct{}

func (mcsf *mockCarrierServiceFinder) FindCarrierServicesForVehicle(vehicleType string) (availableCarrierServices []CarrierService) {
	switch vehicleType {
	case VehicleTypeSmallVan:
		availableCarrierServices = []CarrierService{
			CarrierService{
				Name:         "MockService1",
				Markup:       20,
				DeliveryTime: 1,
			},
			CarrierService{
				Name:         "MockService2",
				Markup:       10,
				DeliveryTime: 5,
			},
		}
	}
	return
}
