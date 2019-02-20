package handlers

import (
	"bytes"
	"github.com/shal/opencars-bot/internal/bot"
	"html/template"
	"log"

	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/shal/opencars-bot/pkg/opencars"
)

type OpenCarsHandler struct {
	OpenCars *opencars.API
}

func (h OpenCarsHandler) Handle(api *tgbotapi.BotAPI, msg *tgbotapi.Message) {
	transport, err := h.OpenCars.Search(msg.Text)

	if err != nil {
		if err := bot.Send(api, msg.Chat, "Вибач. Щось пішло не так 😢"); err != nil {
			log.Printf("send error: %s", err.Error())
		}
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

	if err := bot.SendHTML(api, msg.Chat, buff.String()); err != nil {
		log.Printf("send error: %s\n", err.Error())
	}
}
