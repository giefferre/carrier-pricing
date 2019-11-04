package carrierservicefinders

import (
	"encoding/json"
	"io/ioutil"

	"github.com/giefferre/carrierpricing"
)

// CSFFromJSONFile implements the carrierpricing.CarrierService interface;
// the source of data is a single JSON encoded file from local storage.
type CSFFromJSONFile struct {
	carriers []carrier
}

// NewCSFFromJSONFile returns a fresh CSFFromJSONFile object having the list of available
// carriers loaded in memory. An error is returned if the file is not found or
// it does not contain valid objects.
func NewCSFFromJSONFile(jsonFilePath string) (*CSFFromJSONFile, error) {
	jsonFileContent, err := ioutil.ReadFile(jsonFilePath)
	if err != nil {
		return nil, err
	}

	carriers := []carrier{}
	err = json.Unmarshal(jsonFileContent, &carriers)
	if err != nil {
		return nil, err
	}

	return &CSFFromJSONFile{
		carriers: carriers,
	}, nil
}

// FindCarrierServicesForVehicle finds CarrierService objects for the given vehicleType.
func (csf *CSFFromJSONFile) FindCarrierServicesForVehicle(vehicleType string) []carrierpricing.CarrierService {
	carrierServices := []carrierpricing.CarrierService{}

	// this surely is NOT the most efficient way to store this data,
	// but it is good enough for the purpose of this simple application.
	for _, carrier := range csf.carriers {
		for _, service := range carrier.Services {
			for _, vehicle := range service.Vehicles {
				if vehicleType == vehicle {
					carrierServices = append(carrierServices, carrierpricing.CarrierService{
						Name:         carrier.Name,
						Markup:       carrier.BasePrice + service.Markup,
						DeliveryTime: service.DeliveryTime,
					})
				}
			}
		}
	}

	return carrierServices
}

type carrier struct {
	Name      string    `json:"carrier_name"`
	BasePrice int64     `json:"base_price"`
	Services  []service `json:"services"`
}

type service struct {
	DeliveryTime int64    `json:"delivery_time"`
	Markup       int64    `json:"markup"`
	Vehicles     []string `json:"vehicles"`
}
