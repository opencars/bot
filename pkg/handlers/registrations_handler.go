package handlers

import (
	"bytes"
	"html/template"

	"github.com/opencars/toolkit"

	"github.com/opencars/bot/internal/bot"
	"github.com/opencars/bot/pkg/logger"
)

func (h OpenCarsHandler) RegistrationHandler(msg *bot.Event) {
	if err := msg.SetStatus(bot.ChatTyping); err != nil {
		logger.Errorf("action error: %s", err)
	}

	if msg.Message.Text == "" {
		if err := msg.SendHTML("Номер відсутній"); err != nil {
			logger.Errorf("send error: %s", err)
		}
		return
	}

	registration, err := h.client.Registration().FindByCode(msg.Message.Text)
	if err != nil {
		logger.Errorf("action error: %s", err)
	}

	tpl, err := template.ParseFiles("templates/registrations.tpl")
	if err != nil {
		logger.Errorf("action error: %s", err)
		return
	}

	type payload struct {
		Registrations []toolkit.Registration
		Number, Code  string
	}

	var tmp payload
	tmp.Code = msg.Message.Text
	if registration != nil {
		tmp.Registrations = []toolkit.Registration{*registration}
	}

	buff := bytes.Buffer{}
	if err := tpl.Execute(&buff, tmp); err != nil {
		logger.Errorf("action error: %s", err)
		return
	}

	if err := msg.SendHTML(buff.String()); err != nil {
		logger.Errorf("send error: %s", err.Error())
	}
}
