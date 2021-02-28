package bot

import (
	"context"

	"gopkg.in/tucnak/telebot.v2"

	"github.com/opencars/bot/pkg/domain"
	"github.com/opencars/bot/pkg/logger"
)

type Handler func(m *telebot.Message)

func Wrapper(fn func(ctx context.Context, m *telebot.Message) error) Handler {
	return func(m *telebot.Message) {
		if err := fn(context.TODO(), m); err != nil {
			switch err.(type) {
			case domain.Error:
				logger.Errorf("%s", err)
			default:
				logger.Errorf("%s", err)
			}

			return
		}
	}
}
