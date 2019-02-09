package autoria

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type SearchResultResponse struct {
	Cars   []string `json:"ids"`
	Count  int64    `json:"count"`
	LastID int64    `json:"last_id"`
}

type SearchResponse struct {
	Result struct{
		SearchResult SearchResultResponse `json:"search_result"`
	} `json:"result"`
}

func (api *API) SearchCars(params map[string]string) (*SearchResponse, error) {
	strParams := make([]string, 0)

	for k, v := range params {
		strParams = append(strParams, fmt.Sprintf("%s=%s", k, v))
	}

	resp, err := http.Get(api.BuildURL("/auto/search", strParams...))

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

func ParseSearchParams(url string) (map[string]string, error) {
	query := strings.TrimPrefix(url, "https://auto.ria.com/search/?")
	params := strings.Split(query, "&")

	mapParams := make(map[string]string)

	for _, v := range params {
		lexemes := strings.Split(v, "=")

		if len(lexemes) != 2 {
			return nil, NewErr("invalid parameters")
		}

		mapParams[lexemes[0]] = lexemes[1]
	}

	return mapParams, nil
}
