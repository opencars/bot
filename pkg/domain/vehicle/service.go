package vehicle

import (
	"context"

	"github.com/opencars/bot/pkg/domain"
	"github.com/opencars/bot/pkg/domain/model"
	"github.com/opencars/grpc/pkg/core"
	"github.com/opencars/translit"
	"google.golang.org/grpc"
)

type Service struct {
	c core.VehicleServiceClient
	r domain.Recognizer
}

func NewService(addr string, r domain.Recognizer) (*Service, error) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	return &Service{
		c: core.NewVehicleServiceClient(conn),
		r: r,
	}, nil
}

func (s *Service) FindByNumber(ctx context.Context, number string) (*model.Result, error) {
	resp, err := s.c.FindByNumber(ctx, &core.NumberRequest{Number: number})
	if err != nil {
		return nil, err
	}

	res := convert(resp)
	res.Request = model.Request{Number: &number}

	return res, nil
}

func (s *Service) FindByVIN(ctx context.Context, vin string) (*model.Result, error) {
	resp, err := s.c.FindByVIN(ctx, &core.VINRequest{Vin: vin})
	if err != nil {
		return nil, err
	}

	res := convert(resp)
	res.Request = model.Request{VIN: &vin}

	return res, nil
}

func (s *Service) FindByImage(ctx context.Context, url string) (*model.Result, error) {
	plates, err := s.r.Recognize(ctx, url)
	if err != nil {
		return nil, err
	}

	for i, p := range plates {
		plates[i].Plate = translit.ToUA(p.Plate)
	}

	if len(plates) == 0 {
		return nil, model.ErrNotRecognized
	}

	return s.FindByNumber(ctx, plates[0].Plate)
}
