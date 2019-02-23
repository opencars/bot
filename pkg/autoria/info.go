package autoria

import (
	"encoding/json"
	"net/http"
	"net/url"
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

type CarInfo struct {
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

func (api *API) CarInfo(ID string) (car *CarInfo, err error) {
	values := url.Values{}
	values.Set("auto_id=", ID)
	resp, err := http.Get(api.buildURL("/auto/info", values))

	if err != nil {
		return nil, err
	}

	if err := json.NewDecoder(resp.Body).Decode(&car); err != nil {
		return nil, err
	}

	// Add prefix with website link.
	car.LinkToView = "https://auto.ria.com" + car.LinkToView

	return car, nil
}
