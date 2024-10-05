package opencage

import (
	"context"
	"fmt"
	"github.com/fastly/compute-sdk-go/fsthttp"
	"weather-app/fastly"
)

func GetLocationMetaData(ctx context.Context, lat string, long string) (*fsthttp.Response, error) {
	key, err := fastly.GetSecretStoreKey("keys", "OPEN_CAGE")
	if err != nil {
		return nil, fmt.Errorf("cannot read key: %v", err)
	}

	url := fmt.Sprintf("https://api.opencagedata.com/geocode/v1/json?q=%s+%s&key=%s", lat, long, key)

	req, err := fsthttp.NewRequest("GET", url, nil)
	if err != nil {

		return nil, fmt.Errorf("error making HTTP request: %v", err)
	}

	resp, err := req.Send(ctx, "opencage")

	if resp.StatusCode == fsthttp.StatusOK {
		return resp, nil
	}

	defer resp.Body.Close()

	return nil, fmt.Errorf("error fetching response: %s", err)
}
