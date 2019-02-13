package handlers

import (
	"bytes"
	"fmt"
	"github.com/shal/robot/internal/bot"
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
func (h AutoRiaHandler) FollowHandler(api *tgbotapi.BotAPI, msg *tgbotapi.Message) {
	lexemes := strings.Split(msg.Text, " ")

	if len(lexemes) < 2 || !strings.HasPrefix(lexemes[1], "https://auto.ria.com/search") {
		if err := bot.Send(api, msg.Chat, "Помилковий запит."); err != nil {
			log.Printf("send error: %s\n", err.Error())
		}
		return
	}

	params, err := autoria.ParseSearchParams(lexemes[1])
	if err != nil {
		if err := bot.Send(api, msg.Chat, err.Error()); err != nil {
			log.Printf("send error: %s\n", err.Error())
		}
		return
	}

	// Convert params to old type, because frontend and api have different types.
	params, err = h.API.ConvertNewToOld(params)
	if err != nil {
		if err := bot.Send(api, msg.Chat, err.Error()); err != nil {
			log.Printf("send error: %s\n", err.Error())
		}
		return
	}

	// Create subscription, if it was not created.
	if _, ok := h.Subscriptions[msg.Chat.ID]; !ok {
		h.Subscriptions[msg.Chat.ID] = subscription.New(params)
	}

	h.Subscriptions[msg.Chat.ID].Start(func(quitter chan struct{}) {
		search, err := h.API.SearchCars(params)

		if err != nil {
			if err := bot.Send(api, msg.Chat, err.Error()); err != nil {
				log.Printf("send error: %s\n", err.Error())
			}
			return
		}

		// Get list of new cars.
		newCarIDs := h.Subscriptions[msg.Chat.ID].NewCars(search.Result.SearchResult.Cars)
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

		if err := bot.SendMsgHTML(msg, api); err != nil {
			log.Printf("send error: %s", err.Error())
		}

		time.Sleep(time.Hour)
	})

	// TODO: Save changes to file with data.
	// Add new subscription to data file.
	//api.UpdateData()
}

func (h AutoRiaHandler) StopHandler(api *tgbotapi.BotAPI, msg *tgbotapi.Message) {
	if _, ok := h.Subscriptions[msg.Chat.ID]; !ok {
		if err := bot.Send(api, msg.Chat, "Ви не підписані на оновлення 🤔"); err != nil {
			log.Printf("send error: %s", err.Error())
		}
		return
	}

	h.Subscriptions[msg.Chat.ID].Stop()
}

// TODO: Refactor this handler.
// Analyze first 50 photos, then find best number, that matches the rules.
// Send message firstly.
func (h AutoRiaHandler) CarInfoHandler(api *tgbotapi.BotAPI, msg *tgbotapi.Message) {
	lexemes := strings.Split(msg.Text, "_")

	if len(lexemes) < 2 {
		if err := bot.Send(api, msg.Chat, "Помилковий запит 😮"); err != nil {
			log.Printf("send error: %s", err.Error())
		}
		return
	}

	carID := lexemes[1]

	autoRiaToken := env.MustGet("AUTO_RIA_TOKEN")
	autoRia := autoria.New(autoRiaToken)
	resp, err := autoRia.CarPhotos(carID)

	if err != nil {
		if err := bot.Send(api, msg.Chat, "Неправильний ідентифікатор 🙄️"); err != nil {
			log.Printf("send error: %s", err.Error())
		}
		return
	}

	// Get user know about waiting time.
	text := "Аналіз може зайняти до 1 хвилини 🐌"
	if err := bot.SendHTML(api, msg.Chat, text); err != nil {
		log.Printf("send error: %s\n", err.Error())
	}

	for _, photo := range resp.Photos {
		resp, err := h.Recognizer.Recognize(photo.URL())

		if err != nil {
			log.Println(err)
			continue
		}

		plate, err := resp.Plate()

		if err != nil {
			continue
		}

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

		if err := bot.SendHTML(api, msg.Chat, buff.String()); err != nil {
			log.Printf("send error: %s\n", err.Error())
		}

		return
	}

	if err := bot.Send(api, msg.Chat, "Вибачте, номер не знайдено 😳"); err != nil {
		log.Printf("send error: %s\n", err.Error())
	}
}
