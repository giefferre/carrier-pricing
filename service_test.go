package carrierpricing

import (
	"bytes"
	"errors"
	"log"
	"os"
	"reflect"
	"testing"
)

func TestNewService(t *testing.T) {
	// tests that NewService method returns a valid Service object
	logger := log.New(os.Stdout, "", log.LstdFlags)

	expectedService := &Service{
		logger: logger,
	}

	service := NewService(logger)

	if !reflect.DeepEqual(expectedService, service) {
		t.Fatal("NewService method didn't return the correct Service object")
	}
}

func TestGetBasicQuoteLogs(t *testing.T) {
	// tests that the GetBasicQuote logs the execution
	logDestination := bytes.NewBufferString("")
	logger := log.New(logDestination, "", 0)

	service := NewService(logger)

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
				GetBasicQuoteArgs: GetBasicQuoteArgs{
					PickupPostcode:   "SW1A1AA",
					DeliveryPostcode: "EC2A3LT",
				},
				Price: 316,
			},
			ExpectedError: nil,
		},
	}

	logger := log.New(os.Stdout, "", log.LstdFlags)
	service := NewService(logger)

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

	service := NewService(logger)

	service.GetQuotesByVehicle(GetQuotesByVehicleArgs{
		GetBasicQuoteArgs{
			"FROM",
			"TO",
		},
		"bicycle",
	})

	expectedLogString := "executing GetQuotesByVehicle with args: {{FROM TO} bicycle}\n"
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
				GetBasicQuoteArgs{
					PickupPostcode:   "_",
					DeliveryPostcode: "EC2A3LT",
				},
				"bicycle",
			},
			ExpectedResult: nil,
			ExpectedError:  errors.New("strconv.ParseInt: parsing \"_\": invalid syntax"),
		},
		// case #2 invalid DeliveryPostcode argument
		{
			Arguments: GetQuotesByVehicleArgs{
				GetBasicQuoteArgs{
					PickupPostcode:   "SW1A1AA",
					DeliveryPostcode: "",
				},
				"bicycle",
			},
			ExpectedResult: nil,
			ExpectedError:  errors.New("strconv.ParseInt: parsing \"\": invalid syntax"),
		},
		// case #3 invalid Vehicle argument
		{
			Arguments: GetQuotesByVehicleArgs{
				GetBasicQuoteArgs{
					PickupPostcode:   "SW1A1AA",
					DeliveryPostcode: "EC2A3LT",
				},
				"scooter",
			},
			ExpectedResult: nil,
			ExpectedError:  errors.New("invalid vehicle provided"),
		},
		// case #4 valid request, expected result
		{
			Arguments: GetQuotesByVehicleArgs{
				GetBasicQuoteArgs{
					PickupPostcode:   "SW1A1AA",
					DeliveryPostcode: "EC2A3LT",
				},
				"bicycle",
			},
			ExpectedResult: &GetQuotesByVehicleResponse{
				GetQuotesByVehicleArgs: GetQuotesByVehicleArgs{
					GetBasicQuoteArgs: GetBasicQuoteArgs{
						PickupPostcode:   "SW1A1AA",
						DeliveryPostcode: "EC2A3LT",
					},
					Vehicle: "bicycle",
				},
				Price: 348,
			},
			ExpectedError: nil,
		},
	}

	logger := log.New(os.Stdout, "", log.LstdFlags)
	service := NewService(logger)

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
