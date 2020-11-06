package handlers

import (
	"github.com/opencars/bot/internal/bot"
	"github.com/opencars/bot/pkg/logger"
)

func somethingWentWrong(msg *bot.Event) {
	if err := msg.Send("Вибач. Щось пішло не так 😢"); err != nil {
		logger.Errorf("send error: %s", err.Error())
	}
}

func (h OpenCarsHandler) PhotoHandler(msg *bot.Event) {
	photos := *msg.Message.Photo

	if err := msg.SetStatus(bot.ChatTyping); err != nil {
		logger.Errorf("action error: %s", err.Error())
	}

	if len(photos) < 1 {
		somethingWentWrong(msg)
		return
	}

	url, err := msg.API.GetFileDirectURL(photos[len(photos)-1].FileID)
	if err != nil {

		somethingWentWrong(msg)
		logger.Errorf(err.Error())
		return
	}

	logger.Debugf("Photo: %s", url)

	results, err := h.client.ALPR().Recognize(url)
	if err != nil {
		somethingWentWrong(msg)
		logger.Errorf("recognize: %s", err)
		return
	}

	// If nothing was found, send user notification.
	if len(results) == 0 {
		if err := msg.Send("Номер не знайдено 🤔"); err != nil {
			logger.Errorf("send error: %s", err.Error())
		}
		return
	}

	plates := make([]string, len(results))
	for i := range results {
		plates[i] = results[i].Plate
	}

	// Send number to user.
	if err := msg.Send(plates[0]); err != nil {
		logger.Errorf("send error: %s", err.Error())
	}

	text, err := h.getInfoByNumber(plates[0])
	if err != nil {
		logger.Errorf("failed to get info by number: %s", err)
	}

	if err := msg.SendHTML(text); err != nil {
		logger.Errorf("send error: %s", err.Error())
	}
}
