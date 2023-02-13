package domain

import (
	"context"

	"github.com/opencars/bot/pkg/domain/model"
)

type MessageRepository interface {
	Create(ctx context.Context, msg *model.Message) error
}
