package message

import (
	"context"

	"github.com/opencars/bot/pkg/domain"
	"github.com/opencars/bot/pkg/domain/model"
)

type Service struct {
	r domain.MessageRepository
}

func NewService(r domain.MessageRepository) (*Service, error) {
	return &Service{
		r: r,
	}, nil
}

func (s *Service) Create(ctx context.Context, msg *model.Message) error {
	return s.r.Create(ctx, msg)
}
