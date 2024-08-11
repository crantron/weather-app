<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Weather Report</title>
    <link rel="stylesheet" href="/static/css/tailwind.css" />
    <body style="background: linear-gradient(to top, #09203f 0%, #537895 100%);">
    <h2>Written in GO Lang on the Edge - App Build Version: {{.Version}}</h2>
    <div>
        <strong>City:</strong> {{.City}}<br>
        <strong>Country:</strong> {{.Country}}<br>
        <strong>Timezone:</strong> {{.TimeZone}}<br>
        <strong>Latitude:</strong> {{.Latitude}}<br>
        <strong>Longitude:</strong> {{.Longitude}}<br>
        <strong>IP Address:</strong> {{.IPAddress}}<br><br><br>
    </div>
    <div>
        <h3>Weather in: {{.ResolvedAddress}}</h3>
        <p>Current temperature: {{.Temperature}}°F</p>
        <p>Conditions: {{.Conditions}}</p>
        <h4>Forecast:</h4>
        <div class="container text-center">
            <div class="row">
                {{range .Days}}
                    <div class="col" style="margin-bottom: 2rem;">
                        <div class="card text-bg-light" style="width: 18rem;">
                            <img src="https://placehold.it/300x200" alt="" class="card-img-top">
                            <div class="card-body">
                                <h5 class="card-title">{{.Datetime}}</h5>
                                <p class="card-text">Max: {{.TemperatureMax}}°F, Min: {{.TemperatureMin}}°F, Conditions: {{.Conditions}}</p>
                            </div>
                        </div>
                    </div>
                {{end}}
            </div>
        </div>
    </div>
</body>
</html>
