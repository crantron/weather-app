package location

import (
	_ "embed"
	"github.com/oschwald/geoip2-golang"
	"net"
	"strconv"
)

// Fastly FS workaround
//
//go:embed assets/GeoLite2-City.mmdb
var geoLite2CityMMDB []byte

type Details struct {
	City     string
	Country  string
	TimeZone string
	Lat      string
	Long     string
}

func GetDetails(
	city string,
	country string,
	timezone string,
	lat float64,
	long float64,
) *Details {
	return &Details{
		City:     city,
		Country:  country,
		TimeZone: timezone,
		Lat:      strconv.FormatFloat(lat, 'f', -1, 64),
		Long:     strconv.FormatFloat(long, 'f', -1, 64),
	}
}

func GetDetailsFromIp(ip string) (*Details, error) {
	db, err := geoip2.FromBytes(geoLite2CityMMDB)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	parseIP := net.ParseIP(ip)
	if parseIP == nil {
		return nil, err
	}

	record, err := db.City(parseIP)
	if err != nil {
		return nil, err
	}

	return GetDetails(
		record.City.Names["en"],
		record.Country.Names["en"],
		record.Location.TimeZone,
		record.Location.Latitude,
		record.Location.Longitude,
	), nil
}
