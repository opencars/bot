package handlers

import (
	"fmt"
	"log"

	"github.com/opencars/bot/internal/bot"
)

func StartHandler(msg *bot.Event) {
	if err := msg.SetStatus(bot.ChatTyping); err != nil {
		log.Printf("action error: %s", err.Error())
	}

	text := fmt.Sprintf("Привіт, %s!", msg.Message.Chat.FirstName)
	if err := msg.Send(text); err != nil {
		log.Printf("send error: %s", err.Error())
	}
}
