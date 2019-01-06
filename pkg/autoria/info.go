// Copyright (C) 2019 Ali Shanaakh, github@shanaakh.pro
// This software may be modified and distributed under the terms
// of the MIT license. See the LICENSE file for details.

package autoria

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

func (api *API) GetCarInfo(ID string) (car *CarInfoResponse, err error) {
	resp, err := http.Get(api.BuildURL("/auto/info", fmt.Sprintf("auto_id=%s", ID)))

	fmt.Println("ID: ", ID)
	fmt.Println("Response: ", resp)

	if err != nil {
		return nil, err
	}

	err = json.NewDecoder(resp.Body).Decode(&car)

	if err != nil {
		return nil, err
	}

	// Add prefix with website link.
	car.LinkToView = "https://auto.ria.com" + car.LinkToView

	return car, nil
}
