// Copyright (C) 2019 Ali Shanaakh, github@shanaakh.pro
// This software may be modified and distributed under the terms
// of the MIT license. See the LICENSE file for details.

package bot

import (
	"fmt"
	"log"
	"strings"
	"sync"
	"time"
)

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/shal/robot/pkg/util"
)

var (
	_once sync.Once
	_app *App
)

type Subscription struct {
	ChatID int64
	LastIDs []string
	Quitter chan struct{}
	Running bool
}

// App is a singleton representation of application.
type App struct {
	Bot *tgbotapi.BotAPI
	Subs map[int64]Subscription
}

func NewApp() *App {
	_once.Do(func() {
		_app = &App{
			Bot: newBot(),
			Subs: make(map[int64]Subscription),
		}
	})

	return _app
}

func NewSubscription(ID int64) Subscription {
	lastIDs := make([]string, 0)
	stopper := make(chan struct{})
	return Subscription{ID, lastIDs, stopper, false}
}

func newBot() *tgbotapi.BotAPI {
	bot, err := tgbotapi.NewBotAPI(util.MustGetEnv("BOT_TOKEN"))

	if err != nil {
		panic(err)
	}

	bot.Debug = true
	log.Printf("Bot authorized %s", bot.Self.UserName)

	return bot
}

func (app *App) processUpdate(u tgbotapi.Update) {
	if u.Message != nil {
		app.processMsgUpdate(u.Message)
	}
}

func (app *App) processMsgUpdate(m *tgbotapi.Message) {
	log.Printf("Recieved message from [%d] %s\n", m.Chat.ID, m.Chat.UserName)
	log.Println(m.Text)

	if isFollowCmd(m.Text) {
		app.processFollowMsg(m)
	} else if isStopCmd(m.Text) {
		app.Subs[m.Chat.ID].Stop()
	}
}

func isFollowCmd(lexeme string) bool {
	return strings.HasPrefix(lexeme, "/follow")
}

func isStopCmd(lexeme string) bool {
	return strings.HasPrefix(lexeme, "/stop")
}

func (app *App) processFollowMsg(m *tgbotapi.Message) {
	app.Subs[m.Chat.ID] = NewSubscription(m.Chat.ID)
	app.Subs[m.Chat.ID].Run()
}

func (s Subscription) Stop() {
	// Avoid closing closed channel.
	if !s.Running {
		return
	}

	close(s.Quitter)
	s.Running = false
}

func (s Subscription) Run() {
	// Stop previous goroutine, if it is running.
	if s.Running {
		close(s.Quitter)
		s.Quitter = make(chan struct{})
	}

	go func(q chan struct{}) {
		for {
			select {
			case <-q:
				return
			default:
				fmt.Print("Hi --------> ", s.ChatID)
				time.Sleep(time.Second * 3)
			}
		}
	}(s.Quitter)

	// Mark subscription as "Running".
	s.Running = true
}

func Run() {
	app := NewApp()

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60

	updates, err := app.Bot.GetUpdatesChan(updateConfig)

	if err != nil {
		log.Println(err.Error())
	}

	for update := range updates {
		app.processUpdate(update)
	}
}
