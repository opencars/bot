// Copyright (C) 2019 Ali Shanaakh, github@shanaakh.pro
// This software may be modified and distributed under the terms
// of the MIT license. See the LICENSE file for details.

package openalpr

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"
)

type PlateRecognizerCoordinate struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type ResponseResult struct {
	Confidence       float64                     `json:"confidence"`
	Coordinates      []PlateRecognizerCoordinate `json:"coordinates"`
	Plate            string                      `json:"plate"`
	PlateIndex       int                         `json:"plate_index"`
	ProcessingTimeMs float64                     `json:"processing_time_ms"`
	Region           string                      `json:"region"`
	RegionConfidence int                         `json:"region_confidence"`
	RequestedTopN    int                         `json:"requested_topn"`
}

type Response struct {
	ImgHeight        int              `json:"img_height"`
	ImgWidth         int              `json:"img_width"`
	ProcessingTimeMs float64          `json:"processing_time_ms"`
	Results          []ResponseResult `json:"results"`
}

func (resp Response) Plate() (string, error) {
	if len(resp.Results) < 1 {
		return "", errors.New("plates was not recognized")
	} else if len(resp.Results) > 1 {
		return "", errors.New("too much candidates on the photo")
	}

	return resp.Results[0].Plate, nil
}


type API struct {
	URL string
}

func (client *API) Recognize(imageURL string) (*Response, error) {
	URL := fmt.Sprintf("%s/v2/identify/plate?image_url=%s", client.URL, imageURL)

	resp, err := http.Get(URL)

	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	recognizerResponse := new(Response)

	err = json.Unmarshal(body, recognizerResponse)

	if err != nil {
		return nil, err
	}

	return recognizerResponse, nil
}