package infobot

import (
	"fmt"
	"log"
	"sync"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// общий слушатель сообщений
func (b *Bot) updateListener() {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-b.ctx.Done():
				return

			case update, ok := <-b.updateChan:
				if !ok {
					return
				}
				b.handleUpdate(update)
			}
		}
	}()

	go func() {
		wg.Wait()
	}()
}

// общий слушатель ошибок для потоков
func (b *Bot) errorListener() {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-b.ctx.Done():
				return

			case err, ok := <-b.errErrorChan:
				if !ok {
					return
				}
				if err != nil {
					msg := tgbotapi.NewMessage(b.config.TG_SUPER_ADMIN, fmt.Sprintf("Ошибка в ходе выполенния бота: %v", err))
					b.Send(msg)
					log.Print(err)
				}
			}
		}
	}()

	go func() {
		wg.Wait()
	}()
}
