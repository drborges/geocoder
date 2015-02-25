package google

import (
	"encoding/json"
	"fmt"
	"github.com/drborges/geocoder"
	"io/ioutil"
	"net/http"
)

const (
	ReverseGeocodeEndpoint = "https://maps.googleapis.com/maps/api/geocode/json?latlng=%v,%v"
)

type Geocoder struct {
	HttpClient             *http.Client
	ReverseGeocodeEndpoint string
}

func NewGeocoder() geocoder.Geocoder {
	return &Geocoder{
		HttpClient:             &http.Client{},
		ReverseGeocodeEndpoint: ReverseGeocodeEndpoint,
	}
}

func NewGeocoderWithHttpProvider(c *http.Client) geocoder.Geocoder {
	return &Geocoder{
		HttpClient:             c,
		ReverseGeocodeEndpoint: ReverseGeocodeEndpoint,
	}
}

func (geo *Geocoder) ReverseGeocode(lat float64, lng float64) (*http.Response, error) {
	return geo.HttpClient.Get(fmt.Sprintf(geo.ReverseGeocodeEndpoint, lat, lng))
}

func AddressMapper(res *http.Response) (geocoder.Address, error) {
	if body, err := ioutil.ReadAll(res.Body); err == nil {
		r := new(Response)
		json.Unmarshal(body, &r)
		return r.address(), nil
	} else {
		return geocoder.EmptyAddress, err
	}
}

type QueryResults struct {
	Types             []string            `json:"types"`
	AddressComponents []AddressComponents `json:"address_components"`
}

type AddressComponents struct {
	LongName  string   `json:"long_name"`
	ShortName string   `json:"short_name"`
	Types     []string `json:"types"`
}

type Response struct {
	Results []QueryResults `json:"results"`
}

func (res *Response) address() geocoder.Address {
	address := geocoder.Address{}
	for _, result := range res.Results {
		if result.Types[0] == "postal_code" {
			for _, addrComponent := range result.AddressComponents {
				if addrComponent.Types[0] == "locality" {
					address.City = addrComponent.ShortName
				} else if addrComponent.Types[0] == "administrative_area_level_1" {
					address.State = addrComponent.ShortName
				} else if addrComponent.Types[0] == "country" {
					address.Country = addrComponent.ShortName
				}
			}
		}
	}

	return address
}
