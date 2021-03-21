package bot

import (
	"bytes"
	"context"
	"html/template"

	"gopkg.in/tucnak/telebot.v2"

	"github.com/opencars/bot/pkg/logger"
)

func (b *Bot) FindByNumber(ctx context.Context, m *telebot.Message) error {
	vehicle, err := b.vehicle.FindByNumber(ctx, m.Text)
	if err != nil {
		return err
	}

	opts := &telebot.SendOptions{
		ParseMode: telebot.ModeHTML,
	}

	tpl, err := template.ParseFiles("templates/vehicles.tpl")
	if err != nil {
		return err
	}

	buff := bytes.Buffer{}
	if err := tpl.Execute(&buff, vehicle); err != nil {
		logger.Errorf(err.Error())
	}

	logger.Debugf("%s", buff.String())

	if _, err := b.client.Send(m.Chat, buff.String(), opts); err != nil {
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

	tpl, err := template.ParseFiles("templates/vehicle.tpl")
	if err != nil {
		return err
	}

	buff := bytes.Buffer{}
	if err := tpl.Execute(&buff, vehicle); err != nil {
		logger.Errorf(err.Error())
	}

	if _, err := b.client.Send(m.Chat, buff.String(), opts); err != nil {
		return err
	}

	return nil
}
