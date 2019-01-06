// Copyright (C) 2019 Ali Shanaakh, github@shanaakh.pro
// This software may be modified and distributed under the terms
// of the MIT license. See the LICENSE file for details.

package autoria

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"net/http"
	"strings"
)

type CarSearchResultResponse struct {
	CarsIDs []string `json:"ids"`
	Count   int64    `json:"count"`
	LastID  int64    `json:"last_id"`
}

type CarResultResponse struct {
	SearchResult CarSearchResultResponse `json:"search_result"`
}

type CarSearchResponse struct {
	Result CarResultResponse `json:"result"`
}

func (api *API) GetSearchCars(params map[string]string) (*CarSearchResponse, error) {
	strParams := make([]string, 0)

	fmt.Println(params)

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
			return nil, errors.New("invalid params")
		}

		mapParams[lexemes[0]] = lexemes[1]
	}

	return mapParams, nil
}
