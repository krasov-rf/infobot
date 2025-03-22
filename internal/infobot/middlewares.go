package infobot

import (
	"database/sql"
	"errors"

	"github.com/krasov-rf/infobot/pkg/serializers"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *Bot) AddSiteMiddleware(user *serializers.User, f func(user *serializers.User)) *tgbotapi.MessageConfig {
	if user.GetActionSite() != nil {
		user.SetAction(serializers.ACTION_SITE_UPD)
		f(user)
		return nil
	}
	user.SetAction(serializers.ACTION_SITE_UPD_URL)
	user.SetActionSite(&serializers.SiteSerializer{})

	msg := tgbotapi.NewMessage(user.GetUserId(), "Введите URL сайта для добавления:")
	return &msg
}

// добавление пользователя в базу
func (b *Bot) RegisterUserMiddleware(next handlerFunc) handlerFunc {
	return func(ctx BotContext, update tgbotapi.Update) {
		chat := update.FromChat()

		_, err := b.DB.TelegramUserGet(b.ctx, chat.ID)
		if err != nil {
			if !errors.Is(err, sql.ErrNoRows) {
				b.errErrorChan <- err
				return
			}

			err := b.DB.TelegramUserRegister(b.ctx, &serializers.UserSerializer{
				UserId:    chat.ID,
				UserName:  chat.UserName,
				FirstName: chat.FirstName,
				LastName:  chat.LastName,
			})
			if err != nil {
				b.errErrorChan <- err
				return
			}
		}

		ctx.user = users.Get(chat.ID)

		next(ctx, update)
	}
}
