package foursquare

import (
	"context"
	"fmt"
	"github.com/fastly/compute-sdk-go/fsthttp"
	"weather-app/fastly"
)

func GetNearby(ctx context.Context, lat string, long string, category string) (*fsthttp.Response, error) {
	key, err := fastly.GetSecretStoreKey("keys", "FOUR_SQUARE")
	if err != nil {
		return nil, fmt.Errorf("cannot read key: %v", err)
	}

	cat, err := getCategoryId(category)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("https://api.foursquare.com/v3/places/search?ll=%s,%s&categories=%s", lat, long, cat)

	req, err := fsthttp.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error making HTTP request: %v", err)
	}

	req.Header.Add("Authorization", key)

	resp, err := req.Send(ctx, "foursquare")

	if resp.StatusCode == fsthttp.StatusOK {
		return resp, nil
	}

	defer resp.Body.Close()

	return nil, fmt.Errorf("error fetching response: %s", err)
}

func getCategoryId(name string) (string, error) {
	categoryMap := map[string]string{
		"beaches":     "4bf58dd8d48988d1e2941735",
		"parks":       "4bf58dd8d48988d1e2931735",
		"restaurants": "4bf58dd8d48988d1e0931736",
		"trails":      "4bf58dd8d48988d159941735",
		"breweries":   "50327c8591d4c4b30a586d5d",
	}
	if categoryName, exists := categoryMap[name]; exists {
		return categoryName, nil
	} else {
		return "", fmt.Errorf("category not found")
	}
}
