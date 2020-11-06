package handlers

import (
	"github.com/opencars/bot/internal/bot"
	"github.com/opencars/bot/pkg/logger"
)

func (h OpenCarsHandler) NotImplemented(msg *bot.Event) {
	if err := msg.SendHTML("Not implemented."); err != nil {
		logger.Errorf("send error: %s", err)
	}
}
