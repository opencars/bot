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

func (r *Image) Plate() (string, error) {
	if len(r.Recognized) < 1 {
		return "", errors.New("no plates found")
	} else if len(r.Recognized) > 1 {
		return "", errors.New("too many plates on the image")
	}

	candidates := r.Recognized[0].Candidates

	// Sort by Confidence.
	sort.Slice(candidates, func(i, j int) bool {
		return candidates[i].Confidence > candidates[j].Confidence
	})

	// Find first plates, that matches.
	for i := range candidates {
		if match.EuroPlates(candidates[i].Plate) {
			return candidates[i].Plate, nil
		}
	}

	return candidates[0].Plate, nil
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

	img := new(Image)
	if err := json.NewDecoder(resp.Body).Decode(img); err != nil {
		return nil, errors.New("invalid response body")
	}

	return img, nil
}
