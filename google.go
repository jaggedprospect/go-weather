package main

import (
	"encoding/json"
	"fmt"
	"net/url"
)

type Coordinates struct {
	Lat float64
	Lng float64
}

type GoogleGeocodeResult struct {
	Geometry struct {
		Location Coordinates
	}
}

type GoogleGeocodeResponse struct {
	Status  string
	Results []GoogleGeocodeResult
}

// getCoordForPlace gets the latitude and longitude values of the passed address
// using the Google Maps Geocode API.
func getCoordForPlace(place string) (Coordinates, error) {
	escPlace := url.QueryEscape(place)

	url := fmt.Sprintf("https://maps.googleapis.com/maps/api/geocode/json?key=%s&address=%s", keys.GoogleApiKey, escPlace)

	r, err := httpClient.Get(url)
	if err != nil {
		return Coordinates{}, err
	}

	// defers call until the function completes
	defer r.Body.Close()

	var geocode GoogleGeocodeResponse

	err = json.NewDecoder(r.Body).Decode(&geocode)
	if err != nil {
		return Coordinates{}, err
	}

	if geocode.Status != "OK" || len(geocode.Results) < 1 {
		return Coordinates{}, err
	}

	return geocode.Results[0].Geometry.Location, nil
}
