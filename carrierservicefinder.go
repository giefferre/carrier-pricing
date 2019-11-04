package carrierpricing

// CarrierServiceFinder is a software service used to get the list of all the available carriers
// for a specific vehicle.
type CarrierServiceFinder interface {
	FindCarrierServicesForVehicle(vehicleType string) []CarrierService
}

// CarrierService represents the way a company carrying parcels around can deliver
// parcels according to a specific vehicle. It includes the Carrier name, a Markup
// (composed of both base markup and vehicle-based markup) and a DeliveryTime, in
// minutes.
type CarrierService struct {
	Name         string
	Markup       int64
	DeliveryTime int64
}
