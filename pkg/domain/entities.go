package domain

import (
	"time"
)

type User struct {
	ID           int
	FirstName    string
	LastName     *string
	UserName     *string
	LanguageCode *string
}

type Message struct {
	ID   int
	User User
	Text string
	Time time.Time
}

type Result struct {
	Request  Request
	Vehicles map[string]*Vehicle
}

type Request struct {
	Number *string
	VIN    *string
}

type Vehicle struct {
	VIN           string
	FirstRegDate  time.Time
	Brand         string
	Model         string
	Year          int32
	Registrations []Registration
}

type Registration struct {
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
}
