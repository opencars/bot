package model

type Result struct {
	Request  Request
	Vehicles map[string]*Vehicle
}
