// Copyright (C) 2019 Ali Shanaakh, github@shanaakh.pro
// This software may be modified and distributed under the terms
// of the MIT license. See the LICENSE file for details.

package app

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"
)

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/shal/robot/internal/subscription"
	"github.com/shal/robot/pkg/autoria"
	"github.com/shal/robot/pkg/util"
)

var (
	_once sync.Once
	_app  *App
)

const (
	sleepTime = time.Minute
)

type App struct {
	Bot      *tgbotapi.BotAPI
	Subs     map[int64]*subscription.Subscription
	FilePath string
}

func NewApp(path string) *App {
	_once.Do(func() {
		_app = &App{
			Bot:      newBot(),
			Subs:     make(map[int64]*subscription.Subscription),
			FilePath: path,
		}
	})

	return _app
}

func newBot() *tgbotapi.BotAPI {
	bot, err := tgbotapi.NewBotAPI(util.MustGetEnv("BOT_TOKEN"))

	if err != nil {
		panic(err)
	}

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
		if _, ok := app.Subs[m.Chat.ID]; ok {
			app.Subs[m.Chat.ID].Stop()
			delete(app.Subs, m.Chat.ID)

			app.UpdateDataFile()
		} else {
			app.SendErrorMsg(m.Chat, "You are not subscribed to updates")
		}
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
		app.SendErrorMsg(m.Chat, "Something wrong with command argument")
		return
	} else if !strings.HasPrefix(lexemes[1], "https://auto.ria.com/search") {
		app.SendErrorMsg(m.Chat, "Seems like link is wrong")
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

	// Create subscription, if it was not created.
	if _, ok := app.Subs[m.Chat.ID]; !ok {
		app.Subs[m.Chat.ID] = subscription.NewSubscription(m.Chat, params)
	}

	if err != nil {
		app.SendErrorMsg(m.Chat, err.Error())
		return
	}

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
				log.Println("Quit was called")
				return
			default:
				msg := tgbotapi.NewMessage(m.Chat.ID, car.LinkToView)

				if _, err := app.Bot.Send(msg); err != nil {
					log.Println(err)
				} else {
					log.Printf("Successfully delivered to chat: %d\n", m.Chat.ID)
				}

				time.Sleep(time.Second * 3)
			}
		}

		time.Sleep(sleepTime)
	})

	// Add new subscription to data file.
	app.UpdateDataFile()
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

func (app *App) UpdateDataFile() {
	var values = make([]subscription.Subscription, 0)
	for _, value := range app.Subs {
		values = append(values, *value)
	}

	// Update data file with subscriptions.
	file, err := os.OpenFile(app.FilePath, os.O_WRONLY|os.O_CREATE, 0644)

	file.Truncate(0)
	file.Seek(0,0)

	if err != nil {
		log.Println(err)
	} else if err := json.NewEncoder(file).Encode(values); err != nil {
		log.Printf("Data: %s\n", err.Error())
	}
}