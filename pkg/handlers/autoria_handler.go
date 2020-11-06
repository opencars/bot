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

	"github.com/opencars/bot/pkg/logger"

	"github.com/opencars/toolkit"

	"github.com/opencars/bot/internal/bot"
	"github.com/opencars/bot/internal/subscription"
	"github.com/opencars/bot/pkg/autoria"
)

type AutoRiaHandler struct {
	API           *autoria.API
	Period        time.Duration
	Toolkit       *toolkit.Client
	Subscriptions map[int64]*subscription.Subscription
}

func (h AutoRiaHandler) FollowHandler(msg *bot.Event) {
	if err := msg.SetStatus(bot.ChatTyping); err != nil {
		logger.Errorf("action error: %s", err)
	}

	if !strings.HasPrefix(msg.Message.Text, "https://auto.ria.com/search") {
		if err := msg.Send("Помилковий запит."); err != nil {
			logger.Errorf("send error: %s", err)
		}
		return
	}

	values, err := url.ParseQuery(msg.Message.Text)
	if err != nil {
		if err := msg.Send(err.Error()); err != nil {
			logger.Errorf("send error: %s", err)
		}
		return
	}

	// Convert params to old type, because frontend and api have different types.
	values, err = h.API.ConvertNewToOld(values)
	if err != nil {
		if err := msg.Send(err.Error()); err != nil {
			logger.Errorf("send error: %s", err)
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
			logger.Errorf("Failed to search cars: %s", err)
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
				logger.Errorf("Failed to get car info: %s", err)
				return
			}

			newCars[i] = *car
		}

		tpl, err := template.ParseFiles("templates/message.tpl")
		if err != nil {
			logger.Errorf("Failed to parse template: %s", err)
			return
		}

		res := struct {
			Cars   []autoria.CarInfo
			Amount int64
		}{
			Cars:   newCars,
			Amount: search.Result.SearchResult.Count,
		}

		var buff bytes.Buffer
		if err := tpl.Execute(&buff, res); err != nil {
			logger.Errorf("Failed to execute template: %s", err)
			return
		}

		if err := msg.SendHTML(buff.String()); err != nil {
			logger.Errorf("send error: %s", err)
		}
	})
}

func (h AutoRiaHandler) StopHandler(msg *bot.Event) {
	if err := msg.SetStatus(bot.ChatTyping); err != nil {
		logger.Errorf("action error: %s", err)
	}

	if _, ok := h.Subscriptions[msg.Message.Chat.ID]; !ok {
		if err := msg.Send("Ви не підписані на оновлення 🤔"); err != nil {
			logger.Errorf("send error: %s", err)
		}
		return
	}

	h.Subscriptions[msg.Message.Chat.ID].Stop()

	if err := msg.Send("Підписка призупинена ✅"); err != nil {
		logger.Errorf("send error: %s", err)
	}
}

func (h AutoRiaHandler) AnalyzePhotos(photos []autoria.Photo) string {
	for _, photo := range photos {
		results, err := h.Toolkit.ALPR().Recognize(photo.URL())
		if err != nil {
			logger.Errorf("recognize: %s", err)
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
		logger.Errorf("getNumber: %s", err)
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
		logger.Errorf("action error: %s", err)
	}

	pattern := regexp.MustCompile(`(.*)([0-9]{8})(.*)`)
	id := strings.TrimSpace(pattern.ReplaceAllString(msg.Message.Text, "$2"))

	plate, err := h.getNumber(id)
	if err != nil {
		if err := msg.Send("Номер не знайдено"); err != nil {
			logger.Errorf("send error: %s", err)
		}
	}

	operations, err := h.Toolkit.Operation().FindByNumber(plate)
	if err != nil {
		logger.Errorf("find by number: %s", err)
		return
	}

	tpl, err := template.ParseFiles("templates/operations.tpl")
	if err != nil {
		logger.Errorf("failed to parse files: %s", err)
		return
	}

	var buff bytes.Buffer
	if err := tpl.Execute(&buff, struct {
		Operations []toolkit.Operation
		Number     string
		Code       *string
	}{
		Operations: operations,
		Number:     plate,
	}); err != nil {
		logger.Errorf("execute: %s", err)
		return
	}

	if err := msg.SendHTML(buff.String()); err != nil {
		logger.Errorf("send error: %s", err)
	}

	registrations, err := h.Toolkit.Registration().FindByNumber(plate)
	if err != nil {
		logger.Errorf("find by number: %s", err)
		return
	}

	tpl2, err := template.ParseFiles("templates/registrations.tpl")
	if err != nil {
		logger.Errorf("parse files: %s", err)
		return
	}

	var buff2 bytes.Buffer
	if err := tpl2.Execute(&buff2, struct {
		Registrations []toolkit.Registration
		Number        string
	}{
		Registrations: registrations,
		Number:        plate,
	}); err != nil {
		log.Println(err)
		logger.Debugf("execute: %s", err)

		return
	}

	if err := msg.SendHTML(buff2.String()); err != nil {
		logger.Errorf("send error: %s", err)
	}
}
