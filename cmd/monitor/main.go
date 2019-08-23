package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"html/template"
	"os"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/opencars/bot/pkg/gov"

	"github.com/opencars/bot/pkg/monitor"
)

var (
	pkg    = flag.String("package", "", "Package ID on https://data.gov.ua")
	user   = flag.Int64("user", 0, "Telegram user to be notified")
	token  = flag.String("token", "", "Telegram bot token")
	period = flag.Duration("period", 10*time.Minute, "")
)

func validate() error {
	if *pkg == "" {
		return errors.New("package can not be empty")
	}

	if *user == 0 {
		return errors.New("user can not be empty")
	}

	if *token == "" {
		return errors.New("token can not be empty")
	}

	return nil
}

type Notifier struct {
	bot *tgbotapi.BotAPI
	tpl *template.Template
}

func NewNotifier() (*Notifier, error) {
	bot, err := tgbotapi.NewBotAPI(*token)
	if err != nil {
		return nil, err
	}

	tpl, err := template.ParseFiles("templates/monitor.tpl")
	if err != nil {
		return nil, err
	}

	return &Notifier{
		bot: bot,
		tpl: tpl,
	}, nil
}

func (n *Notifier) notify(pkg *gov.Package) error {
	buff := bytes.Buffer{}
	if err := n.tpl.Execute(&buff, struct {
		Package  string
		Resource gov.Resource
	}{
		pkg.Title, pkg.Resources[len(pkg.Resources)-1],
	}); err != nil {
		return err
	}

	msg := tgbotapi.NewMessage(*user, buff.String())
	msg.ParseMode = tgbotapi.ModeHTML

	if _, err := n.bot.Send(msg); err != nil {
		return err
	}
	return nil
}

func main() {
	flag.Parse()

	telegram, err := NewNotifier()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	if err := validate(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	m, err := monitor.New(*pkg, *user, gov.NewClient())
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	go func() {
		for event := range m.Events {
			fmt.Println("Package was modified!")
			if err := telegram.notify(event); err != nil {
				fmt.Fprintln(os.Stderr, "Error skipped!", err)
				continue
			}

			<-time.After(*period)
		}
	}()

	if err := m.Listen(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
