package test

import (
	"log"
	"runtime"
	"testing"

	"github.com/krasov-rf/infobot/internal/infobot"
	infobot_settings "github.com/krasov-rf/infobot/pkg/settings/infobot"
)

var (
	err    error
	config *infobot_settings.Config
)

func TestInfoBot(t *testing.T) {
	bot, err := infobot.New(config)
	if err != nil {
		t.Fatal()
	}
	bot.Start()
}

func TestMain(m *testing.M) {
	config, err = infobot_settings.InitEnv()
	if err != nil {
		log.Fatal(err)
	}

	runtime.GOMAXPROCS(runtime.NumCPU() - 1)
	m.Run()
}
