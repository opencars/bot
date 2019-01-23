// Copyright (C) 2019 Ali Shanaakh, github@shanaakh.pro
// This software may be modified and distributed under the terms
// of the MIT license. See the LICENSE file for details.

package app

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/shal/robot/pkg/openalpr"
	"html/template"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/shal/robot/internal/subscription"
	"github.com/shal/robot/pkg/autoria"
	"github.com/shal/robot/pkg/env"
)

var (
	_once sync.Once
	_app  *App
)

const (
	sleepTime = time.Hour
)

type App struct {
	Bot        *tgbotapi.BotAPI
	Recognizer *openalpr.API
	Subs       map[int64]*subscription.Subscription
	FilePath   string
}

func NewApp(path, url string) *App {
	_once.Do(func() {
		_app = &App{
			Bot:        newBot(),
			Recognizer: &openalpr.API{URL: url},
			Subs:       make(map[int64]*subscription.Subscription),
			FilePath:   path,
		}
	})

	return _app
}

func newBot() *tgbotapi.BotAPI {
	bot, err := tgbotapi.NewBotAPI(env.MustGet("BOT_TOKEN"))

	bot.Debug = env.Get("LOG_LEVEL", "DEBUG") == "DEBUG"

	if err != nil {
		panic(err)
	}

	log.Printf("Bot authorized %s\n", bot.Self.UserName)

	return bot
}

func (app *App) HandleUpdate(u tgbotapi.Update) {
	if u.Message != nil {
		app.HandleMsg(u.Message)
	}
}

func (app *App) HandleMsg(m *tgbotapi.Message) {
	log.Printf("Recieved message from [%d] %s\n", m.Chat.ID, m.Chat.UserName)
	log.Println(m.Text)

	if isFollowCmd(m.Text) {
		app.HandleFollowing(m)
	} else if isStopCmd(m.Text) {
		if _, ok := app.Subs[m.Chat.ID]; ok {
			app.Subs[m.Chat.ID].Stop()

			app.UpdateData()
		} else {
			app.SendErrorMsg(m.Chat, "You are not subscribed to updates")
		}
	} else if isAutoInfoCmd(m.Text) {
		app.HandleCarInfo(m)
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

func isAutoInfoCmd(lexeme string) bool {
	return strings.HasPrefix(lexeme, "/auto_")
}

func (app *App) HandleCarInfo(m *tgbotapi.Message) {
	lexemes := strings.Split(m.Text, "_")

	if len(lexemes) < 2 {
		app.SendErrorMsg(m.Chat, "Something wrong with command argument")
		return
	}

	carID := lexemes[1]

	autoRia := autoria.NewAPI(env.MustGet("RIA_API_KEY"))
	resp, err := autoRia.CarPhotos(carID)

	if err != nil {
		app.SendErrorMsg(m.Chat, "Seems like ID is wrong")
		return
	}

	for _, photo := range resp.Photos {
		resp, err := app.Recognizer.Recognize(photo.URL())

		if err != nil {
			log.Println(err)
			continue
		}

		plate, err := resp.Plate()

		if err == nil {
			msg := tgbotapi.NewMessage(m.Chat.ID, "Номер: " + plate)
			if _, err := app.Bot.Send(msg); err != nil {
				log.Println(err)
			} else {
				log.Printf("Successfully delivered to chat: %d\n", m.Chat.ID)
			}
			return
		}
	}

	msg := tgbotapi.NewMessage(m.Chat.ID, "Номер не найден")
	if _, err := app.Bot.Send(msg); err != nil {
		log.Println(err)
	} else {
		log.Printf("Successfully delivered to chat: %d\n", m.Chat.ID)
	}
}

func (app *App) HandleFollowing(m *tgbotapi.Message) {
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

	autoRia := autoria.NewAPI(env.MustGet("RIA_API_KEY"))

	// Convert params to old type, because frontend and api have different types.
	params, err = autoRia.ConvertNewToOld(params)

	if err != nil {
		app.SendErrorMsg(m.Chat, err.Error())
		return
	}

	// Create subscription, if it was not created.
	if _, ok := app.Subs[m.Chat.ID]; !ok {
		app.Subs[m.Chat.ID] = subscription.NewSubscription(params)
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

		// Get list of new cars.
		newCarIDs := app.Subs[m.Chat.ID].GetNewCars(search.Result.SearchResult.Cars)
		// Store latest result.
		app.Subs[m.Chat.ID].Cars = search.Result.SearchResult.Cars

		newCars := make([]autoria.CarInfoResponse, len(newCarIDs))

		for i, ID := range newCarIDs {
			car, err := autoRia.CarInfo(ID)

			if err != nil {
				log.Println(err)
			}

			newCars[i] = *car

		}

		tpl, err := template.ParseFiles("templates/message.tpl")
		if err != nil {
			log.Println(err)
		}

		buff := bytes.Buffer{}
		if err := tpl.Execute(&buff, newCars); err != nil {
			log.Println(err)
		}

		msg := tgbotapi.NewMessage(m.Chat.ID, buff.String())
		msg.ParseMode = tgbotapi.ModeHTML
		msg.DisableWebPagePreview = true

		if _, err := app.Bot.Send(msg); err != nil {
			log.Println(err)
		} else {
			log.Printf("Successfully delivered to chat: %d\n", m.Chat.ID)
		}

		time.Sleep(sleepTime)
	})

	// Add new subscription to data file.
	app.UpdateData()
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

func (app *App) UpdateData() {
	// Update data file with subscriptions.
	file, err := os.OpenFile(app.FilePath, os.O_WRONLY|os.O_CREATE, 0644)

	file.Truncate(0)
	file.Seek(0, 0)

	if err != nil {
		log.Println(err)
	} else if err := json.NewEncoder(file).Encode(app.Subs); err != nil {
		log.Printf("Data: %s\n", err.Error())
	}
}
