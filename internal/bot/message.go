package bot

import "github.com/go-telegram-bot-api/telegram-bot-api"

//
type Message struct {
	API *tgbotapi.BotAPI
	msg *tgbotapi.Message
}

//
func (msg *Message) Chat() *tgbotapi.Chat {
	return msg.msg.Chat
}

//
func (msg *Message) Text() string {
	return msg.msg.Text
}

//
func (msg *Message) Photo() *[]tgbotapi.PhotoSize {
	return msg.msg.Photo
}

//
type Handler interface {
	Handle(msg *Message)
}

func (msg *Message) send(message tgbotapi.MessageConfig) error {
	message.DisableWebPagePreview = !WebPagePreview
	if _, err := msg.API.Send(message); err != nil {
		return err
	}

	return nil
}

//
func (msg *Message) Send(text string) error {
	return msg.send(tgbotapi.NewMessage(msg.msg.Chat.ID, text))
}

//
func (msg *Message) SetStatus(status string) error {
	action := tgbotapi.NewChatAction(msg.msg.Chat.ID, status)
	if _, err := msg.API.Send(action); err != nil {
		return err
	}

	return nil
}

// SendHTML sends message to the chat with text formatted as HTML.
func (msg *Message) SendHTML(text string) error {
	res := tgbotapi.NewMessage(msg.msg.Chat.ID, text)
	res.ParseMode = tgbotapi.ModeHTML

	return msg.send(res)
}
