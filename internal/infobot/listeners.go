package infobot

import (
	"errors"
	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// общий слушатель сообщений
func (b *Bot) updateListener() {
	go func() {
		for {
			select {
			case <-b.ctx.Done():
				return

			case update, ok := <-b.updateChan:
				if !ok {
					return
				}
				b.handleUpdate(b.ctx, update)
			}
		}
	}()
}

// общий слушатель ошибок для потоков
func (b *Bot) errorListener() {
	go func() {
		for {
			select {
			case <-b.ctx.Done():
				return

			case err, ok := <-b.errChan:
				if !ok {
					return
				}
				if err != nil {
					_, err_tg := b.Send(tgbotapi.NewMessage(b.config.TG_SUPER_ADMIN, fmt.Sprintf("Ошибка в ходе выполенния бота: %v", err)))
					if err_tg != nil {
						err = errors.Join(err, err_tg)
					}
					log.Print(err)
				}
			}
		}
	}()
}
