package bot

import (
	"bytes"
	"context"
	"errors"
	"html/template"

	"gopkg.in/tucnak/telebot.v2"

	"github.com/opencars/bot/pkg/domain"
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

	if _, err := b.client.Send(m.Chat, buff.String(), opts); err != nil {
		return err
	}

	return nil
}

func (b *Bot) FindByImage(ctx context.Context, m *telebot.Message) error {
	var fileID string

	if m.Document != nil {
		fileID = m.Document.FileID
	}

	if m.Photo != nil {
		fileID = m.Photo.FileID
	}

	url, err := b.client.FileURLByID(fileID)
	if err != nil {
		return err
	}

	vehicle, err := b.vehicle.FindByImage(ctx, url)
	if errors.Is(err, domain.ErrNotRecognized) {
		if _, err := b.client.Reply(m, "not found"); err != nil {
			return err
		}
		return nil
	}

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

	if _, err := b.client.Send(m.Chat, buff.String(), opts); err != nil {
		return err
	}

	return nil
}
