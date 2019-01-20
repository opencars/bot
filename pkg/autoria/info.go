// Copyright (C) 2019 Ali Shanaakh, github@shanaakh.pro
// This software may be modified and distributed under the terms
// of the MIT license. See the LICENSE file for details.

package autoria

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Car struct {
	Description  string `json:"description"`
	Version      string `json:"version"`
	OnModeration bool   `json:"onModeration"`
	Year         int    `json:"year"`
	AutoID       int    `json:"autoId"`
	StatusID     int    `json:"statusId"`
	WithVideo    bool   `json:"withVideo"`
	Race         string `json:"race"`
	RaceInt      int    `json:"raceInt"`
	FuelName     string `json:"fuelName"`
	GearboxName  string `json:"gearboxName"`
	IsSold       bool   `json:"isSold"`
	MainCurrency string `json:"mainCurrency"`
	CategoryID   int    `json:"categoryId"`
}

type CarInfoResponse struct {
	LocationCityName string `json:"locationCityName"`
	MarkName         string `json:"markName"`
	MarkID           uint32 `json:"markId"`
	ModelName        string `json:"modelName"`
	ModelID          uint32 `json:"modelId"`
	LinkToView       string `json:"linkToView"`
	PriceUSD         int64  `json:"USD"`
	PriceHRN         int64  `json:"UAH"`
	PriceEUR         int64  `json:"EUR"`
	Car              Car    `json:"autoData"`
}

func (api *API) CarInfo(ID string) (car *CarInfoResponse, err error) {
	resp, err := http.Get(api.BuildURL("/auto/info", fmt.Sprintf("auto_id=%s", ID)))

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
