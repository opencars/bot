package operation

import (
	"context"
	"fmt"

	"github.com/opencars/toolkit"

	"github.com/opencars/bot/pkg/domain"
)

type Service struct {
	client *toolkit.OperationClient
}

func NewService(c *toolkit.Client) *Service {
	return &Service{client: c.Operation()}
}

func (s *Service) FindByNumber(_ context.Context, number string) ([]domain.Operation, error) {
	operations, err := s.client.FindByNumber(number)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", domain.ErrRemoteFailed, err)
	}

	result := make([]domain.Operation, 0)
	for i := range operations {
		result = append(result, domain.Operation{
			Person:      operations[i].Person,
			RegAddress:  operations[i].RegAddress,
			RegCode:     operations[i].RegCode,
			Reg:         operations[i].Reg,
			Date:        operations[i].Date,
			DepCode:     operations[i].DepCode,
			Dep:         operations[i].Dep,
			Brand:       operations[i].Brand,
			Model:       operations[i].Model,
			Year:        operations[i].Year,
			Color:       operations[i].Color,
			Kind:        operations[i].Kind,
			Body:        operations[i].Body,
			Purpose:     operations[i].Purpose,
			Fuel:        operations[i].Fuel,
			Capacity:    operations[i].Capacity,
			OwnWeight:   operations[i].OwnWeight,
			TotalWeight: operations[i].TotalWeight,
			Number:      operations[i].Number,
		})
	}

	return result, nil
}
