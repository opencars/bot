package bot

import (
	"context"

	"gopkg.in/tucnak/telebot.v2"

	"github.com/opencars/bot/pkg/domain"
	"github.com/opencars/bot/pkg/domain/model"
	"github.com/opencars/seedwork/logger"
)

type Poller struct {
	repo domain.MessageService
	next telebot.Poller
}

func NewPoller(r domain.MessageRepository, next telebot.Poller) *Poller {
	return &Poller{
		repo: r,
		next: next,
	}
}

func (p *Poller) Poll(b *telebot.Bot, updates chan telebot.Update, stop chan struct{}) {
	middle := make(chan telebot.Update, 1)
	stopPoller := make(chan struct{})

	go p.next.Poll(b, middle, stop)

	for {
		select {
		case <-stop:
			close(stopPoller)
			return
		case upd := <-middle:
			logger.WithFields(logger.Fields{
				"id":        upd.Message.ID,
				"chat_id":   upd.Message.Chat.ID,
				"chat_type": upd.Message.Chat.Type,
				"text":      upd.Message.Text,
				"time":      upd.Message.Time(),
			}).Infof("incoming message")

			msg := model.Message{
				ID: upd.Message.ID,
				User: model.User{
					ID:        int(upd.Message.Chat.ID),
					FirstName: upd.Message.Chat.FirstName,
				},
				Text: upd.Message.Text,
				Time: upd.Message.Time(),
			}

			if upd.Message.Chat.LastName != "" {
				msg.User.LastName = &upd.Message.Chat.LastName
			}

			if upd.Message.Chat.Username != "" {
				msg.User.UserName = &upd.Message.Chat.Username
			}

			if err := p.repo.Create(context.Background(), &msg); err != nil {
				logger.Errorf("message: %s", err)
			}

			updates <- upd
		}
	}
}
