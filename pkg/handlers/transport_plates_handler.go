package handlers

import (
	"bytes"
	"html/template"
	"log"

	"github.com/shal/opencars-bot/internal/bot"
	"github.com/shal/opencars-bot/pkg/opencars"
)

type OpenCarsHandler struct {
	OpenCars *opencars.API
}

func (h OpenCarsHandler) Handle(msg *bot.Message) {
	if err := msg.Send(bot.ChatTyping); err != nil {
		log.Printf("action error: %s", err.Error())
	}

	transport, err := h.OpenCars.Search(msg.Text())

	if err != nil {
		if err := msg.Send("Вибач. Щось пішло не так 😢"); err != nil {
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
		msg.Text(),
	}); err != nil {
		log.Println(err)
	}

	if err := msg.SendHTML(buff.String()); err != nil {
		log.Printf("send error: %s\n", err.Error())
	}
}
