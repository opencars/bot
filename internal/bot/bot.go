package bot

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/shal/opencars-bot/pkg/env"
	"github.com/shal/opencars-bot/pkg/openalpr"
	"github.com/shal/opencars-bot/pkg/opencars"
)

// Constant values for ChatActions.
const (
	ChatTyping = "typing"
)

const (
	PhotoEvent   = "/photo"
	StickerEvent = "/sticker"
)

var (
	WebPagePreview = true
)

type HandlerFunc func(msg *Message)

func (f HandlerFunc) Handle(msg *Message) {
	f(msg)
}

type MuxEntry struct {
	handler Handler
	match   func(x string) bool
}

// Bot is structure representation of Bot instance.
type Bot struct {
	API        *tgbotapi.BotAPI
	Recognizer *openalpr.API
	Storage    *opencars.API
	FilePath   string
	Mux        []MuxEntry
}

// Handle registers handler by key.
func (bot *Bot) Handle(key string, handler Handler) {
	bot.Mux = append(bot.Mux, MuxEntry{
		handler: handler,
		match: func(text string) bool {
			return strings.HasPrefix(text, key)
		},
	})
}

// HandleRegexp registers handler by regular expression.
func (bot *Bot) HandleRegexp(regexp *regexp.Regexp, handler Handler) {
	bot.Mux = append(bot.Mux, MuxEntry{
		handler: handler,
		match: func(x string) bool {
			return regexp.MatchString(x)
		},
	})
}

// HandleFuncRegexp registers handler function by regular expression.
func (bot *Bot) HandleFuncRegexp(regexp *regexp.Regexp, handler func(*Message)) {
	bot.Mux = append(bot.Mux, MuxEntry{
		handler: HandlerFunc(handler),
		match: func(x string) bool {
			return regexp.MatchString(x)
		},
	})
}

// HandleFunc registers handler function by key.
func (bot *Bot) HandleFunc(key string, handler func(*Message)) {
	bot.Mux = append(bot.Mux, MuxEntry{
		handler: HandlerFunc(handler),
		match: func(text string) bool {
			return strings.HasPrefix(text, key)
		},
	})
}

func (bot *Bot) handleMsg(request *tgbotapi.Message) {
	msg := &Message{bot.API, request}

	for _, entry := range bot.Mux {
		if entry.match(PhotoEvent) {
			if request.Photo != nil {
				entry.handler.Handle(msg)
			}
			return
		}

		if entry.match(StickerEvent) {
			if request.Sticker != nil {
				entry.handler.Handle(msg)
			}
			return
		}

		if entry.match(request.Text) {
			entry.handler.Handle(msg)
			return
		}
	}
}

func (bot *Bot) handle(update tgbotapi.Update) {
	if update.Message != nil {
		bot.handleMsg(update.Message)
	}
}

// Listen for telegram updates.
func (bot *Bot) Listen(host, port string) error {
	URL := fmt.Sprintf("%s/tg/%s", host, bot.API.Token)
	_, err := bot.API.SetWebhook(tgbotapi.NewWebhook(URL))
	if err != nil {
		log.Fatal(err)
	}

	// Should be thread-safe out of the box.	
	path := fmt.Sprintf("/tg/%s", bot.API.Token)

	http.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		bytes, _ := ioutil.ReadAll(r.Body)
		r.Body.Close()

		update := tgbotapi.Update{}

		if err := json.Unmarshal(bytes, &update); err != nil {
			log.Printf("update error: %s", err.Error())
		}

		bot.handle(update)
	})

	return http.ListenAndServe(":"+port, http.DefaultServeMux)
}

// New creates new instance of the Bot.
// Idiomatically, there is only one "Bot" instance per application.
func New(path, recognizerUrl, storageUrl string) *Bot {
	return &Bot{
		API:        newAPI(),
		Recognizer: &openalpr.API{URI: recognizerUrl},
		Storage:    &opencars.API{URI: storageUrl},
		Mux:        make([]MuxEntry, 0),
		FilePath:   path,
	}
}

func newAPI() *tgbotapi.BotAPI {
	telegramToken := env.MustFetch("TELEGRAM_TOKEN")
	bot, err := tgbotapi.NewBotAPI(telegramToken)
	if err != nil {
		panic(err)
	}

	//bot.Debug = env.Fetch("LOG_LEVEL", "DEBUG") == "DEBUG"
	bot.Debug = false
	log.Printf("API authorized %s\n", bot.Self.UserName)

	return bot
}
