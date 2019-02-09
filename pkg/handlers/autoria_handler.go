package handlers

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"strings"
	"time"

	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/shal/robot/internal/subscription"
	"github.com/shal/robot/pkg/autoria"
	"github.com/shal/robot/pkg/env"
	"github.com/shal/robot/pkg/openalpr"
	"github.com/shal/robot/pkg/opencars"
)

type AutoRiaHandler struct {
	API        *autoria.API
	Recognizer *openalpr.API
	Storage    *opencars.API

	Subscriptions map[int64]*subscription.Subscription
	FilePath      string
}

// TODO: Split this method into few methods aka delegate code.
func (h AutoRiaHandler) FollowHandler(bot *tgbotapi.BotAPI, msg *tgbotapi.Message) {
	lexemes := strings.Split(msg.Text, " ")

	if len(lexemes) < 2 || !strings.HasPrefix(lexemes[1], "https://auto.ria.com/search") {
		Send(bot, msg.Chat, "Помилковий запит.")
		return
	}

	params, err := autoria.ParseCarSearchParams(lexemes[1])
	if err != nil {
		Send(bot, msg.Chat, err.Error())
		return
	}

	// Convert params to old type, because frontend and api have different types.
	params, err = h.API.ConvertNewToOld(params)
	if err != nil {
		Send(bot, msg.Chat, err.Error())
		return
	}

	// Create subscription, if it was not created.
	if _, ok := h.Subscriptions[msg.Chat.ID]; !ok {
		h.Subscriptions[msg.Chat.ID] = subscription.NewSubscription(params)
	}

	h.Subscriptions[msg.Chat.ID].Start(func(quitter chan struct{}) {
		search, err := h.API.GetSearchCars(params)

		if err != nil {
			Send(bot, msg.Chat, err.Error())
			return
		}

		// Get list of new cars.
		newCarIDs := h.Subscriptions[msg.Chat.ID].GetNewCars(search.Result.SearchResult.Cars)
		// Store latest result.
		h.Subscriptions[msg.Chat.ID].Cars = search.Result.SearchResult.Cars

		newCars := make([]autoria.CarInfo, len(newCarIDs))

		for i, ID := range newCarIDs {
			car, err := h.API.CarInfo(ID)

			if err != nil {
				log.Println(err)
			}

			fmt.Println(car)
			fmt.Println(newCars)

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

		msg := tgbotapi.NewMessage(msg.Chat.ID, buff.String())
		msg.DisableWebPagePreview = true
		SendMsgHTML(msg, bot)

		time.Sleep(time.Hour)
	})

	// TODO: Save changes to file with data.
	// Add new subscription to data file.
	//bot.UpdateData()
}

func (h AutoRiaHandler) StopHandler(bot *tgbotapi.BotAPI, msg *tgbotapi.Message) {
	if _, ok := h.Subscriptions[msg.Chat.ID]; !ok {
		Send(bot, msg.Chat, "Ви не підписані на оновлення.")
		return
	}

	h.Subscriptions[msg.Chat.ID].Stop()
}

func (h AutoRiaHandler) CarInfoHandler(bot *tgbotapi.BotAPI, msg *tgbotapi.Message) {
	lexemes := strings.Split(msg.Text, "_")

	if len(lexemes) < 2 {
		Send(bot, msg.Chat, "Something wrong with command argument")
		return
	}

	carID := lexemes[1]

	autoRia := autoria.NewAPI(env.MustGet("RIA_API_KEY"))
	resp, err := autoRia.CarPhotos(carID)

	if err != nil {
		Send(bot, msg.Chat, "Неправильний ідентифікатор ☹️")
		return
	}

	for _, photo := range resp.Photos {
		resp, err := h.Recognizer.Recognize(photo.URL())

		if err != nil {
			log.Println(err)
			continue
		}

		plate, err := resp.Plate()
		if err == nil {
			transport, err := h.Storage.Search(plate)

			fmt.Println(transport)

			if err != nil {
				log.Println(err)
			}

			tpl, err := template.ParseFiles("templates/car_info.tpl")
			if err != nil {
				log.Println(err)
			}

			buff := bytes.Buffer{}
			if err := tpl.Execute(&buff, struct {
				Cars   []opencars.Transport
				Number string
			}{
				transport, plate,
			}); err != nil {
				log.Println(err)
			}

			SendHTML(bot, msg.Chat, buff.String())

			return
		}
	}

	Send(bot, msg.Chat, "Номер не знайдено!")
}
