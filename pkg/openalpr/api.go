package openalpr

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"time"

	"github.com/opencars/bot/pkg/match"
)

type API struct {
	client *http.Client
	URI    string
}

func New(URI string) *API {
	return &API{
		URI: URI,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (r *Image) Plates() ([]string, error) {
	if len(r.Recognized) < 1 {
		return nil, errors.New("no plates found")
	}

	plates := make([]string, 0)
	for _, recognized := range r.Recognized {
		candidates := recognized.Candidates

		// Sort by confidence.
		sort.Slice(candidates, func(i, j int) bool {
			return candidates[i].Confidence > candidates[j].Confidence
		})


		// Find first plates, that matches.
		plate := candidates[0].Plate
		for _, candidate := range candidates {
			if match.EuroPlates(candidate.Plate) {
				plate = candidate.Plate
				break
			}
		}

		plates = append(plates, plate)
	}

	return plates, nil
}

func (api *API) Recognize(uri string) (*Image, error) {
	URI := fmt.Sprintf("%s/api/v1/alpr/private/recognize?image_url=\"=%s", api.URI, uri)

	if _, err := url.ParseRequestURI(uri); err != nil {
		return nil, err
	}

	resp, err := api.client.Get(URI)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var img Image
	if err := json.NewDecoder(resp.Body).Decode(&img); err != nil {
		return nil, fmt.Errorf("invalid response body: %w", err)
	}

	return &img, nil
}
