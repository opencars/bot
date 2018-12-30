// Copyright (C) 2019 Ali Shanaakh, github@shanaakh.pro
// This software may be modified and distributed under the terms
// of the MIT license. See the LICENSE file for details.

package autoria_api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type CarInfoResponse struct {
	LocationCityName string `json:"locationCityName"`
	MarkName         string `json:"markName"`
	MakeID           uint32 `json:"markId"`
	ModelName        string `json:"modelName"`
	ModelID          uint32 `json:"modelId"`
	LinkToView       string `json:"linkToView"`
}

func (api *API) GetCarInfo(ID string) (car CarInfoResponse) {
	url := fmt.Sprintf("/info%s?api_key=%s&auto_id=%s", api.base, api.key, ID)

	fmt.Print(url)

	resp, err := http.Get(url)

	if err != nil {
		panic(err.Error())
	}

	err = json.NewDecoder(resp.Body).Decode(&car)

	if err != nil {
		panic(err.Error())
	}

	return car
}
