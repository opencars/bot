package autoria

import (
	"encoding/json"
	"net/http"
	"net/url"
)

type SearchResultResponse struct {
	Cars   []string `json:"ids"`
	Count  int64    `json:"count"`
	LastID int64    `json:"last_id"`
}

type SearchResponse struct {
	Result struct {
		SearchResult SearchResultResponse `json:"search_result"`
	} `json:"result"`
}

func (api *API) SearchCars(values url.Values) (*SearchResponse, error) {
	resp, err := http.Get(api.buildURL("/auto/search", values))

	if err != nil {
		return nil, err
	}

	search := &SearchResponse{}
	err = json.NewDecoder(resp.Body).Decode(search)

	if err != nil {
		return nil, err
	}

	return search, nil
}
