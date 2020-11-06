package handlers

import (
	"strings"

	"github.com/opencars/bot/internal/bot"
	"github.com/opencars/bot/pkg/logger"
)

func (h OpenCarsHandler) PlatesHandler(msg *bot.Event) {
	if err := msg.SetStatus(bot.ChatTyping); err != nil {
		logger.Errorf("action error: %s", err)
	}

	number := strings.TrimPrefix(msg.Message.Text, "/number")
	number = strings.TrimSpace(number)
	number = strings.ToUpper(number)

	if number == "" {
		if err := msg.SendHTML("Номер відсутній"); err != nil {
			logger.Errorf("send error: %s", err)
		}
		return
	}

	text, err := h.getInfoByNumber(number)
	if err != nil {
		logger.Errorf("failed to get info by number: %s", err)
		return
	}

	if err := msg.SendHTML(text); err != nil {
		logger.Errorf("send error: %s", err)
		return
	}

	text, err = h.getRegistrationsByNumber(number)
	if err != nil {
		logger.Errorf("failed to get registrations by number: %s", err)
		return
	}

	if err := msg.SendHTML(text); err != nil {
		logger.Errorf("send error: %s", err)
		return
	}
}

func (h OpenCarsHandler) ReportByVIN(msg *bot.Event) {
	if err := msg.SetStatus(bot.ChatTyping); err != nil {
		logger.Errorf("action error: %s", err)
		return
	}

	vin := strings.TrimSpace(strings.TrimPrefix(msg.Message.Text, "/vin"))
	if vin == "" {
		if err := msg.SendHTML("Номер відсутній"); err != nil {
			logger.Errorf("send error: %s", err)
		}
		return
	}

	text, err := h.GetReportByVIN(vin)
	if err != nil {
		logger.Errorf("error: %s", err)
		return
	}

	if err := msg.SendHTML(text); err != nil {
		logger.Errorf("send error: %s", err)
		return
	}
}
