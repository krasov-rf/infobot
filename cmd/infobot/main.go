package main

import (
	"log"

	"github.com/krasov-rf/infobot/internal/infobot"
	infobot_settings "github.com/krasov-rf/infobot/pkg/settings/infobot"
)

func main() {
	config, err := infobot_settings.InitEnv()
	if err != nil {
		log.Fatal(err)
	}

	// Создаем нового бота
	bot, err := infobot.New(config)
	if err != nil {
		log.Fatalf("Failed to create bot: %v", err)
	}

	// Запускаем бота
	bot.Start()
}
