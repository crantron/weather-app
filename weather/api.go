package weather

import (
	"compute-starter-kit-go/fastly"
	"context"
	"encoding/json"
	"fmt"
	"github.com/fastly/compute-sdk-go/fsthttp"
	"time"
)

const API_BASE_URL string = "https://weather.visualcrossing.com/VisualCrossingWebServices/rest/services/timeline/"

type WeatherResponse struct {
	ResolvedAddress   string            `json:"resolvedAddress"`
	CurrentConditions CurrentConditions `json:"currentConditions"`
	Days              []Day             `json:"days"`
}

type CurrentConditions struct {
	Temperature float64 `json:"temp"`
	Conditions  string  `json:"conditions"`
	WindSpeed   float64 `json:"wspd"`
	Humidity    float64 `json:"humidity"`
}

type Day struct {
	Datetime       string  `json:"datetime"`
	TemperatureMax float64 `json:"tempmax"`
	TemperatureMin float64 `json:"tempmin"`
	Conditions     string  `json:"conditions"`
}

func GetWeather(ctx context.Context, lat string, long string) (*WeatherResponse, error) {
	key, err := fastly.GetSecretStoreKey("keys", "VISUAL_CROSSING_WEATHER")
	if err != nil {
		return nil, err
	}

	url := API_BASE_URL + lat + "," + long + "?key=" + key

	req, err := fsthttp.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := req.Send(ctx, "visualcrossing2")
	var weather WeatherResponse
	if err := json.NewDecoder(resp.Body).Decode(&weather); err != nil {
		return nil, fmt.Errorf("error decoding weather response: %s", err)
	}

	return &weather, nil
}

func GetDay(inputDate string) (string, error) {
	layout := "2006-01-02"
	date, err := time.Parse(layout, inputDate)
	if err != nil {
		return "", err
	}

	return date.Weekday().String(), nil
}

func GetImgNameMap(condition string) (string, error) {
	m := make(map[string]string)
	m["Clear"] = "Clear"
	m["Partially cloudy"] = "partiallycloudy"

	return m[condition], nil
}
