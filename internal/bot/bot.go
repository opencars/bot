package bot

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

	"github.com/opencars/bot/pkg/env"
)

// Constant values for ChatActions.
const (
	ChatTyping = "typing"
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
	api          *tgbotapi.BotAPI
	mux          []MuxEntry
	photoHandler Handler
}

// Handle registers handler by key.
func (bot *Bot) Handle(key string, handler Handler) {
	bot.mux = append(bot.mux, MuxEntry{
		handler: handler,
		match: func(text string) bool {
			return strings.HasPrefix(text, key)
		},
	})
}

// HandleRegexp registers handler by regular expression.
func (bot *Bot) HandleRegexp(regexp *regexp.Regexp, handler Handler) {
	bot.mux = append(bot.mux, MuxEntry{
		handler: handler,
		match:   regexp.MatchString,
	})
}

// HandleFuncRegexp registers handler function by regular expression.
func (bot *Bot) HandleFuncRegexp(regexp *regexp.Regexp, handler func(*Event)) {
	bot.HandleRegexp(regexp, HandlerFunc(handler))
}

// HandleFunc registers handler function by key.
func (bot *Bot) HandleFunc(key string, handler func(*Event)) {
	bot.Handle(key, HandlerFunc(handler))
}

func (bot *Bot) HandlePhoto(handler func(*Event)) {
	bot.photoHandler = HandlerFunc(handler)
}

func (bot *Bot) handleMsg(request *tgbotapi.Message) {
	log.Printf("User { ID: %d, name: %s }\n", request.From.ID, request.From)

	if request.Photo != nil && bot.photoHandler != nil {
		bot.photoHandler.Handle(&Event{bot.api, request})
	}

	for _, entry := range bot.mux {
		if entry.match(request.Text) {
			log.Printf("Text: %s\n", request.Text)
			entry.handler.Handle(&Event{bot.api, request})
			return
		}
	}
}

func (bot *Bot) handle(update *tgbotapi.Update) {
	if update.Message != nil {
		bot.handleMsg(update.Message)
	}
}

// Listen for telegram updates.
func (bot *Bot) Listen(host, port string) error {
	URL := fmt.Sprintf("%s/tg/%s", host, bot.api.Token)
	_, err := bot.api.SetWebhook(tgbotapi.NewWebhook(URL))
	if err != nil {
		log.Fatal(err)
	}

	path := fmt.Sprintf("/tg/%s", bot.api.Token)

	http.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		bytes, _ := ioutil.ReadAll(r.Body)
		r.Body.Close()

		update := new(tgbotapi.Update)
		if err := json.Unmarshal(bytes, update); err != nil {
			log.Printf("update error: %s", err.Error())
		}

		// TODO: Log message content and user ID.
		bot.handle(update)
	})

	return http.ListenAndServe(":"+port, http.DefaultServeMux)
}

// New creates new instance of the Bot.
// Idiomatically, there is only one "Bot" instance per application.
func New() *Bot {
	return &Bot{
		api: newAPI(),
		mux: make([]MuxEntry, 0),
	}
}

func newAPI() *tgbotapi.BotAPI {
	telegramToken := env.MustFetch("TELEGRAM_TOKEN")
	bot, err := tgbotapi.NewBotAPI(telegramToken)
	if err != nil {
		panic(err)
	}

	bot.Debug = false
	log.Printf("API authorized %s\n", bot.Self.UserName)

	return bot
}
