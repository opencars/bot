package handlers

import (
	"bytes"
	"html/template"
	"log"

	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/shal/robot/pkg/opencars"
)

type OpenCarsHandler struct {
	OpenCars *opencars.API
}

func (h OpenCarsHandler) Handle(bot *tgbotapi.BotAPI, msg *tgbotapi.Message) {
	transport, err := h.OpenCars.Search(msg.Text)

	if err != nil {
		Send(bot, msg.Chat, "Вибач. Щось пішло не так 😢")
		return
	}

	tpl, err := template.ParseFiles("templates/opencars_info.tpl")
	if err != nil {
		log.Println(err)
	}

	buff := bytes.Buffer{}
	if err := tpl.Execute(&buff, struct {
		Cars   []opencars.Transport
		Number string
	}{
		transport,
		msg.Text,
	}); err != nil {
		log.Println(err)
	}

	if err := SendHTML(bot, msg.Chat, buff.String()); err != nil {
		log.Printf("send error: %s\n", err.Error())
	}
}
