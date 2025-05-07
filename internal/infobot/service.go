package infobot

import (
	"context"
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

	config *settings.Config
	cron   *cron.Cron
	ctx    context.Context

	DB infobotdb.IInfoBotDB

	updateChan tgbotapi.UpdatesChannel
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

	return &Bot{
		BotAPI: bot,
		DB:     db,
		config: c,
		cron:   cron.New(),
		Router: NewRouter(),
	}, nil
}

func (b *Bot) Start() {
	var cancel context.CancelFunc

	b.ctx, cancel = context.WithCancel(context.Background())
	defer cancel()

	defer b.DB.Close(b.ctx)

	logFile, err := os.OpenFile("bot.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Ошибка при открытии файла логов: %v", err)
	}
	defer logFile.Close()
	log.SetOutput(logFile)

	b.InitializeRoutes()

	err = b.InitializeCron()
	if err != nil {
		log.Fatal(err)
		return
	}

	b.cron.Start()
	defer b.cron.Stop()

	updater := tgbotapi.NewUpdate(0)
	updater.Timeout = 60
	b.updateChan = b.GetUpdatesChan(updater)
	go b.updateListener()

	b.errChan = make(chan error, 10)
	defer close(b.errChan)
	go b.errorListener()

	_, err = b.Send(tgbotapi.NewMessage(b.config.TG_SUPER_ADMIN, "Bot started"))
	if err != nil {
		log.Fatal(err)
		return
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs

	log.Print("Бот остановлен")
}
