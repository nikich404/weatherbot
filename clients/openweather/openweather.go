package openweather

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type OpenWeatherClient struct {
	apiKey string
}

func New(apiKey string) *OpenWeatherClient {
	return &OpenWeatherClient{
		apiKey: apiKey,
	}
}

func (o OpenWeatherClient) Coordinates(city string) (Coordinates, error) {
	url := "http://api.openweathermap.org/geo/1.0/direct?q=%s&limit=5&appid=%s"
	resp, err := http.Get(fmt.Sprintf(url, city, o.apiKey))
	if err != nil {
		return Coordinates{}, err
	}

	if resp.StatusCode != 200 {
		return Coordinates{}, fmt.Errorf("API error: status code %d", resp.StatusCode)
	}

	var coordinatesResponse []CordinatesResponse
	err = json.NewDecoder(resp.Body).Decode(&coordinatesResponse)
	if err != nil {
		return Coordinates{}, err
	}
	if len(coordinatesResponse) == 0 {
		return Coordinates{}, fmt.Errorf("city '%s' not found", city)
	}
	return Coordinates{
		Lat: coordinatesResponse[0].Lat,
		Lon: coordinatesResponse[0].Lon,
	}, nil
}

func (o OpenWeatherClient) Weather(lat, lon float64) (Weather, error) {

	url := fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?lat=%f&lon=%f&appid=%s&units=metric", lat, lon, o.apiKey)

	resp, err := http.Get(url)
	if err != nil {
		return Weather{}, fmt.Errorf("failed to get weather: %w", err)
	}

	var weatherResponse WeatherResponse
	err = json.NewDecoder(resp.Body).Decode(&weatherResponse)
	if err != nil {
		return Weather{}, fmt.Errorf("failed to decode response: %w", err)
	}

	return Weather{
		Temp: weatherResponse.Main.Temp,
	}, nil
}
