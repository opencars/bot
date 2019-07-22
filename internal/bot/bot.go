package bot

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/go-telegram-bot-api/telegram-bot-api"

	"github.com/opencars/bot/pkg/env"
	"github.com/opencars/bot/pkg/openalpr"
	"github.com/opencars/toolkit/sdk"
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

type HandlerFunc func(msg *Event)

func (f HandlerFunc) Handle(msg *Event) {
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
	Storage    *sdk.Client
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
func (bot *Bot) HandleFuncRegexp(regexp *regexp.Regexp, handler func(*Event)) {
	bot.Mux = append(bot.Mux, MuxEntry{
		handler: HandlerFunc(handler),
		match: func(x string) bool {
			return regexp.MatchString(x)
		},
	})
}

// HandleFunc registers handler function by key.
func (bot *Bot) HandleFunc(key string, handler func(*Event)) {
	bot.Mux = append(bot.Mux, MuxEntry{
		handler: HandlerFunc(handler),
		match: func(text string) bool {
			log.Println(text, strings.HasPrefix(text, key))
			return strings.HasPrefix(text, key)
		},
	})
}

func (bot *Bot) handleMsg(request *tgbotapi.Message) {
	msg := &Event{bot.API, request}

	if request.Photo != nil {
		for _, entry := range bot.Mux {
			if entry.match(PhotoEvent) {
				entry.handler.Handle(msg)
			}
		}
	}

	if request.Sticker != nil {
		for _, entry := range bot.Mux {
			if entry.match(StickerEvent) {
				entry.handler.Handle(msg)
			}
		}
	}

	for _, entry := range bot.Mux {
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

	path := fmt.Sprintf("/tg/%s", bot.API.Token)

	http.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		bytes, _ := ioutil.ReadAll(r.Body)
		r.Body.Close()

		update := tgbotapi.Update{}

		if err := json.Unmarshal(bytes, &update); err != nil {
			log.Printf("update error: %s", err.Error())
		}

		start := time.Now()
		fmt.Printf("Started time: %s\n", start.Format("2006-01-02 15:04:05"))
		bot.handle(update)
		fmt.Printf("Finished time: %s\n", time.Now().Format("2006-01-02 15:04:05"))
		fmt.Printf("Execution time: %f\n", time.Since(start).Seconds())
	})

	return http.ListenAndServe(":"+port, http.DefaultServeMux)
}

// New creates new instance of the Bot.
// Idiomatically, there is only one "Bot" instance per application.
func New(path, recognizerUrl, storageUrl string) *Bot {
	return &Bot{
		API:        newAPI(),
		Recognizer: &openalpr.API{URI: recognizerUrl},
		Storage:    sdk.New(storageUrl),
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
