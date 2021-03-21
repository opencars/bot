package vehicle

import (
	"context"

	"github.com/opencars/grpc/pkg/core"
	"google.golang.org/grpc"
)

type Service struct {
	c core.VehicleServiceClient
}

func NewService(addr string) (*Service, error) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	return &Service{
		c: core.NewVehicleServiceClient(conn),
	}, nil
}

func (s *Service) FindByNumber(ctx context.Context, number string) (*core.Result, error) {
	resp, err := s.c.FindByNumber(ctx, &core.NumberRequest{Number: number})
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *Service) FindByVIN(ctx context.Context, vin string) (*core.Result, error) {
	resp, err := s.c.FindByVIN(ctx, &core.VINRequest{Vin: vin})
	if err != nil {
		return nil, err
	}

	return resp, nil
}
