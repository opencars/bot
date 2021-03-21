package domain

import "context"

type MessageRepository interface {
	Create(ctx context.Context, msg *Message) error
}
