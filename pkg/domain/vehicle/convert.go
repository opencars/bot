package vehicle

import (
	"time"

	"github.com/opencars/bot/pkg/domain"
	"github.com/opencars/grpc/pkg/core"
)

func convert(in *core.Result) *domain.Result {
	result := domain.Result{
		Vehicles: make(map[string]*domain.Vehicle, len(in.Vehicles)),
	}

	for _, v := range in.Vehicles {
		vehicle := domain.Vehicle{
			Brand: v.Brand,
			Model: v.Model,
			Year:  v.Year,
		}

		if v.Vin != nil {
			vehicle.VIN = v.Vin.Value
		}

		if v.FirstRegDate != nil {
			vehicle.FirstRegDate = time.Date(
				int(v.FirstRegDate.Year),
				time.Month(v.FirstRegDate.Month),
				int(v.FirstRegDate.Day),
				0, 0, 0, 0,
				time.UTC,
			)
		}

		for _, r := range v.Registrations {
			vehicle.Registrations = append(vehicle.Registrations, domain.Registration{
				VIN:         r.Vin,
				Code:        r.Code,
				Number:      r.Number,
				Brand:       r.Brand,
				Model:       r.Model,
				Color:       r.Color,
				Kind:        r.Kind,
				Year:        r.Year,
				TotalWeight: r.TotalWeight,
				OwnWeight:   r.OwnWeight,
				Capacity:    r.Capacity,
				Fuel:        r.Fuel,
				Category:    r.Category,
				NumSeating:  r.NumSeating,
				Date: time.Date(
					int(r.Date.Year),
					time.Month(r.Date.Month),
					int(r.Date.Day),
					0, 0, 0, 0,
					time.UTC,
				),
			})
		}

		for _, r := range v.Operations {
			vehicle.Operations = append(vehicle.Operations, domain.Operation{
				VIN:         r.Vin,
				Number:      r.Number,
				Brand:       r.Brand,
				Model:       r.Model,
				Color:       r.Color,
				Kind:        r.Kind,
				Year:        r.Year,
				TotalWeight: r.TotalWeight,
				OwnWeight:   r.OwnWeight,
				Capacity:    r.Capacity,
				Fuel:        r.Fuel,
				Date: time.Date(
					int(r.Date.Year),
					time.Month(r.Date.Month),
					int(r.Date.Day),
					0, 0, 0, 0,
					time.UTC,
				),
			})
		}

		result.Vehicles[vehicle.VIN] = &vehicle
	}

	return &result
}
