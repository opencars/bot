package handlers

import (
	"fmt"

	"github.com/opencars/bot/internal/bot"
	"github.com/opencars/bot/pkg/logger"
)

func StartHandler(msg *bot.Event) {
	if err := msg.SetStatus(bot.ChatTyping); err != nil {
		logger.Errorf("action error: %s", err)
	}

	text := fmt.Sprintf("Привіт, %s!", msg.Message.Chat.FirstName)
	if err := msg.Send(text); err != nil {
		logger.Errorf("send error: %s", err)
	}
}
