package bot

import (
	"context"
	"regexp"
	"time"

	"gopkg.in/tucnak/telebot.v2"

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

func NewBot(v domain.VehicleService, , token string) (*Bot, error) {

	telebot.Webhook{
		Listen:         "",
		MaxConnections: 0,
		AllowedUpdates: nil,
		HasCustomCert:  false,
		PendingUpdates: 0,
		ErrorUnixtime:  0,
		ErrorMessage:   "",
		TLS:            nil,
		Endpoint:       nil,
	}

	bot, err := telebot.NewBot(telebot.Settings{
		URL: telebot.DefaultApiURL,

		Token: token,
		Poller: &telebot.LongPoller{
			Timeout: 10 * time.Second,
		},
	})

	if err != nil {
		return nil, err
	}

	return &Bot{
		client:  bot,
		vehicle: v,
	}, nil
}

// TODO: Add middleware for logging and saving all requests into the DB.
func (b *Bot) Process(ctx context.Context) error {
	logger.Infof("me %#v", b.client.Me)

	b.client.Handle("/start", b.StartHandler())

	b.client.Handle(telebot.OnText, func(m *telebot.Message) {
		if number.MatchString(m.Text) {
			b.FindByNumber(context.TODO(), m)
		}

		if vin.MatchString(m.Text) {
			b.FindByVIN(context.TODO(), m)
		}
	})

	if err := b.client.RemoveWebhook(); err != nil {
		return err
	}

	stop := make(chan struct{})
	go b.client.Poller.Poll(b.client, b.client.Updates, stop)

	for {
		select {
		case upd := <-b.client.Updates:
			b.client.ProcessUpdate(upd)
		case <-stop:
			return nil
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}
