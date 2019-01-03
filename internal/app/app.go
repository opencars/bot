// Copyright (C) 2019 Ali Shanaakh, github@shanaakh.pro
// This software may be modified and distributed under the terms
// of the MIT license. See the LICENSE file for details.

package app

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/shal/robot/internal/subscription"
	"github.com/shal/robot/pkg/util"
	"log"
	"strings"
	"sync"
)

var (
	_once sync.Once
	_app *App
)

type App struct {
	Bot *tgbotapi.BotAPI
	Subs map[int64]*subscription.Subscription
}

func NewApp() *App {
	_once.Do(func() {
		_app = &App{
			Bot: newBot(),
			Subs: make(map[int64]*subscription.Subscription),
		}
	})

	return _app
}

func newBot() *tgbotapi.BotAPI {
	bot, err := tgbotapi.NewBotAPI(util.MustGetEnv("BOT_TOKEN"))

	if err != nil {
		panic(err)
	}

	//bot.Debug = true

	log.Printf("Bot authorized %s", bot.Self.UserName)

	return bot
}

func (app *App) ProcessUpdate(u tgbotapi.Update) {
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
	return strings.HasPrefix(lexeme, "/start")
}

func isStopCmd(lexeme string) bool {
	return strings.HasPrefix(lexeme, "/stop")
}

func (app *App) processFollowMsg(m *tgbotapi.Message) {
	app.Subs[m.Chat.ID] = subscription.NewSubscription(m.Chat)
	app.Subs[m.Chat.ID].Start(app.Bot)
}
