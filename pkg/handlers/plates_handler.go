package handlers

import (
	"log"

	"github.com/shal/opencars-bot/internal/bot"
)

func (h OpenCarsHandler) PlatesHandler(msg *bot.Event) {
	if err := msg.SetStatus(bot.ChatTyping); err != nil {
		log.Printf("action error: %s", err.Error())
	}

	text, err := h.getInfoByPlates(msg.Message.Text)
	if err != nil {
		log.Println(err.Error())
	}

	if err := msg.SendHTML(text); err != nil {
		log.Printf("send error: %s\n", err.Error())
	}
}
