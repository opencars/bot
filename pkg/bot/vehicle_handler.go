package bot

import (
	"context"

	"gopkg.in/tucnak/telebot.v2"
)

func (b *Bot) StartHandler() Handler {
	return Wrapper(func(ctx context.Context, m *telebot.Message) error {
		if _, err := b.client.Send(m.Chat, "Привіт!"); err != nil {
			return err
		}

		return nil
	})
}

func (b *Bot) FindByNumber(ctx context.Context, m *telebot.Message) error {
	vehicle, err := b.vehicle.FindByNumber(ctx, m.Text)
	if err != nil {
		return err
	}

	opts := &telebot.SendOptions{
		ParseMode: telebot.ModeHTML,
	}

	if _, err := b.client.Send(m.Chat, vehicle, opts); err != nil {
		return err
	}

	return nil
}

func (b *Bot) FindByVIN(ctx context.Context, m *telebot.Message) error {
	vehicle, err := b.vehicle.FindByVIN(ctx, m.Text)
	if err != nil {
		return err
	}

	opts := telebot.SendOptions{
		ParseMode: telebot.ModeHTML,
	}

	if _, err := b.client.Send(m.Chat, vehicle, opts); err != nil {
		return err
	}

	return nil
}
