package bot

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

//
type Event struct {
	API     *tgbotapi.BotAPI
	Message *tgbotapi.Message
}

//
type Handler interface {
	Handle(msg *Event)
}

func (event *Event) send(message tgbotapi.MessageConfig) error {
	message.DisableWebPagePreview = !WebPagePreview
	if _, err := event.API.Send(message); err != nil {
		return err
	}

	return nil
}

//
func (event *Event) Send(text string) error {
	return event.send(tgbotapi.NewMessage(event.Message.Chat.ID, text))
}

//
func (event *Event) SetStatus(status string) error {
	action := tgbotapi.NewChatAction(event.Message.Chat.ID, status)
	if _, err := event.API.Send(action); err != nil {
		return err
	}

	return nil
}

// SendHTML sends message to the chat with text formatted as HTML.
func (event *Event) SendHTML(text string) error {
	res := tgbotapi.NewMessage(event.Message.Chat.ID, text)
	res.ParseMode = tgbotapi.ModeHTML

	return event.send(res)
}
