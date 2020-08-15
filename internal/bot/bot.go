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

	"github.com/opencars/bot/pkg/model"
	"github.com/opencars/bot/pkg/store"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

	"github.com/opencars/bot/pkg/env"
)

// Constant values for ChatActions.
const (
	ChatTyping = "typing"
)

var (
	WebPagePreview = false
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
	store        store.Store
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
func (bot *Bot) HandleFuncRegexp(expr string, handler func(*Event)) {
	pattern, err := regexp.Compile(expr)
	if err != nil {
		panic(err)
	}

	bot.HandleRegexp(pattern, HandlerFunc(handler))
}

// HandleFunc registers handler function by key.
func (bot *Bot) HandleFunc(key string, handler func(*Event)) {
	bot.Handle(key, HandlerFunc(handler))
}

func (bot *Bot) HandlePhoto(handler func(*Event)) {
	bot.photoHandler = HandlerFunc(handler)
}

func (bot *Bot) handle(update *tgbotapi.Update) {
	if update.Message == nil {
		return
	}

	event := Event{
		API:     bot.api,
		Message: update.Message,
	}

	lastName := &update.Message.From.LastName
	if *lastName == "" {
		lastName = nil
	}

	userName := &update.Message.From.UserName
	if *userName == "" {
		userName = nil
	}

	lang := &update.Message.From.LanguageCode
	if *lang == "" {
		lang = nil
	}

	user := model.User{
		ID:           update.Message.From.ID,
		FirstName:    update.Message.From.FirstName,
		LastName:     lastName,
		UserName:     userName,
		LanguageCode: lang,
	}

	if err := bot.store.User().Create(&user); err != nil {
		log.Println(err)
	}

	if update.Message.Photo != nil && bot.photoHandler != nil {
		bot.photoHandler.Handle(&event)
	}

	for _, entry := range bot.mux {
		if !entry.match(update.Message.Text) {
			continue
		}

		tmp := model.Update{
			ID:     update.Message.MessageID,
			UserID: update.Message.From.ID,
			Text:   update.Message.Text,
			Time:   time.Unix(int64(update.Message.Date), 0),
		}

		if err := bot.store.Update().Create(&tmp); err != nil {
			log.Println(err)
		}

		log.Printf("Matched Text: %s\n", update.Message.Text)
		entry.handler.Handle(&event)
		return
	}

	log.Printf("Text: %s\n", update.Message.Text)
}

// Listen for telegram updates.
func (bot *Bot) Listen(host, port string) error {
	URL := fmt.Sprintf("%s/api/v1/bot/tg/%s", host, bot.api.Token)
	_, err := bot.api.SetWebhook(tgbotapi.NewWebhook(URL))
	if err != nil {
		log.Fatal(err)
	}

	path := fmt.Sprintf("/api/v1/bot/tg/%s", bot.api.Token)

	http.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		bytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Printf("failed to parse body: %v", err)
			return
		}

		_ = r.Body.Close()

		update := new(tgbotapi.Update)
		if err := json.Unmarshal(bytes, update); err != nil {
			log.Printf("update error: %s", err.Error())
		}

		bot.handle(update)
	})

	return http.ListenAndServe(":"+port, http.DefaultServeMux)
}

// New creates new instance of the Bot.
// Idiomatically, there is only one "Bot" instance per application.
func New(store store.Store) *Bot {
	return &Bot{
		api:   newAPI(),
		mux:   make([]MuxEntry, 0),
		store: store,
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
