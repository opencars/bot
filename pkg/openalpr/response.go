package openalpr

type Candidate struct {
	Confidence float64
	Plate      string
}

type Coordinate struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type Plate struct {
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

type Image struct {
	Height         int     `json:"img_height"`
	Width          int     `json:"img_width"`
	ProcessingTime float64 `json:"processing_time_ms"`
	Recognized     []Plate `json:"results"`
}
