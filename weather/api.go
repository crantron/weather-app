package weather

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/fastly/compute-sdk-go/fsthttp"
	"github.com/fastly/compute-sdk-go/rtlog"
	"time"
	"weather-app/fastly"
)

const API_BASE_URL string = "https://weather.visualcrossing.com/VisualCrossingWebServices/rest/services/timeline/"

type WeatherResponse struct {
	ResolvedAddress   string            `json:"resolvedAddress"`
	Description       string            `json:"description"`
	Address           string            `json:"address"`
	CurrentConditions CurrentConditions `json:"currentConditions"`
	Days              []Day             `json:"days"`
}

type CurrentConditions struct {
	Temperature float64 `json:"temp"`
	Conditions  string  `json:"conditions"`
	WindSpeed   float64 `json:"wspd"`
	Humidity    float64 `json:"humidity"`
	Icon        string  `json:"icon"`
}

type Day struct {
	Datetime       string  `json:"datetime"`
	TemperatureMax float64 `json:"tempmax"`
	TemperatureMin float64 `json:"tempmin"`
	Conditions     string  `json:"conditions"`
	Icon           string  `json:"icon"`
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

func GetWeatherRaw(ctx context.Context, lat string, long string) (*fsthttp.Response, error) {
	var err error
	s3logging := rtlog.Open("s3-logging")
	key, err := fastly.GetSecretStoreKey("keys", "VISUAL_CROSSING_WEATHER")
	if err != nil {
		return nil, err
	}

	retries := 3
	url := API_BASE_URL + lat + "," + long + "?key=" + key

	for i := 0; i < retries; i++ {
		req, err := fsthttp.NewRequest("GET", url, nil)
		if err != nil {
			return nil, err
		}
		resp, err := req.Send(ctx, "visualcrossing2")
		if err != nil {
			fmt.Fprintln(s3logging, "write %v", err.Error())
			return nil, err
		}
		if resp.StatusCode == fsthttp.StatusOK {
			defer resp.Body.Close()
			return resp, nil
		}
		defer resp.Body.Close()
		time.Sleep(2 * time.Second)
	}

	return nil, fmt.Errorf("error fetching response: %s", err)
}

func GetDay(inputDate string) (string, error) {
	layout := "2006-01-02"
	date, err := time.Parse(layout, inputDate)
	if err != nil {
		return "", err
	}

	return date.Weekday().String(), nil
}
