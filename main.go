package main

import (
	"context"
	"embed"
	_ "embed"
	"fmt"
	"github.com/fastly/compute-sdk-go/fsthttp"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"weather-app/foursquare"
	"weather-app/location"
	"weather-app/opencage"
	"weather-app/templates"
	"weather-app/weather"
)

//go:embed templates
var templateFS embed.FS

//go:generate npm run build

func main() {
	fsthttp.ServeFunc(func(ctx context.Context, w fsthttp.ResponseWriter, r *fsthttp.Request) {
		checkCors(w, r)
		filterHTTPMethods(w, r)
		handleHTTPOptions(w, r)

		// STATIC FRONT PAGE ROUTE
		if r.URL.Path == "/" {
			serveFrontPage(ctx, w, r)
			return
		}

		// TIMELINE ROUTE
		if r.URL.Path == "/v2/timeline/" {
			w.Header().Set("Cache-Control", "public, max-age=300")
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			var lat, long string
			var err error
			var resp *fsthttp.Response

			lat = r.URL.Query().Get("lat")
			long = r.URL.Query().Get("lon")

			if lat != "" && long != "" {
				resp, err = weather.GetWeatherRaw(ctx, lat, long)
				if err != nil {
					log.Println("Error: lat and long from backend: Failed getting weather data", err)
				}
				if err != nil {
					log.Println("Error: failed to get opencage reverse geocoding", err)
				}

			} else {
				//get details from embedded geo-reverse db
				details, err := location.GetDetailsFromIp(r.RemoteAddr)
				if err != nil {

					log.Println(w, "Was not able to read from database: ", err)
				}
				lat = details.Lat
				long = details.Long
				resp, err = weather.GetWeatherRaw(ctx, lat, long)
				if err != nil {
					log.Println("Error: lat and long from browser. Failed getting weather data", err)
				}
			}

			io.Copy(w, resp.Body)
			return
		}

		//LOCATION NAME RESOLVER ROUTE
		if r.URL.Path == "/v2/location-name-resolver/" {
			w.Header().Set("Cache-Control", "public, max-age=300")
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			var lat string = r.URL.Query().Get("lat")
			var long string = r.URL.Query().Get("lon")
			resp, err := opencage.GetLocationMetaData(ctx, lat, long)
			if err != nil {
				fmt.Fprintf(w, "Error reading response body: %v\n", err)
			}
			io.Copy(w, resp.Body)
			return
		}

		if r.URL.Path == "/v2/nearby/" {
			w.Header().Set("Cache-Control", "public, max-age=300")
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			var lat string = r.URL.Query().Get("lat")
			var long string = r.URL.Query().Get("lon")
			var locType string = r.URL.Query().Get("type")
			resp, err := foursquare.GetNearby(ctx, lat, long, locType)
			if err != nil {
				fmt.Fprintf(w, fmt.Sprintf("Error fetching nearby beaches: %v", err), http.StatusInternalServerError)
				return
			}
			if resp.Body == nil {
				fmt.Fprintf(w, "Received empty response body", http.StatusInternalServerError)
				return
			}

			_, copyErr := io.Copy(w, resp.Body)
			if copyErr != nil {
				fmt.Fprintf(w, fmt.Sprintf("Error reading response body: %v", copyErr), http.StatusInternalServerError)
			}
			return
		}

		w.WriteHeader(fsthttp.StatusNotFound)
		fmt.Fprintf(w, "The page you requested could not be found\n")
	})
}

func checkCors(w fsthttp.ResponseWriter, r *fsthttp.Request) {
	acceptedOrigins := map[string]bool{
		"http://localhost:3000":                        true,
		"https://weather-app-react-psi-one.vercel.app": true,
		"https://www.localsonly.today":                 true,
		"https://localsonly.today":                     true,
	}

	origin := r.Header.Get("Origin")
	if acceptedOrigins[origin] {
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	}
}

func filterHTTPMethods(w fsthttp.ResponseWriter, r *fsthttp.Request) {
	if r.Method == "POST" || r.Method == "PUT" || r.Method == "PATCH" || r.Method == "DELETE" {
		w.WriteHeader(fsthttp.StatusMethodNotAllowed)
		fmt.Fprintf(w, "This method is not  allowed\n")
		return
	}
}

func handleHTTPOptions(w fsthttp.ResponseWriter, r *fsthttp.Request) {
	if r.Method == fsthttp.MethodOptions {
		w.WriteHeader(fsthttp.StatusOK)
		return
	}
}

func serveFrontPage(ctx context.Context, w fsthttp.ResponseWriter, r *fsthttp.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	r.CacheOptions.TTL = 0

	details, err := location.GetDetailsFromIp(r.RemoteAddr)
	if err != nil {
		fmt.Fprintf(w, "Was not able to read from databse: ", err)
	}

	weatherData, err := weather.GetWeather(ctx, details.Lat, details.Long)
	if err != nil {
		fmt.Fprintf(w, "Error reading weather data: ", err)
	}

	funcMap := template.FuncMap{
		"GetDay": weather.GetDay,
	}

	tmpl, err := template.New("weather.html").Funcs(funcMap).ParseFS(templateFS, "templates/weather.html")
	if err != nil {
		fmt.Fprintf(w, "Error loading template:", err)
	}

	weatherTplData := templates.GetWeatherData(
		os.Getenv("FASTLY_SERVICE_VERSION"),
		details.City,
		details.Country,
		details.TimeZone,
		details.Lat,
		details.Long,
		r.RemoteAddr,
		weatherData.ResolvedAddress,
		weatherData.CurrentConditions.Temperature,
		weatherData.CurrentConditions.Conditions,
		weatherData.Description,
		weatherData.Address,
		weatherData.CurrentConditions.Icon,
		weatherData.Days,
	)

	err = tmpl.Execute(w, weatherTplData)
	if err != nil {
		fmt.Fprintf(w, "Error loading template:", err)
	}

	return
}
