package infobot

import (
	"context"
	"log"

	"github.com/krasov-rf/infobot/pkg/serializers"
	settings "github.com/krasov-rf/infobot/pkg/settings/infobot"
	infobotdb "github.com/krasov-rf/infobot/pkg/storage/infobot"
	infobotdb_pg "github.com/krasov-rf/infobot/pkg/storage/infobot/postgres"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Bot struct {
	*Router
	*tgbotapi.BotAPI

	ctx    context.Context
	config *settings.Config

	DB infobotdb.IInfoBotDB

	updateChan   chan tgbotapi.Update
	errErrorChan chan error
}

type BotContext struct {
	context.Context
	user *serializers.User
}

type BotContextKeys int

const (
	CTX_KEY_DATA BotContextKeys = iota
	CTX_KEY_SITE
)

var users = serializers.NewUsers()

func New(c *settings.Config) (*Bot, error) {
	bot, err := tgbotapi.NewBotAPI(c.TG_TOKEN)
	if err != nil {
		return nil, err
	}

	db, err := infobotdb_pg.New(c)
	if err != nil {
		return nil, err
	}

	ctx := context.Background()

	return &Bot{
		BotAPI: bot,
		DB:     db,
		ctx:    ctx,
		config: c,
		Router: NewRouter(),
	}, nil
}

// старт
func (b *Bot) Start() {
	b.InitializeRoutes()

	msg := tgbotapi.NewMessage(b.config.TG_SUPER_ADMIN, "Bot started")
	_, err := b.Send(msg)
	if err != nil {
		log.Fatal(err)
	}

	b.errErrorChan = make(chan error, 10)
	defer close(b.errErrorChan)
	b.updateChan = make(chan tgbotapi.Update, 10)
	defer close(b.updateChan)

	go b.errorListener()
	go b.updateListener()

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.GetUpdatesChan(u)

	for update := range updates {
		b.handleUpdate(update)
	}
}
