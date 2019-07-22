package handlers

import (
	"log"

	"github.com/opencars/bot/internal/bot"
)

func (h OpenCarsHandler) NotImplemented(msg *bot.Event) {
	if err := msg.SendHTML("Not implemented."); err != nil {
		log.Printf("send error: %s\n", err.Error())
	}
}
