package bot

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

	"github.com/opencars/bot/pkg/logger"
	"github.com/opencars/bot/pkg/model"
	"github.com/opencars/bot/pkg/store"
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

	l := logger.WithFields(logger.Fields{
		"id":           update.UpdateID,
		"text":         update.Message.Text,
		"message_date": update.Message.Date,
	})

	event := Event{
		API:     bot.api,
		Message: update.Message,
	}

	user := model.User{
		ID:        update.Message.From.ID,
		FirstName: update.Message.From.FirstName,
	}

	if update.Message.From.LastName != "" {
		user.LastName = &update.Message.From.LastName
	}

	if update.Message.From.UserName != "" {
		user.UserName = &update.Message.From.UserName
	}

	if update.Message.From.LanguageCode == "" {
		user.LanguageCode = &update.Message.From.LanguageCode
	}

	if err := bot.store.User().Create(&user); err != nil {
		l.Errorf("failed to create user: %s", err)
	}

	if update.Message.Photo != nil && bot.photoHandler != nil {
		bot.photoHandler.Handle(&event)
	}

	for _, entry := range bot.mux {
		if !entry.match(update.Message.Text) {
			continue
		}

		dto := model.Update{
			ID:     update.Message.MessageID,
			UserID: update.Message.From.ID,
			Text:   update.Message.Text,
			Time:   time.Unix(int64(update.Message.Date), 0),
		}

		if err := bot.store.Update().Create(&dto); err != nil {
			l.Errorf("failed to save update: %s", err)
		}

		l.WithFields(logger.Fields{
			"user_id":  update.Message.From.ID,
			"username": update.Message.From.UserName,
		}).Debugf("handle event")

		entry.handler.Handle(&event)
		return
	}
}

// Listen for telegram updates.
func (bot *Bot) Listen(host string, addr string) error {
	URL := fmt.Sprintf("%s/api/v1/bot/tg/%s", host, bot.api.Token)
	_, err := bot.api.SetWebhook(tgbotapi.NewWebhook(URL))
	if err != nil {
		logger.Errorf("sql: %s", err)
	}

	path := fmt.Sprintf("/api/v1/bot/tg/%s", bot.api.Token)

	http.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		bytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			logger.Errorf("failed to parse body: %s", err)
			return
		}

		_ = r.Body.Close()

		update := new(tgbotapi.Update)
		if err := json.Unmarshal(bytes, update); err != nil {
			logger.Errorf("update error: %s", err)
		}

		bot.handle(update)
	})

	return http.ListenAndServe(addr, http.DefaultServeMux)
}

// New creates new instance of the Bot.
// Idiomatically, there is only one "Bot" instance per application.
func New(token string, store store.Store) (*Bot, error) {
	api, err := newAPI(token)
	if err != nil {
		return nil, err
	}

	return &Bot{
		api:   api,
		mux:   make([]MuxEntry, 0),
		store: store,
	}, nil
}

func newAPI(token string) (*tgbotapi.BotAPI, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}

	bot.Debug = false

	logger.WithFields(logger.Fields{
		"username": bot.Self.UserName,
	}).Infof("Bot authorized")

	return bot, nil
}
