package handlers

import (
	"bytes"
	"github.com/opencars/bot/internal/bot"
	"github.com/opencars/toolkit"
	"html/template"
	"log"
)

func (h OpenCarsHandler) RegistrationHandler(msg *bot.Event) {
	if err := msg.SetStatus(bot.ChatTyping); err != nil {
		log.Printf("action error: %s", err.Error())
	}

	if msg.Message.Text == "" {
		if err := msg.SendHTML("Номер відсутній"); err != nil {
			log.Printf("send error: %s\n", err.Error())
		}
		return
	}

	registration, err := h.client.Registration().FindByCode(msg.Message.Text)
	if err != nil {
		log.Printf("action error: %s\n", err)
	}

	tpl, err := template.ParseFiles("templates/registrations.tpl")
	if err != nil {
		log.Printf("action error: %s\n", err)
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
		log.Printf("action error: %s\n", err)
		return
	}

	if err := msg.SendHTML(buff.String()); err != nil {
		log.Printf("send error: %s\n", err.Error())
	}
}
