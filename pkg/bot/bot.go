package bot

import (
	"context"
	"regexp"

	"gopkg.in/tucnak/telebot.v2"

	"github.com/opencars/bot/pkg/config"
	"github.com/opencars/bot/pkg/domain"
	"github.com/opencars/bot/pkg/logger"
)

var (
	number = regexp.MustCompile(`^\p{L}{2}\d{4}\p{L}{2}$`)
	vin    = regexp.MustCompile(`^[A-HJ-NPR-Z0-9]{17}$`)
)

type Bot struct {
	client *telebot.Bot

	vehicle domain.VehicleService
}

func NewBot(conf *config.Bot, v domain.VehicleService, m domain.MessageService, addr string) (*Bot, error) {
	w := telebot.Webhook{
		Listen: addr,
		Endpoint: &telebot.WebhookEndpoint{
			PublicURL: conf.URL,
		},
	}

	if conf.MaxConnections != nil {
		w.MaxConnections = *conf.MaxConnections
	}

	bot, err := telebot.NewBot(telebot.Settings{
		URL:    telebot.DefaultApiURL,
		Token:  conf.Token,
		Poller: NewPoller(m, &w),
	})
	if err != nil {
		return nil, err
	}

	return &Bot{
		client:  bot,
		vehicle: v,
	}, nil
}

func (b *Bot) Process(ctx context.Context) error {
	logger.Infof("Logged in as %s", b.client.Me.Username)

	b.client.Handle(telebot.OnText, func(m *telebot.Message) {
		if number.MatchString(m.Text) {
			if err := b.FindByNumber(context.TODO(), m); err != nil {
				logger.Errorf("number: %s", err)
			}
		}

		if vin.MatchString(m.Text) {
			if err := b.FindByVIN(context.TODO(), m); err != nil {
				logger.Errorf("vin: %s", err)
			}
		}
	})

	stop := make(chan struct{})
	go b.client.Poller.Poll(b.client, b.client.Updates, stop)

	for {
		select {
		case upd := <-b.client.Updates:
			b.client.ProcessUpdate(upd)
		case <-stop:
			return nil
		case <-ctx.Done():
			return nil
		}
	}
}
