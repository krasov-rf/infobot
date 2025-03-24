package infobot

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/krasov-rf/infobot/pkg/serializers"
	settings "github.com/krasov-rf/infobot/pkg/settings/infobot"
	infobotdb "github.com/krasov-rf/infobot/pkg/storage/infobot"
	infobotdb_pg "github.com/krasov-rf/infobot/pkg/storage/infobot/postgres"
	"github.com/robfig/cron/v3"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Bot struct {
	*Router
	*tgbotapi.BotAPI

	ctx    context.Context
	config *settings.Config
	cron   *cron.Cron

	DB infobotdb.IInfoBotDB

	updateChan chan tgbotapi.Update
	errChan    chan error
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
		cron:   cron.New(),
		Router: NewRouter(),
	}, nil
}

func (b *Bot) Start() {
	b.InitializeRoutes()

	err := b.InitializeCron()
	if err != nil {
		log.Fatal(err)
	}

	b.cron.Start()
	defer b.cron.Stop()

	b.errChan = make(chan error, 10)
	defer close(b.errChan)
	go b.errorListener()

	b.updateChan = make(chan tgbotapi.Update, 10)
	defer close(b.updateChan)
	go b.updateListener()

	_, err = b.Send(tgbotapi.NewMessage(b.config.TG_SUPER_ADMIN, "Bot started"))
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		u := tgbotapi.NewUpdate(0)
		u.Timeout = 60
		for update := range b.GetUpdatesChan(u) {
			b.handleUpdate(b.ctx, update)
		}
	}()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	<-sigs

	fmt.Println("Остановка бота...")
}
