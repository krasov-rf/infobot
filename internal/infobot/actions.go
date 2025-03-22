package infobot

import (
	"errors"
	"fmt"
	"strconv"

	er "github.com/krasov-rf/infobot/pkg/errors"
	"github.com/krasov-rf/infobot/pkg/serializers"
	infobotdb "github.com/krasov-rf/infobot/pkg/storage/infobot"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// сместить позицию листинга
func (b *Bot) updateOffset(ctx BotContext, update tgbotapi.Update) {
	switch ctx.user.GetAction() {
	case serializers.ACTION_SITE_LIST:
		b.HB_Sites(ctx, update)
	case serializers.ACTION_FEEDBACK_LIST:
		b.HB_Feedbacks(ctx, update)
	}
}

// сместить позицию листинга далее
func (b *Bot) HB_UpdateOffsetNext(ctx BotContext, update tgbotapi.Update) {
	ctx.user.SetOffset(ctx.user.GetOffset() + infobotdb.QUERY_LIMIT)
	b.updateOffset(ctx, update)
}

// сместить позицию листинга назад
func (b *Bot) HB_UpdateOffsetPrevious(ctx BotContext, update tgbotapi.Update) {
	ctx.user.SetOffset(ctx.user.GetOffset() - infobotdb.QUERY_LIMIT)
	b.updateOffset(ctx, update)
}

// вывести главную страницу
func (b *Bot) HB_HomePage(ctx BotContext, update tgbotapi.Update) {
	ctx.user.SetAction(serializers.ACTION_NONE)
	_, err := b.Send(tgbotapi.NewEditMessageTextAndMarkup(
		update.CallbackQuery.Message.Chat.ID,
		update.CallbackQuery.Message.MessageID,
		"Главное меню:",
		*b.KeyboardHomePage(ctx.user),
	))
	if err != nil {
		b.errErrorChan <- err
	}
}

// вывести главную страницу
func (b *Bot) MSG_HomePage(ctx BotContext, update tgbotapi.Update) {
	msg := tgbotapi.NewMessage(ctx.user.GetUserId(), "Добро пожаловать!")
	ctx.user.SetAction(serializers.ACTION_NONE)
	msg.ReplyMarkup = b.KeyboardHomePage(ctx.user)
	_, err := b.Send(msg)
	if err != nil {
		b.errErrorChan <- err
	}
}

// вывести хелпу
func (b *Bot) HB_Help(ctx BotContext, update tgbotapi.Update) {
	text := `
		Доступные команды:  
		/start - Запустить бота
	`
	chatId := ctx.user.GetUserId()
	if chatId == b.config.TG_SUPER_ADMIN {
		text += ""
	}
	msg := tgbotapi.NewMessage(chatId, text)
	msg.ParseMode = "HTML"
	_, err := b.Send(msg)
	if err != nil {
		b.errErrorChan <- err
	}
}

// вывести телеграм id пользователя
func (b *Bot) HB_TelegramId(ctx BotContext, update tgbotapi.Update) {
	chatId := ctx.user.GetUserId()
	msg := tgbotapi.NewMessage(chatId, fmt.Sprintf("Ваш телеграм ID: `%d`", chatId))
	msg.ParseMode = "Markdown"
	_, err := b.Send(msg)
	if err != nil {
		b.errErrorChan <- err
	}
}

// вывести обращения пользователей
func (b *Bot) HB_Feedbacks(ctx BotContext, update tgbotapi.Update) {
	ctx.user.SetAction(serializers.ACTION_FEEDBACK_LIST)
	keyboard, err := b.KeyboardFeedbacks(ctx.user)
	if err != nil {
		b.errErrorChan <- err
	}
	_, err = b.Send(tgbotapi.NewEditMessageTextAndMarkup(
		update.CallbackQuery.Message.Chat.ID,
		update.CallbackQuery.Message.MessageID,
		"Обращения пользователей:",
		*keyboard,
	))
	if err != nil {
		b.errErrorChan <- err
	}
}

// вывести обращение пользователя
func (b *Bot) HB_Feedback(ctx BotContext, update tgbotapi.Update) {
	var v int

	if val := ctx.Value(CTX_KEY_DATA); val != nil {
		data, ok := val.(string)
		if !ok {
			b.errErrorChan <- errors.New("ошибка преобразование в int")
			return
		}
		var err error
		v, err = strconv.Atoi(data)
		if err != nil {
			b.errErrorChan <- err
			return
		}
	} else {
		b.errErrorChan <- errors.New("не найдено значение в контексте")
		return
	}

	opt := infobotdb.NewInfoBotOptions(
		infobotdb.WithId(v),
	)
	feedbacks, _, err := b.DB.Feedbacks(ctx, opt)
	if err != nil {
		b.errErrorChan <- err
		return
	}
	if len(feedbacks) == 0 {
		b.errErrorChan <- er.ErrorNotExist
		return
	}

	feedback := feedbacks[0]
	text := fmt.Sprintf(
		"Имя: %s\nКонтактная информация: %s\nСообщение: %s\nДата обращения: %s",
		feedback.Name,
		feedback.Contact,
		feedback.Message,
		feedback.CreatedAt.Format("2006-01-02 15:04:05"),
	)

	if update.CallbackQuery.Message.Text == text {
		return
	}

	newMsg := tgbotapi.NewEditMessageTextAndMarkup(
		ctx.user.GetUserId(),
		update.CallbackQuery.Message.MessageID,
		text,
		*update.CallbackQuery.Message.ReplyMarkup,
	)
	if _, err := b.Send(newMsg); err != nil {
		b.errErrorChan <- err
	}
}

// удалить сайт
func (b *Bot) HB_DelSite(ctx BotContext, update tgbotapi.Update) {
	actionSite := ctx.user.GetActionSite()
	if actionSite == nil {
		b.HB_Sites(ctx, update)
		return
	}

	err := b.DB.MonitoringSiteDelete(b.ctx, ctx.user.GetUserId(), actionSite.Id)
	if err != nil {
		b.errErrorChan <- err
		return
	}
	ctx.user.SetActionSite(nil)
	b.HB_Sites(ctx, update)
}

// вывести сайты
func (b *Bot) HB_Sites(ctx BotContext, update tgbotapi.Update) {
	ctx.user.SetAction(serializers.ACTION_SITE_LIST)
	ctx.user.SetOffset(0)
	keyboard, err := b.KeyboardSites(ctx.user)
	if err != nil {
		b.errErrorChan <- err
		return
	}
	_, err = b.Send(tgbotapi.NewEditMessageTextAndMarkup(
		update.CallbackQuery.Message.Chat.ID,
		update.CallbackQuery.Message.MessageID,
		"Ваши добавленные сайты:",
		*keyboard,
	))
	if err != nil {
		b.errErrorChan <- err
	}
}

// войти в процесс добавления инфы о сайте
func (b *Bot) HB_SiteAdd(ctx BotContext, update tgbotapi.Update) {
	ctx.user.SetActionSite(nil)
	msg := b.AddSiteMiddleware(ctx.user, func(user *serializers.User) {})
	_, err := b.Send(msg)
	if err != nil {
		b.errErrorChan <- err
	}
}

// войти в процесс обновления инфы о сайте
func (b *Bot) HB_SiteUpdate(ctx BotContext, update tgbotapi.Update) {
	var site_id int

	if val := ctx.Value(CTX_KEY_DATA); val != nil {
		data, ok := val.(string)
		if !ok {
			b.errErrorChan <- errors.New("ошибка преобразование в int")
			return
		}
		var err error
		site_id, err = strconv.Atoi(data)
		if err != nil {
			b.errErrorChan <- err
			return
		}
	} else {
		b.errErrorChan <- errors.New("не найдено значение в контексте")
		return
	}

	ctx.user.SetAction(serializers.ACTION_SITE_UPD)

	opts := infobotdb.NewInfoBotOptions(
		infobotdb.WithUserId(ctx.user.GetUserId()),
		infobotdb.WithId(site_id),
	)
	sites, _, err := b.DB.MonitoringSites(b.ctx, opts)
	if err != nil {
		b.errErrorChan <- err
		return
	}

	ctx.user.SetActionSite(sites[0])
	b.HB_SiteInfoUpdate(ctx, update)
}

// Обновить инфу о сайте
func (b *Bot) HB_SiteInfoUpdate(ctx BotContext, update tgbotapi.Update) {
	chatId := ctx.user.GetUserId()
	site := ctx.user.GetActionSite()

	_, err := b.DB.MonitoringSiteUpdate(ctx, chatId, site)
	if err != nil {
		b.errErrorChan <- err
	}

	_, err = b.Send(tgbotapi.NewEditMessageTextAndMarkup(
		update.CallbackQuery.Message.Chat.ID,
		update.CallbackQuery.Message.MessageID,
		fmt.Sprintf("Обновление информации о сайте %s:", site.Url),
		*KeyboardSiteSettings(site),
	))
	if err != nil {
		b.errErrorChan <- err
	}
}

// Обновление сайта, кнопка "Мониторить"
func (b *Bot) HB_UpdateSiteMonitorYes(ctx BotContext, update tgbotapi.Update) {
	actionSite := ctx.user.GetActionSite()
	actionSite.Monitoring = true
	keyboard := KeyboardSiteSettings(ctx.user.GetActionSite())
	_, err := b.Send(tgbotapi.NewEditMessageReplyMarkup(
		update.CallbackQuery.Message.Chat.ID,
		update.CallbackQuery.Message.MessageID,
		*keyboard,
	))
	if err != nil {
		b.errErrorChan <- err
	}
}

// Обновление сайта, кнопка "Не мониторить"
func (b *Bot) HB_UpdateSiteMonitorNo(ctx BotContext, update tgbotapi.Update) {
	actionSite := ctx.user.GetActionSite()
	actionSite.Monitoring = false
	keyboard := KeyboardSiteSettings(ctx.user.GetActionSite())
	_, err := b.Send(tgbotapi.NewEditMessageReplyMarkup(
		update.CallbackQuery.Message.Chat.ID,
		update.CallbackQuery.Message.MessageID,
		*keyboard,
	))
	if err != nil {
		b.errErrorChan <- err
	}
}

// Обновление сайта, кнопка "Мониторить каждые 10 минут"
func (b *Bot) HB_UpdateSiteMonitorDuration10(ctx BotContext, update tgbotapi.Update) {
	actionSite := ctx.user.GetActionSite()
	actionSite.DurationMinutes = 10
	keyboard := KeyboardSiteSettings(ctx.user.GetActionSite())
	_, err := b.Send(tgbotapi.NewEditMessageReplyMarkup(
		update.CallbackQuery.Message.Chat.ID,
		update.CallbackQuery.Message.MessageID,
		*keyboard,
	))
	if err != nil {
		b.errErrorChan <- err
	}
}

// Обновление сайта, кнопка "Мониторить каждые 15 минут"
func (b *Bot) HB_UpdateSiteMonitorDuration15(ctx BotContext, update tgbotapi.Update) {
	actionSite := ctx.user.GetActionSite()
	actionSite.DurationMinutes = 15
	keyboard := KeyboardSiteSettings(ctx.user.GetActionSite())
	_, err := b.Send(tgbotapi.NewEditMessageReplyMarkup(
		update.CallbackQuery.Message.Chat.ID,
		update.CallbackQuery.Message.MessageID,
		*keyboard,
	))
	if err != nil {
		b.errErrorChan <- err
	}
}

// Обновление сайта, кнопка "Мониторить каждые 20 минут"
func (b *Bot) HB_UpdateSiteMonitorDuration20(ctx BotContext, update tgbotapi.Update) {
	actionSite := ctx.user.GetActionSite()
	actionSite.DurationMinutes = 20
	keyboard := KeyboardSiteSettings(ctx.user.GetActionSite())
	_, err := b.Send(tgbotapi.NewEditMessageReplyMarkup(
		update.CallbackQuery.Message.Chat.ID,
		update.CallbackQuery.Message.MessageID,
		*keyboard,
	))
	if err != nil {
		b.errErrorChan <- err
	}
}
