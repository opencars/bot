package handlers

import "github.com/go-telegram-bot-api/telegram-bot-api"

func send(bot *tgbotapi.BotAPI, message tgbotapi.MessageConfig) error {
	if _, err := bot.Send(message); err != nil {
		return err
	}

	return nil
}

// SendHTML send message to the chat with regular text.
func Send(bot *tgbotapi.BotAPI, chat *tgbotapi.Chat, text string) error {
	msg := tgbotapi.NewMessage(chat.ID, text)
	return send(bot, msg)
}

// SendHTML send message to the chat with text formatted as HTML.
func SendHTML(bot *tgbotapi.BotAPI, chat *tgbotapi.Chat, text string) error {
	msg := tgbotapi.NewMessage(chat.ID, text)
	msg.ParseMode = tgbotapi.ModeHTML

	return send(bot, msg)
}

// SendHTML send message to the chat with text formatted as HTML.
func SendMsgHTML(msg tgbotapi.MessageConfig, bot *tgbotapi.BotAPI) error {
	msg.ParseMode = tgbotapi.ModeHTML

	return send(bot, msg)
}

// SendMarkdown send message to the chat with text formatted as Markdown.
func SendMarkdown(bot *tgbotapi.BotAPI, chat *tgbotapi.Chat, text string) error {
	msg := tgbotapi.NewMessage(chat.ID, text)
	msg.ParseMode = tgbotapi.ModeMarkdown

	return send(bot, msg)
}
