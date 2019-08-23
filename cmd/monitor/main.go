package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"os"
	"sort"
	"text/template"
	"time"

	"github.com/go-telegram-bot-api/telegram-bot-api"

	"github.com/opencars/bot/pkg/gov"
)

// Sorter sorts array by LastModified field.
type Sorter []gov.Resource

func (e Sorter) Len() int      { return len(e) }
func (e Sorter) Swap(i, j int) { e[i], e[j] = e[j], e[i] }
func (e Sorter) Less(i, j int) bool {
	return e[i].LastModified.Before(e[j].LastModified.Time)
}

var (
	pkg    = flag.String("package", "", "Package ID on https://data.gov.ua")
	user   = flag.Int64("user", 0, "Telegram user to be notified")
	token  = flag.String("token", "", "Telegram bot token")
	period = flag.Duration("period", 10*time.Minute, "")
)

// Monitor is responsible from monitoring government data registry.
type Monitor struct {
	timestamp gov.Time
	client    *gov.Client
	bot       *tgbotapi.BotAPI
}

// New creates new instance of Monitor.
func New() (*Monitor, error) {
	bot, err := tgbotapi.NewBotAPI(*token)
	if err != nil {
		return nil, err
	}

	return &Monitor{
		// Month ago.
		timestamp: gov.Time{Time: time.Now().Add(- 730 * time.Hour)},
		client:    gov.NewClient(),
		bot:       bot,
	}, nil
}

func (monitor *Monitor) notify(pkg *gov.Package) error {
	tpl, err := template.ParseFiles("templates/monitor.tpl")
	if err != nil {
		return err
	}

	buff := bytes.Buffer{}
	if err := tpl.Execute(&buff, struct {
		Package  string
		Resource gov.Resource
	}{
		pkg.Title, pkg.Resources[len(pkg.Resources)-1],
	}); err != nil {
		return err
	}

	msg := tgbotapi.NewMessage(*user, buff.String())
	msg.ParseMode = tgbotapi.ModeHTML

	_, err = monitor.bot.Send(msg)
	if err != nil {
		return err
	}
	return nil
}

func (monitor *Monitor) check() error {
	res, err := monitor.client.Package(*pkg)
	if err != nil {
		return err
	}

	// Sort by modification time.
	sort.Sort(Sorter(res.Resources))
	resource := res.Resources[len(res.Resources)-1]

	if resource.LastModified.After(monitor.timestamp.Time) {
		fmt.Println("Package was modified!")

		// Update latest timestamp.
		monitor.timestamp = resource.LastModified

		err := monitor.notify(res)
		if err != nil {
			return err
		}
	}

	return nil
}

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

func main() {
	flag.Parse()

	if err := validate(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	monitor, err := New()
	if err != nil {
		panic(err)
	}

	for {
		err := monitor.check()
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error skipped!", err)
			continue
		}

		<-time.After(*period)
	}
}
