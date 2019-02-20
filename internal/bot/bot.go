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
func (bot *Bot) HandleFuncRegexp(regexp *regexp.Regexp, handler func(*tgbotapi.BotAPI, *tgbotapi.Message)) {
	bot.Mux = append(bot.Mux, MuxEntry{
		handler: HandlerFunc(handler),
		match: func(x string) bool {
			return regexp.MatchString(x)
		},
	})
}

// HandleFunc registers handler function by key.
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

		fmt.Printf("Incoming request %v\n", r)
		// Handle "Update".
		bot.handle(update)
	})

	return http.ListenAndServe(":"+port, http.DefaultServeMux)
}

// New creates new instanse of the Bot.
// Idiomatically, there is only one "Bot" instance per application.
func New(path, recognizerUrl, storageUrl string) *Bot {
	return &Bot{
		API:        NewAPI(),
		Recognizer: &openalpr.API{URL: recognizerUrl},
		Storage:    &opencars.API{URI: storageUrl},
		Mux:        make([]MuxEntry, 0),
		FilePath:   path,
	}
}

// NewAPI creates new instance without Debug logs by default.
// Export DEBUG=true to enable debug logs.
func NewAPI() *tgbotapi.BotAPI {
	telegramToken := env.MustGet("TELEGRAM_TOKEN")
	bot, err := tgbotapi.NewBotAPI(telegramToken)

	//bot.
	bot.Debug = env.Get("LOG_LEVEL", "DEBUG") == "DEBUG"

	if err != nil {
		panic(err)
	}

	log.Printf("API authorized %s\n", bot.Self.UserName)

	return bot
}

func send(bot *tgbotapi.BotAPI, message tgbotapi.MessageConfig) error {
	if _, err := bot.Send(message); err != nil {
		return err
	}

	return nil
}

// Send sends message to the chat with regular text.
func Send(bot *tgbotapi.BotAPI, chat *tgbotapi.Chat, text string) error {
	msg := tgbotapi.NewMessage(chat.ID, text)
	return send(bot, msg)
}

// SendHTML sends message to the chat with text formatted as HTML.
func SendHTML(bot *tgbotapi.BotAPI, chat *tgbotapi.Chat, text string) error {
	msg := tgbotapi.NewMessage(chat.ID, text)
	msg.ParseMode = tgbotapi.ModeHTML

	return send(bot, msg)
}

// SendMsgHTML sends message to the chat with text formatted as HTML.
func SendMsgHTML(msg tgbotapi.MessageConfig, bot *tgbotapi.BotAPI) error {
	msg.ParseMode = tgbotapi.ModeHTML

	return send(bot, msg)
}

// SendMarkdown sends message to the chat with text formatted as Markdown.
func SendMarkdown(bot *tgbotapi.BotAPI, chat *tgbotapi.Chat, text string) error {
	msg := tgbotapi.NewMessage(chat.ID, text)
	msg.ParseMode = tgbotapi.ModeMarkdown

	return send(bot, msg)
}
