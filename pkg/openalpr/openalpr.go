package openalpr

import (
	"encoding/json"
	"fmt"
	"github.com/shal/robot/pkg/match"
	"io/ioutil"
	"net/http"
	"sort"
)

type API struct {
	URL string
}

type Candidate struct {
	Confidence float64
	Plate      string
}

type Coordinate struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type Result struct {
	Candidates       []Candidate  `json:"candidates"`
	Confidence       float64      `json:"confidence"`
	Coordinates      []Coordinate `json:"coordinates"`
	Plate            string       `json:"plate"`
	PlateIndex       int          `json:"plate_index"`
	ProcessingTimeMs float64      `json:"processing_time_ms"`
	Region           string       `json:"region"`
	RegionConfidence int          `json:"region_confidence"`
	RequestedTopN    int          `json:"requested_topn"`
}

type Response struct {
	ImageHeight    int      `json:"img_height"`
	ImageWidth     int      `json:"img_width"`
	ProcessingTime float64  `json:"processing_time_ms"`
	Results        []Result `json:"results"`
}

func (c Candidate) Priority() int {
	priority := int(c.Confidence)

	if len(c.Plate) == 8 {
		priority += 100
	}

	if match.EuroPlates(c.Plate) {
		priority += 100
	}

	return priority
}

func (r Response) Plate() (string, error) {
	if len(r.Results) < 1 {
		return "", New("plates was not recognized")
	} else if len(r.Results) > 1 {
		return "", New("too much candidates on the photo")
	}

	candidates := r.Results[0].Candidates

	sort.Slice(candidates, func(i, j int) bool {
		return candidates[i].Priority() > candidates[j].Priority()
	})

	return candidates[0].Plate, nil
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
