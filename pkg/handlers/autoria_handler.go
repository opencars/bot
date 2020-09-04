package handlers

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/opencars/bot/internal/bot"
	"github.com/opencars/bot/internal/subscription"
	"github.com/opencars/bot/pkg/autoria"
	"github.com/opencars/toolkit"
)

type AutoRiaHandler struct {
	API           *autoria.API
	Period        time.Duration
	Toolkit       *toolkit.Client
	Subscriptions map[int64]*subscription.Subscription
}

func (h AutoRiaHandler) FollowHandler(msg *bot.Event) {
	if err := msg.SetStatus(bot.ChatTyping); err != nil {
		log.Printf("action error: %s", err.Error())
	}

	if !strings.HasPrefix(msg.Message.Text, "https://auto.ria.com/search") {
		if err := msg.Send("Помилковий запит."); err != nil {
			log.Printf("send error: %s\n", err.Error())
		}
		return
	}

	values, err := url.ParseQuery(msg.Message.Text)
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
	if _, ok := h.Subscriptions[msg.Message.Chat.ID]; !ok {
		h.Subscriptions[msg.Message.Chat.ID] = subscription.New(h.Period)
	}

	h.Subscriptions[msg.Message.Chat.ID].Start(func() {
		search, err := h.API.SearchCars(values)

		if err != nil {
			log.Printf("Failed to search cars: %s\n", err)
			return
		}

		// Fetch list of new cars.
		newCarIDs := h.Subscriptions[msg.Message.Chat.ID].NewCars(search.Result.SearchResult.Cars)

		// Skip, if newCarIDs is empty.
		if len(newCarIDs) == 0 {
			return
		}

		// Store latest result.
		h.Subscriptions[msg.Message.Chat.ID].Cars = search.Result.SearchResult.Cars

		newCars := make([]autoria.CarInfo, len(newCarIDs))
		for i, ID := range newCarIDs {
			car, err := h.API.CarInfo(ID)
			if err != nil {
				log.Printf("Failed to get car info: %s\n", err)
				return
			}

			newCars[i] = *car
		}

		tpl, err := template.ParseFiles("templates/message.tpl")
		if err != nil {
			log.Printf("Failed to parse template: %s\n", err)
			return
		}

		res := struct {
			Cars   []autoria.CarInfo
			Amount int64
		}{
			Cars:   newCars,
			Amount: search.Result.SearchResult.Count,
		}

		buff := bytes.Buffer{}
		if err := tpl.Execute(&buff, res); err != nil {
			log.Printf("Failed to execute template: %s\n", err)
			return
		}

		if err := msg.SendHTML(buff.String()); err != nil {
			log.Printf("send error: %s", err.Error())
		}
	})
}

func (h AutoRiaHandler) StopHandler(msg *bot.Event) {
	if err := msg.SetStatus(bot.ChatTyping); err != nil {
		log.Printf("action error: %s", err.Error())
	}

	if _, ok := h.Subscriptions[msg.Message.Chat.ID]; !ok {
		if err := msg.Send("Ви не підписані на оновлення 🤔"); err != nil {
			log.Printf("send error: %s", err.Error())
		}
		return
	}

	h.Subscriptions[msg.Message.Chat.ID].Stop()

	if err := msg.Send("Підписка призупинена ✅"); err != nil {
		log.Printf("send error: %s", err.Error())
	}
}

func (h AutoRiaHandler) AnalyzePhotos(photos []autoria.Photo) string {
	for _, photo := range photos {
		results, err := h.Toolkit.ALPR().Recognize(photo.URL())
		if err != nil {
			log.Println(err)
			continue
		}

		if len(results) == 0 {
			continue
		}

		return results[0].Plate
	}

	return ""
}

func (h AutoRiaHandler) getNumber(id string) (string, error) {
	car, err := h.API.CarInfo(id)
	if err != nil {
		log.Printf("error: %s", err.Error())
		return "", err
	}

	if car.PlateNumber != nil {
		res := strings.Join(strings.Split(*car.PlateNumber, " "), "")
		return res, nil
	}

	resp, err := h.API.CarPhotos(id)
	if err != nil {
		return "", err
	}

	plate := h.AnalyzePhotos(resp.Photos)
	if plate == "" {
		return "", fmt.Errorf("not found")
	}

	return plate, nil
}

// Analyze first 50 photos, then find best number, that matches the rules.
// Send message firstly.
func (h AutoRiaHandler) CarInfoHandler(msg *bot.Event) {
	if err := msg.SetStatus(bot.ChatTyping); err != nil {
		log.Printf("action error: %s", err.Error())
	}

	pattern := regexp.MustCompile(`(.*)([0-9]{8})(.*)`)
	id := strings.TrimSpace(pattern.ReplaceAllString(msg.Message.Text, "$2"))

	plate, err := h.getNumber(id)
	if err != nil {
		if err := msg.Send("Номер не знайдено"); err != nil {
			log.Printf("send error: %s\n", err.Error())
		}
	}

	operations, err := h.Toolkit.Operation().FindByNumber(plate)
	if err != nil {
		log.Println(err)
		return
	}

	tpl, err := template.ParseFiles("templates/operations.tpl")
	if err != nil {
		log.Println(err)
		return
	}

	buff := bytes.Buffer{}
	if err := tpl.Execute(&buff, struct {
		Operations []toolkit.Operation
		Number     string
	}{
		Operations: operations,
		Number:     plate,
	}); err != nil {
		log.Println(err)
		return
	}

	if err := msg.SendHTML(buff.String()); err != nil {
		log.Printf("send error: %s\n", err.Error())
	}

	registrations, err := h.Toolkit.Registration().FindByNumber(plate)
	if err != nil {
		log.Println(err)
		return
	}

	tpl2, err := template.ParseFiles("templates/registrations.tpl")
	if err != nil {
		log.Println(err)
		return
	}

	buff2 := bytes.Buffer{}
	if err := tpl2.Execute(&buff2, struct {
		Registrations []toolkit.Registration
		Number        string
	}{
		Registrations: registrations,
		Number:        plate,
	}); err != nil {
		log.Println(err)
		return
	}

	if err := msg.SendHTML(buff2.String()); err != nil {
		log.Printf("send error: %s\n", err.Error())
	}
}
