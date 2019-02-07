package autoria

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type CarSearchResultResponse struct {
	Cars   []string `json:"ids"`
	Count  int64    `json:"count"`
	LastID int64    `json:"last_id"`
}

type CarResultResponse struct {
	SearchResult CarSearchResultResponse `json:"search_result"`
}

type CarSearchResponse struct {
	Result CarResultResponse `json:"result"`
}

func (api *API) GetSearchCars(params map[string]string) (*CarSearchResponse, error) {
	strParams := make([]string, 0)

	for k, v := range params {
		strParams = append(strParams, fmt.Sprintf("%s=%s", k, v))
	}

	resp, err := http.Get(api.BuildURL("/auto/search", strParams...))

	if err != nil {
		return nil, err
	}

	search := &CarSearchResponse{}
	err = json.NewDecoder(resp.Body).Decode(search)

	if err != nil {
		return nil, err
	}

	return search, nil
}

func ParseCarSearchParams(url string) (map[string]string, error) {
	query := strings.TrimPrefix(url, "https://auto.ria.com/search/?")
	params := strings.Split(query, "&")

	mapParams := make(map[string]string)

	for _, v := range params {
		lexemes := strings.Split(v, "=")

		if len(lexemes) != 2 {
			return nil, New("invalid parameters")
		}

		mapParams[lexemes[0]] = lexemes[1]
	}

	return mapParams, nil
}
