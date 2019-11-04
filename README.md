# Carrier Pricing

REST API golang example, used to calculate delivery pricing according to different vehicles/carriers 🚚

The repository is structured to provide an application which actually serves REST API requests, to calculate the price of the parcel delivery service. In the same time, such repository could be imported in other applications as well, as the HTTP server is separated from the service package itself.

## Available APIs

- `/quotes` or `/quotes/basic`: provides users with a basic calculation of the delivery service price between two post codes
- `/quotes/byvehicle`: provides users with a calculation of the delivery service price between two post codes; the price will change according to the specific vehicle the user wants
- `/quotes/bycarrier` (work in progress): provides users with the list of all the prices for a delivery of a parcel using different vehicles and different carriers

REST examples are available in the [docs/examples](docs/examples) folder.
They are meant to be used on [VSCode](https://code.visualstudio.com) [REST Client plugin](https://github.com/Huachao/vscode-restclient).

## Running the application

### Requirements

- Docker engine v. >= 19.03.4
- Linux / Unix machine w/GNU make installed
- Go distribution

### Commands

To start the application, run:

```bash
make start
```

PLEASE NOTE: to provide SSL features, the application will use [Caddy Server](https://caddyserver.com/) using a self-signed certificate. Such certificate will be generated using the [mkcert tool](https://github.com/FiloSottile/mkcert).

## Tests

You can run tests by executing the `make tests` command.

## Implementing Carrier Finder services

A CarrierFinder service is piece of software used from the package for the `GetQuotesByCarrier` method.

Its aim is to provide the carrierpricing a list of Carriers supporting a specific vehicle.

In this repository a simple CarrierFinder is implemented, but you can implement your own.

All you need to do is to implement a struct with the following interface:

```go
    FindCarrierServicesForVehicle(vehicleType string) []carrierpricing.CarrierService
```

### Example

```go
package mockcarrierservicefinder

import (
	"github.com/giefferre/carrierpricing"
)

// MockCarrierServiceFinder implements a CarrierServiceFinder interface,
// returning static data.
type MockCarrierServiceFinder struct{}

// New returns a new MockCarrierServiceFinder object.
func New() *MockCarrierServiceFinder {
	return &MockCarrierServiceFinder{}
}

// FindCarrierServicesForVehicle returns mock, static data.
func (mcsf *MockCarrierServiceFinder) FindCarrierServicesForVehicle(vehicleType string) (availableCarrierServices []carrierpricing.CarrierService) {
	switch vehicleType {
	case carrierpricing.VehicleTypeSmallVan:
		availableCarrierServices = []carrierpricing.CarrierService{
			carrierpricing.CarrierService{
				Name:         "RoyalPackages",
				Markup:       80,
				DeliveryTime: 1,
			},
			carrierpricing.CarrierService{
				Name:         "Hercules",
				Markup:       35,
				DeliveryTime: 5,
			},
			carrierpricing.CarrierService{
				Name:         "CollectTimes",
				Markup:       70,
				DeliveryTime: 1,
			},
		}
	}
	return
}
```