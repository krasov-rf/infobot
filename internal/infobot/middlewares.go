package infobot

import (
	"database/sql"
	"errors"

	"github.com/krasov-rf/infobot/pkg/serializers"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// добавление пользователя в базу
func (b *Bot) RegisterUserMiddleware(next handlerFunc) handlerFunc {
	return func(ctx *BotContext, update tgbotapi.Update) {
		chat := update.FromChat()

		_, err := b.DB.TelegramUserGet(ctx, chat.ID)
		if err != nil {
			if !errors.Is(err, sql.ErrNoRows) {
				b.errChan <- err
				return
			}

			err := b.DB.TelegramUserRegister(ctx, &serializers.UserSerializer{
				UserId:    chat.ID,
				UserName:  chat.UserName,
				FirstName: chat.FirstName,
				LastName:  chat.LastName,
			})
			if err != nil {
				b.errChan <- err
				return
			}
		}

		ctx.user = users.Get(chat.ID)

		next(ctx, update)
	}
}
