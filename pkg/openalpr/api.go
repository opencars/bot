package openalpr

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"sort"

	"github.com/shal/opencars-bot/pkg/match"
)

type API struct {
	URI string
}

func New(URI string) *API {
	return &API{URI: URI}
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
	URI := fmt.Sprintf("%s/v2/identify/plate?image_url=%s", api.URI, uri)

	if _, err := url.ParseRequestURI(uri); err != nil {
		return nil, err
	}

	resp, err := http.Get(URI)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	img := Image{}
	if err := json.NewDecoder(resp.Body).Decode(&img); err != nil {
		return nil, errors.New("invalid response body")
	}

	return &img, nil
}
