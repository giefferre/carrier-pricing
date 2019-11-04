/*
Package carrierpricing allows to calculate quotes for carrier pricing
according to three different methods.

Available methods:

- GetBasicQuote:		returns the price just according to pickup and delivery postcodes
- GetQuoteByVehicle:	returns the price according to pickup and delivery postcodes and the vehicle used;
						available vehicles are: "bicycle", "motorbike", "parcel_car", "small_van", "large_van"
- GetQuoteByCarrier:	returns the price according to pickup and delivery postcodes and the vehicle used,
						giving the list of the all available carriers used; all prices are sorted by price
*/
package carrierpricing
