package alpr

import (
	"context"

	"github.com/opencars/toolkit"
)

type Service struct {
	client *toolkit.Client
}

func NewService(client *toolkit.Client) *Service {
	return &Service{
		client: client,
	}
}
func (s *Service) Recognize(ctx context.Context, url string) ([]toolkit.ResultALPR, error) {
	resp, err := s.client.ALPR().Recognize(url)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
