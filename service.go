package carrierpricing

import (
	"errors"
	"log"
	"math"
	"sort"
	"strconv"
)

const (
	// VehicleTypeBicycle means that a bicycle is used to carry parcels around.
	VehicleTypeBicycle = "bicycle"

	// VehicleTypeMotorbike means that a motorbike is used to carry parcels around.
	VehicleTypeMotorbike = "motorbike"

	// VehicleTypeParcelCar means that a parcel car is used to carry parcels around.
	VehicleTypeParcelCar = "parcel_car"

	// VehicleTypeSmallVan means that a small van is used to carry parcels around.
	VehicleTypeSmallVan = "small_van"

	// VehicleTypeLargeVan means that a large van is used to carry parcels around.
	VehicleTypeLargeVan = "large_van"
)

// ValidVehicleTypes is the list of all the available vehicles
var ValidVehicleTypes = []string{
	VehicleTypeBicycle,
	VehicleTypeMotorbike,
	VehicleTypeParcelCar,
	VehicleTypeSmallVan,
	VehicleTypeLargeVan,
}

// VehiclesMarkupTable indicates the markup to be applied to the base price for each vehicle type.
var VehiclesMarkupTable = map[string]float64{
	VehicleTypeBicycle:   1.1,
	VehicleTypeMotorbike: 1.15,
	VehicleTypeParcelCar: 1.2,
	VehicleTypeSmallVan:  1.3,
	VehicleTypeLargeVan:  1.4,
}

var (
	errInvalidVehicle                       = errors.New("invalid vehicle provided")
	errMarkupNotPresent                     = errors.New("markup not present")
	errNoAvailableCarrierServicesForVehicle = errors.New("no available carrier services for the given vehicle")
)

// GetBasicQuoteArgs contains arguments for the GetBasicQuote method.
type GetBasicQuoteArgs struct {
	PickupPostcode   string `json:"pickup_postcode"`
	DeliveryPostcode string `json:"delivery_postcode"`
}

// GetBasicQuoteResponse is the response object for the GetBasicQuote method.
type GetBasicQuoteResponse struct {
	GetBasicQuoteArgs
	Price int64 `json:"price"`
}

// GetQuotesByVehicleArgs contains arguments for the GetQuotesByVehicle method.
type GetQuotesByVehicleArgs struct {
	GetBasicQuoteArgs
	Vehicle string `json:"vehicle"`
}

// GetQuotesByVehicleResponse is the response object for the GetQuotesByVehicle method.
type GetQuotesByVehicleResponse struct {
	GetQuotesByVehicleArgs
	Price int64 `json:"price"`
}

// GetQuotesByCarrierArgs contains arguments for the GetQuotesByCarrier method.
type GetQuotesByCarrierArgs GetQuotesByVehicleArgs

// GetQuotesByCarrierResponse is the response object for the GetQuotesByCarrier method.
type GetQuotesByCarrierResponse struct {
	GetQuotesByCarrierArgs
	PriceList PriceByCarrierList `json:"price_list"`
}

// PriceByCarrier is the object returned in the GetQuotesByCarrierResponse
// indicating the service price and delivery time for a specific carrier
// matching the request.
type PriceByCarrier struct {
	CarrierName  string `json:"service"`
	Amount       int64  `json:"price"`
	DeliveryTime int64  `json:"delivery_time"`
}

// PriceByCarrierList is a list of PriceByCarrier.
// This struct has been created to apply sorting convenience methods.
type PriceByCarrierList []PriceByCarrier

func (pbcl PriceByCarrierList) Len() int           { return len(pbcl) }
func (pbcl PriceByCarrierList) Swap(i, j int)      { pbcl[i], pbcl[j] = pbcl[j], pbcl[i] }
func (pbcl PriceByCarrierList) Less(i, j int) bool { return pbcl[i].Amount < pbcl[j].Amount }

// ServiceInterface defines the interface of the Service.
// This is meant to be used from main/external packages, allowing to mock the service itself.
type ServiceInterface interface {
	GetBasicQuote(args GetBasicQuoteArgs) (*GetBasicQuoteResponse, error)
	GetQuotesByVehicle(args GetQuotesByVehicleArgs) (*GetQuotesByVehicleResponse, error)
	GetQuotesByCarrier(args GetQuotesByCarrierArgs) (*GetQuotesByCarrierResponse, error)
}

// Service implements the ServiceInterface exposing the required methods.
type Service struct {
	carrierServiceFinder CarrierServiceFinder
	logger               *log.Logger
}

// NewService returns a new Service initialized with the given parameters.
func NewService(logger *log.Logger, carrierServiceFinder CarrierServiceFinder) *Service {
	return &Service{
		carrierServiceFinder: carrierServiceFinder,
		logger:               logger,
	}
}

// GetBasicQuote calculates the basic price of the delivery between pickup and delivery
// post codes, provided via the given args.
func (s *Service) GetBasicQuote(args GetBasicQuoteArgs) (*GetBasicQuoteResponse, error) {
	s.logger.Printf("executing GetBasicQuote with args: %v\n", args)

	basePrice, err := s.calculateBasePrice(args.PickupPostcode, args.DeliveryPostcode)
	if err != nil {
		return nil, err
	}

	return &GetBasicQuoteResponse{
		GetBasicQuoteArgs: GetBasicQuoteArgs{
			PickupPostcode:   args.PickupPostcode,
			DeliveryPostcode: args.DeliveryPostcode,
		},
		Price: *basePrice,
	}, nil
}

// GetQuotesByVehicle calculates the price of the delivery betweeen pickup and delivery
// post codes according to a specific vehicle, multiplying the basic price for the
// relative markup.
func (s *Service) GetQuotesByVehicle(args GetQuotesByVehicleArgs) (*GetQuotesByVehicleResponse, error) {
	s.logger.Printf("executing GetQuotesByVehicle with args: %v\n", args)

	if !s.isVehicleValid(args.Vehicle) {
		return nil, errInvalidVehicle
	}

	basePrice, err := s.calculateBasePrice(args.PickupPostcode, args.DeliveryPostcode)
	if err != nil {
		return nil, err
	}

	priceByVehicle := s.applyVehicleMarkup(*basePrice, args.Vehicle)

	return &GetQuotesByVehicleResponse{
		GetQuotesByVehicleArgs: GetQuotesByVehicleArgs{
			GetBasicQuoteArgs: GetBasicQuoteArgs{
				PickupPostcode:   args.PickupPostcode,
				DeliveryPostcode: args.DeliveryPostcode,
			},
			Vehicle: args.Vehicle,
		},
		Price: priceByVehicle,
	}, nil
}

// GetQuotesByCarrier calculates the price of the delivery between pickup and delivery
// post codes according to a specified vehicle and all the available carriers. A
// markup is applied to the basic price for both the vehicle type and the carriers.
func (s *Service) GetQuotesByCarrier(args GetQuotesByCarrierArgs) (*GetQuotesByCarrierResponse, error) {
	s.logger.Printf("executing GetQuotesByCarrier with args: %v\n", args)

	if !s.isVehicleValid(args.Vehicle) {
		return nil, errInvalidVehicle
	}

	basePrice, err := s.calculateBasePrice(args.PickupPostcode, args.DeliveryPostcode)
	if err != nil {
		return nil, err
	}

	priceByVehicle := s.applyVehicleMarkup(*basePrice, args.Vehicle)

	availableCarrierServices := s.carrierServiceFinder.FindCarrierServicesForVehicle(args.Vehicle)
	if len(availableCarrierServices) == 0 {
		return nil, errNoAvailableCarrierServicesForVehicle
	}

	priceList := s.getPriceListFromPriceAndCarrierServices(priceByVehicle, availableCarrierServices)

	return &GetQuotesByCarrierResponse{
		GetQuotesByCarrierArgs: GetQuotesByCarrierArgs{
			GetBasicQuoteArgs: GetBasicQuoteArgs{
				PickupPostcode:   args.PickupPostcode,
				DeliveryPostcode: args.DeliveryPostcode,
			},
			Vehicle: args.Vehicle,
		},
		PriceList: priceList,
	}, nil
}

func (s *Service) calculateBasePrice(pickupPostcode, deliveryPostcode string) (*int64, error) {
	pickup, err := strconv.ParseInt(pickupPostcode, 36, 64)
	if err != nil {
		return nil, err
	}

	delivery, err := strconv.ParseInt(deliveryPostcode, 36, 64)
	if err != nil {
		return nil, err
	}

	const someLargeNumber = 100000000
	result := int64(math.Abs(float64(pickup)-float64(delivery))) / someLargeNumber

	return &result, nil
}

func (s *Service) isVehicleValid(vehicleLabelToVerify string) bool {
	for _, vehicleType := range ValidVehicleTypes {
		if vehicleType == vehicleLabelToVerify {
			return true
		}
	}
	return false
}

func (s *Service) applyVehicleMarkup(basePrice int64, vehicleType string) int64 {
	markup, exists := VehiclesMarkupTable[vehicleType]
	if !exists {
		return basePrice
	}

	return int64(math.RoundToEven(float64(basePrice) * markup))
}

func (s *Service) getPriceListFromPriceAndCarrierServices(priceByVehicle int64, availableCarrierServices []CarrierService) PriceByCarrierList {
	priceList := PriceByCarrierList{}
	for _, carrierService := range availableCarrierServices {
		priceList = append(priceList, PriceByCarrier{
			CarrierName:  carrierService.Name,
			Amount:       priceByVehicle + carrierService.Markup,
			DeliveryTime: carrierService.DeliveryTime,
		})
	}

	sort.Sort(priceList)

	return priceList
}
