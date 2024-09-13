package templates

import "weather-app/weather"

type Weather struct {
	Version         string
	City            string
	Country         string
	TimeZone        string
	Latitude        string
	Longitude       string
	IPAddress       string
	ResolvedAddress string
	Temperature     float64
	Conditions      string
	Description     string
	Address         string
	Icon            string
	Days            []weather.Day
}

func GetWeatherData(
	version string,
	city string,
	country string,
	timeZone string,
	latitude string,
	longitude string,
	ipAddress string,
	resolvedAddress string,
	temperature float64,
	conditions string,
	description string,
	address string,
	icon string,
	days []weather.Day,
) *Weather {
	return &Weather{
		Version:         version,
		City:            city,
		Country:         country,
		TimeZone:        timeZone,
		Latitude:        latitude,
		Longitude:       longitude,
		IPAddress:       ipAddress,
		ResolvedAddress: resolvedAddress,
		Temperature:     temperature,
		Conditions:      conditions,
		Description:     description,
		Address:         address,
		Icon:            icon,
		Days:            days,
	}
}
