// Copyright (C) 2019 Ali Shanaakh, github@shanaakh.pro
// This software may be modified and distributed under the terms
// of the MIT license. See the LICENSE file for details.

package app

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/shal/robot/internal/subscription"
	"github.com/shal/robot/pkg/autoria-api"
	"github.com/shal/robot/pkg/util"
	"log"
	"strings"
	"sync"
	"time"
)

var (
	_once sync.Once
	_app *App
)

const (
	sleepTime = time.Second * 30
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

	log.Printf("Bot authorized %s\n", bot.Self.UserName)

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
	} else {
		// TODO: Answer that command is invalid.
	}
}

func isFollowCmd(lexeme string) bool {
	return strings.HasPrefix(lexeme, "/follow")
}

func isStopCmd(lexeme string) bool {
	return strings.HasPrefix(lexeme, "/stop")
}

func (app *App) processFollowMsg(m *tgbotapi.Message) {
	lexemes := strings.Split(m.Text, " ")

	if len(lexemes) < 2 {
		// TODO: Show error.
		return
	} else if !strings.HasPrefix(lexemes[1], "https://auto.ria.com/search") {
		// TODO: Show error.
		return
	}

	params, err := autoria_api.ParseCarSearchParams(lexemes[1])

	if err != nil {
		// TODO: Show error.
		log.Print(err)
		return
	}

	autoRia := autoria_api.NewAPI(util.MustGetEnv("RIA_API_KEY"))

	// Convert params to old type, because frontend and api have different types.
	params, err = autoRia.ConvertNewToOld(params)

	if err != nil {
		// TODO: Show error.
		log.Print(err)
		return
	}

	app.Subs[m.Chat.ID] = subscription.NewSubscription(m.Chat, params)
	app.Subs[m.Chat.ID].Start(func() {
		chat := app.Subs[m.Chat.ID].Chat

		search, err := autoRia.GetSearchCars(params)

		if err != nil {
			// TODO: Show error.
			log.Print(err)
			return
		}

		for _, ID := range search.Result.SearchResult.CarsIDs {
			car, err := autoRia.GetCarInfo(ID)

			if err != nil {
				// TODO: Show error.
				log.Print(err)
				return
			}

			msg := tgbotapi.NewMessage(chat.ID, car.LinkToView)

			if _, err := app.Bot.Send(msg); err != nil {
				log.Print(err)
			} else {
				log.Printf("Successfully delivered to Chat: %d\n", chat.ID)
			}

			time.Sleep(time.Second * 3)
		}

		time.Sleep(sleepTime)
	})
}
