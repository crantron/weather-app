package main

import (
	"compute-starter-kit-go/location"
	"compute-starter-kit-go/templates"
	"compute-starter-kit-go/weather"
	"context"
	"embed"
	_ "embed"
	"fmt"
	"github.com/fastly/compute-sdk-go/fsthttp"
	"html/template"
	"os"
)

//go:embed templates
var templateFS embed.FS

//go:generate npm run build

func main() {
	fsthttp.ServeFunc(func(ctx context.Context, w fsthttp.ResponseWriter, r *fsthttp.Request) {
		// Filter requests that have unexpected methods.
		if r.Method == "POST" || r.Method == "PUT" || r.Method == "PATCH" || r.Method == "DELETE" {
			w.WriteHeader(fsthttp.StatusMethodNotAllowed)
			fmt.Fprintf(w, "This method is not allowed\n")
			return
		}

		if r.URL.Path == "/" {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			r.CacheOptions.TTL = 5

			details, err := location.GetDetailsFromIp(r.RemoteAddr)
			if err != nil {
				fmt.Fprintf(w, "Was not able to read from databse: ", err)
			}

			weatherData, err := weather.GetWeather(ctx, details.Lat, details.Long)
			if err != nil {
				fmt.Fprintf(w, "Error reading weather data: ", err)
			}

			funcMap := template.FuncMap{
				"GetDay":        weather.GetDay,
				"GetImgNameMap": weather.GetImgNameMap,
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
				weatherData.Days,
			)

			err = tmpl.Execute(w, weatherTplData)
			if err != nil {
				fmt.Fprintf(w, "Error loading template:", err)
			}

			return
		}

		// Catch all other requests and return a 404.
		w.WriteHeader(fsthttp.StatusNotFound)
		fmt.Fprintf(w, "The page you requested could not be found\n")
	})
}
