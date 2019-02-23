package handlers

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/url"
	"strings"
	"time"

	"github.com/shal/opencars-bot/internal/bot"
	"github.com/shal/opencars-bot/internal/subscription"
	"github.com/shal/opencars-bot/pkg/autoria"
	"github.com/shal/opencars-bot/pkg/env"
	"github.com/shal/opencars-bot/pkg/openalpr"
	"github.com/shal/opencars-bot/pkg/opencars"
)

type AutoRiaHandler struct {
	API        *autoria.API
	Recognizer *openalpr.API
	Storage    *opencars.API

	Subscriptions map[int64]*subscription.Subscription
	FilePath      string
}

func (h AutoRiaHandler) FollowHandler(msg *bot.Message) {
	if err := msg.SetStatus(bot.ChatTyping); err != nil {
		log.Printf("action error: %s", err.Error())
	}

	lexemes := strings.Split(msg.Text(), " ")
	if len(lexemes) < 2 || !strings.HasPrefix(lexemes[1], "https://auto.ria.com/search") {
		if err := msg.Send("Помилковий запит."); err != nil {
			log.Printf("send error: %s\n", err.Error())
		}
		return
	}

	values, err := url.ParseQuery(lexemes[1])
	if err != nil {
		if err := msg.Send(err.Error()); err != nil {
			log.Printf("send error: %s\n", err.Error())
		}
		return
	}

	// Convert params to old type, because frontend and api have different types.
	values, err = h.API.ConvertNewToOld(values)
	if err != nil {
		if err := msg.Send(err.Error()); err != nil {
			log.Printf("send error: %s\n", err.Error())
		}
		return
	}

	// Create subscription, if it was not created.
	if _, ok := h.Subscriptions[msg.Chat().ID]; !ok {
		h.Subscriptions[msg.Chat().ID] = subscription.New()
	}

	h.Subscriptions[msg.Chat().ID].Start(func(quitter chan struct{}) {
		search, err := h.API.SearchCars(values)

		if err != nil {
			if err := msg.Send(err.Error()); err != nil {
				log.Printf("send error: %s\n", err.Error())
			}
			return
		}

		// Fetch list of new cars.
		newCarIDs := h.Subscriptions[msg.Chat().ID].NewCars(search.Result.SearchResult.Cars)
		// Store latest result.
		h.Subscriptions[msg.Chat().ID].Cars = search.Result.SearchResult.Cars

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

		if err := msg.SendHTML(buff.String()); err != nil {
			log.Printf("send error: %s", err.Error())
		}

		time.Sleep(time.Hour)
	})

	// TODO: Save changes to file with data.
	// Add new subscription to data file.
	//api.UpdateData()
}

func (h AutoRiaHandler) StopHandler(msg *bot.Message) {
	if err := msg.SetStatus(bot.ChatTyping); err != nil {
		log.Printf("action error: %s", err.Error())
	}

	if _, ok := h.Subscriptions[msg.Chat().ID]; !ok {
		if err := msg.Send("Ви не підписані на оновлення 🤔"); err != nil {
			log.Printf("send error: %s", err.Error())
		}
		return
	}

	h.Subscriptions[msg.Chat().ID].Stop()
}

// TODO: Refactor this handler.
// Analyze first 50 photos, then find best number, that matches the rules.
// Send message firstly.
func (h AutoRiaHandler) CarInfoHandler(msg *bot.Message) {
	lexemes := strings.Split(msg.Text(), "_")

	if len(lexemes) < 2 {
		if err := msg.Send("Помилковий запит 😮"); err != nil {
			log.Printf("send error: %s", err.Error())
		}
		return
	}

	carID := lexemes[1]

	autoRiaToken := env.MustFetch("AUTO_RIA_TOKEN")
	autoRia := autoria.New(autoRiaToken)
	resp, err := autoRia.CarPhotos(carID)

	if err != nil {
		if err := msg.Send("Неправильний ідентифікатор 🙄️"); err != nil {
			log.Printf("send error: %s", err.Error())
		}
		return
	}

	// Fetch user know about waiting time.
	text := "Аналіз може зайняти до 1 хвилини 🐌"
	if err := msg.Send(text); err != nil {
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
		if err != nil {
			log.Println(err)
			return
		}

		tpl, err := template.ParseFiles("templates/car_info.tpl")
		if err != nil {
			log.Println(err)
			return
		}

		buff := bytes.Buffer{}
		if err := tpl.Execute(&buff, struct {
			Cars   []opencars.Transport
			Number string
		}{
			transport, plate,
		}); err != nil {
			log.Println(err)
			return
		}

		if err := msg.Send(buff.String()); err != nil {
			log.Printf("send error: %s\n", err.Error())
		}

		return
	}

	if err := msg.Send("Вибачте, номер не знайдено 😳"); err != nil {
		log.Printf("send error: %s\n", err.Error())
	}
}
