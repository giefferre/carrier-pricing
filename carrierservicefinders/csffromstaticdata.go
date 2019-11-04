package carrierservicefinders

import (
	"github.com/giefferre/carrierpricing"
)

// CSFFromStaticData implements a CarrierServiceFinder interface,
// returning static data.
type CSFFromStaticData struct{}

// NewCSFFromStaticData returns a new CSFFromStaticData object.
func NewCSFFromStaticData() *CSFFromStaticData {
	return &CSFFromStaticData{}
}

// FindCarrierServicesForVehicle returns mock, static data.
func (csf *CSFFromStaticData) FindCarrierServicesForVehicle(vehicleType string) (availableCarrierServices []carrierpricing.CarrierService) {
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
		break
	}
	return
}
