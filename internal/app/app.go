// Copyright (C) 2019 Ali Shanaakh, github@shanaakh.pro
// This software may be modified and distributed under the terms
// of the MIT license. See the LICENSE file for details.

package app

import (
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/shal/robot/internal/subscription"
	"github.com/shal/robot/pkg/autoria"
	"github.com/shal/robot/pkg/util"
	"log"
	"strings"
	"sync"
	"time"
)

var (
	_once sync.Once
	_app  *App
)

const (
	sleepTime = time.Hour
)

type App struct {
	Bot  *tgbotapi.BotAPI
	Subs map[int64]*subscription.Subscription
}

func NewApp() *App {
	_once.Do(func() {
		_app = &App{
			Bot:  newBot(),
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
		//if _, ok := app.Subs[m.Chat.ID]; ok {
		app.Subs[m.Chat.ID].Stop()
		//} else {/
		//	app.SendErrorMsg(m.Chat, "You are not subscribed to updates")
		//}
	} else {
		text := fmt.Sprintf("Invalid command %s", strings.Split(m.Text, " ")[0])
		app.SendErrorMsg(m.Chat, text)
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
		text := fmt.Sprintf("Something wrong with command argument")
		app.SendErrorMsg(m.Chat, text)
		return
	} else if !strings.HasPrefix(lexemes[1], "https://auto.ria.com/search") {
		text := fmt.Sprintf("Seems like link is wrong")
		app.SendErrorMsg(m.Chat, text)
		return
	}

	params, err := autoria.ParseCarSearchParams(lexemes[1])

	if err != nil {
		app.SendErrorMsg(m.Chat, err.Error())
		return
	}

	autoRia := autoria.NewAPI(util.MustGetEnv("RIA_API_KEY"))

	// Convert params to old type, because frontend and api have different types.
	params, err = autoRia.ConvertNewToOld(params)

	if err != nil {
		app.SendErrorMsg(m.Chat, err.Error())
		return
	}

	app.Subs[m.Chat.ID] = subscription.NewSubscription(m.Chat, params)
	app.Subs[m.Chat.ID].Start(func(quitter chan struct{}) {
		search, err := autoRia.GetSearchCars(params)

		if err != nil {
			app.SendErrorMsg(m.Chat, err.Error())
			return
		}

		for _, ID := range search.Result.SearchResult.CarsIDs {
			car, err := autoRia.GetCarInfo(ID)

			if err != nil {
				app.SendErrorMsg(m.Chat, err.Error())
				return
			}

			select {
			case <-quitter:
				fmt.Println("Quit was called. Test 2")
				return
			default:
				msg := tgbotapi.NewMessage(m.Chat.ID, car.LinkToView)

				if _, err := app.Bot.Send(msg); err != nil {
					log.Print(err)
				} else {
					log.Printf("Successfully delivered to chat: %d\n", m.Chat.ID)
				}

				time.Sleep(time.Second * 3)
			}
		}
	})
}
func (app *App) SendErrorMsg(chat *tgbotapi.Chat, text string) {
	msg := tgbotapi.NewMessage(chat.ID, text)
	log.Println(text)

	if _, err := app.Bot.Send(msg); err != nil {
		log.Println(err)
	} else {
		log.Printf("Successfully delivered to chat: %d\n", chat.ID)
	}
}
