package model

import "time"

type Vehicle struct {
	VIN          string
	FirstRegDate time.Time
	Brand        string
	Model        string
	Year         int32
	Actions      []Action
}

type Action struct {
	VIN         string
	Code        string
	Number      string
	Brand       string
	Model       string
	Color       string
	Kind        string
	Year        int32
	TotalWeight int32
	OwnWeight   int32
	Capacity    int32
	Fuel        string
	Category    string
	NumSeating  int32
	Date        time.Time
	Body        string
	Purpose     string
}
