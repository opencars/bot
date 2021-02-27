package registration

import (
	"context"
	"fmt"

	"github.com/opencars/toolkit"

	"github.com/opencars/bot/pkg/domain"
)

type Service struct {
	client *toolkit.RegistrationClient
}

func NewService(c *toolkit.Client) *Service {
	return &Service{client: c.Registration()}
}

// TODO: Add context support into toolkit, if gRPC is not implemented.
func (s *Service) FindByNumber(_ context.Context, number string) ([]domain.Registration, error) {
	registrations, err := s.client.FindByNumber(number)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", domain.ErrRemoteFailed, err)
	}

	result := make([]domain.Registration, 0)
	for i := range registrations {
		result = append(result, domain.Registration{
			Brand:        registrations[i].Brand,
			Capacity:     registrations[i].Capacity,
			Color:        registrations[i].Color,
			FirstRegDate: registrations[i].FirstRegDate,
			Date:         registrations[i].Date,
			Fuel:         registrations[i].Fuel,
			Kind:         registrations[i].Kind,
			Year:         registrations[i].Year,
			Model:        registrations[i].Model,
			Code:         registrations[i].Code,
			Number:       registrations[i].Number,
			NumSeating:   registrations[i].NumSeating,
			NumStanding:  registrations[i].NumStanding,
			OwnWeight:    registrations[i].OwnWeight,
			RankCategory: registrations[i].RankCategory,
			TotalWeight:  registrations[i].TotalWeight,
			VIN:          registrations[i].VIN,
		})
	}

	return result, nil
}

// TODO: Add context support into toolkit, if gRPC is not implemented.
func (s *Service) FindByVIN(_ context.Context, vin string) ([]domain.Registration, error) {
	registrations, err := s.client.FindByVIN(vin)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", domain.ErrRemoteFailed, err)
	}

	result := make([]domain.Registration, 0)
	for i := range registrations {
		result = append(result, domain.Registration{
			Brand:        registrations[i].Brand,
			Capacity:     registrations[i].Capacity,
			Color:        registrations[i].Color,
			FirstRegDate: registrations[i].FirstRegDate,
			Date:         registrations[i].Date,
			Fuel:         registrations[i].Fuel,
			Kind:         registrations[i].Kind,
			Year:         registrations[i].Year,
			Model:        registrations[i].Model,
			Code:         registrations[i].Code,
			Number:       registrations[i].Number,
			NumSeating:   registrations[i].NumSeating,
			NumStanding:  registrations[i].NumStanding,
			OwnWeight:    registrations[i].OwnWeight,
			RankCategory: registrations[i].RankCategory,
			TotalWeight:  registrations[i].TotalWeight,
			VIN:          registrations[i].VIN,
		})
	}

	return result, nil
}
