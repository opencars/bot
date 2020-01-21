package handlers

import (
	"fmt"
	"log"
	"strings"

	"github.com/opencars/bot/internal/bot"
)

func (h OpenCarsHandler) PlatesHandler(msg *bot.Event) {
	if err := msg.SetStatus(bot.ChatTyping); err != nil {
		log.Printf("action error: %s", err)
	}

	number := strings.TrimSpace(strings.TrimPrefix(msg.Message.Text, "/number"))
	if number == "" {
		if err := msg.SendHTML("Номер відсутній"); err != nil {
			log.Printf("send error: %s\n", err)
		}
		return
	}

	text, err := h.getInfoByNumber(number)
	if err != nil {
		log.Println(err)
		return
	}

	if err := msg.SendHTML(text); err != nil {
		log.Printf("send error: %s\n", err)
		return
	}

	text, err = h.getRegistrationsByNumber(number)
	if err != nil {
		log.Println(err)
		return
	}

	if err := msg.SendHTML(text); err != nil {
		log.Printf("send error: %s\n", err)
		return
	}
}

func (h OpenCarsHandler) ReportByVIN(msg *bot.Event) {
	if err := msg.SetStatus(bot.ChatTyping); err != nil {
		log.Printf("action error: %s\n", err)
		return
	}

	vin := strings.TrimSpace(strings.TrimPrefix(msg.Message.Text, "/vin"))
	if vin == "" {
		if err := msg.SendHTML("Номер відсутній"); err != nil {
			log.Printf("send error: %s\n", err)
		}
		return
	}

	fmt.Println(vin)
	text, err := h.GetReportByVIN(vin)
	if err != nil {
		log.Printf("error: %s\n", err)
		return
	}

	if err := msg.SendHTML(text); err != nil {
		log.Printf("send error: %s\n", err)
		return
	}
}
