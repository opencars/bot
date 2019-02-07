package bot

import (
	"log"
	"regexp"
	"strings"

	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/shal/robot/pkg/env"
	"github.com/shal/robot/pkg/openalpr"
	"github.com/shal/robot/pkg/opencars"
)

type Handler interface {
	Handle(bot *tgbotapi.BotAPI, msg *tgbotapi.Message)
}

type HandlerFunc func(*tgbotapi.BotAPI, *tgbotapi.Message)

func (f HandlerFunc) Handle(bot *tgbotapi.BotAPI, msg *tgbotapi.Message) {
	f(bot, msg)
}

type MuxEntry struct {
	handler Handler
	match   func(x string) bool
}

type Bot struct {
	API        *tgbotapi.BotAPI
	Recognizer *openalpr.API
	Storage    *opencars.API
	FilePath   string
	Mux        []MuxEntry
}

func (bot *Bot) Handle(key string, handler Handler) {
	bot.Mux = append(bot.Mux, MuxEntry{
		handler: handler,
		match: func(text string) bool {
			return strings.HasPrefix(text, key)
		},
	})
}

func (bot *Bot) HandleRegexp(regexp *regexp.Regexp, handler Handler) {
	bot.Mux = append(bot.Mux, MuxEntry{
		handler: handler,
		match: func(x string) bool {
			return regexp.MatchString(x)
		},
	})
}

func (bot *Bot) HandleFuncRegexp(regexp *regexp.Regexp, handler func(*tgbotapi.BotAPI, *tgbotapi.Message)) {
	bot.Mux = append(bot.Mux, MuxEntry{
		handler: HandlerFunc(handler),
		match: func(x string) bool {
			return regexp.MatchString(x)
		},
	})
}

func (bot *Bot) HandleFunc(key string, handler func(*tgbotapi.BotAPI, *tgbotapi.Message)) {
	bot.Mux = append(bot.Mux, MuxEntry{
		handler: HandlerFunc(handler),
		match: func(text string) bool {
			return strings.HasPrefix(text, key)
		},
	})
}

func (bot *Bot) handle(update tgbotapi.Update) {
	for _, entry := range bot.Mux {
		if update.Message == nil {
			break
		}

		if entry.match(update.Message.Text) {
			entry.handler.Handle(bot.API, update.Message)
		}
	}
}

func (bot *Bot) Listen() error {
	// API updates configuration.
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60

	updates, err := bot.API.GetUpdatesChan(updateConfig)

	if err != nil {
		return err
	}

	// TODO: This should be concurrent, using webhook and server.
	// Triggers, when new update pushed to channel.
	for update := range updates {
		bot.handle(update)
	}

	return nil
}

func New(path, recognizerUrl, storageUrl string) *Bot {
	return &Bot{
		API:        NewAPI(),
		Recognizer: &openalpr.API{URL: recognizerUrl},
		Storage:    &opencars.API{URI: storageUrl},
		Mux:        make([]MuxEntry, 0),
		FilePath:   path,
	}
}

func NewAPI() *tgbotapi.BotAPI {
	bot, err := tgbotapi.NewBotAPI(env.MustGet("BOT_TOKEN"))

	bot.Debug = env.Get("LOG_LEVEL", "DEBUG") == "DEBUG"

	if err != nil {
		panic(err)
	}

	log.Printf("API authorized %s\n", bot.Self.UserName)

	return bot
}
