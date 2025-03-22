package infobot

import (
	"errors"
	"fmt"

	"github.com/krasov-rf/infobot/pkg/serializers"
	infobotdb "github.com/krasov-rf/infobot/pkg/storage/infobot"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Главная клавиатура
func (b *Bot) KeyboardHomePage(a *serializers.User) *tgbotapi.InlineKeyboardMarkup {
	var rows [][]tgbotapi.InlineKeyboardButton

	if a.GetChatId() == b.config.TG_SUPER_ADMIN {
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(serializers.BTN_SITES))
	}

	rows = append(rows,
		tgbotapi.NewInlineKeyboardRow(serializers.BTN_TG_ID),
		tgbotapi.NewInlineKeyboardRow(serializers.BTN_HELP),
	)

	keyboard := tgbotapi.NewInlineKeyboardMarkup(rows...)
	return &keyboard
}

// Клавиатура вывода всех сайтов
func (b *Bot) KeyboardFeedbacks(a *serializers.User) (*tgbotapi.InlineKeyboardMarkup, error) {
	site := a.GetActionSite()
	if site == nil {
		return nil, errors.New("не указан сайт")
	}

	offset := a.GetOffset()
	opts := infobotdb.NewInfoBotOptions(
		infobotdb.WithOffset(offset),
		infobotdb.WithSiteId(site.Id),
	)

	feedbacks, cnt, err := b.DB.Feedbacks(b.ctx, opts)
	if err != nil {
		return nil, err
	}

	rows := [][]tgbotapi.InlineKeyboardButton{}
	for _, feedback := range feedbacks {
		rows = append(rows,
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(
					fmt.Sprintf("%s | %s", feedback.Name, feedback.Message),
					fmt.Sprintf("%s|%d", serializers.DATA_BTN_FEEDBACK, feedback.Id)),
			),
		)
	}

	var navigations []tgbotapi.InlineKeyboardButton
	if offset != 0 {
		navigations = append(navigations, serializers.BTN_PREVIOUS)
	}
	if offset+infobotdb.QUERY_LIMIT < cnt {
		navigations = append(navigations, serializers.BTN_NEXT)
	}

	if len(navigations) != 0 {
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(navigations...))
	}

	rows = append(rows,
		tgbotapi.NewInlineKeyboardRow(serializers.BTN_SITES),
		tgbotapi.NewInlineKeyboardRow(serializers.BTN_HOMEPAGE),
	)

	keyboard := tgbotapi.NewInlineKeyboardMarkup(rows...)

	return &keyboard, nil
}

// Клавиатура вывода всех сайтов
func (b *Bot) KeyboardSites(a *serializers.User) (*tgbotapi.InlineKeyboardMarkup, error) {
	offset := a.GetOffset()
	opts := infobotdb.NewInfoBotOptions(
		infobotdb.WithOffset(offset),
		infobotdb.WithUserId(a.GetChatId()),
	)

	sites, cnt, err := b.DB.MonitoringSites(b.ctx, opts)
	if err != nil {
		return nil, err
	}

	rows := [][]tgbotapi.InlineKeyboardButton{}
	for _, site := range sites {
		rows = append(rows,
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(site.Url, fmt.Sprintf("%s|%d", serializers.DATA_BTN_UPD_SITE, site.Id)),
			),
		)
	}

	var navigations []tgbotapi.InlineKeyboardButton
	if offset != 0 {
		navigations = append(navigations, serializers.BTN_PREVIOUS)
	}
	if offset+infobotdb.QUERY_LIMIT < cnt {
		navigations = append(navigations, serializers.BTN_NEXT)
	}

	if len(navigations) != 0 {
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(navigations...))
	}

	rows = append(rows,
		tgbotapi.NewInlineKeyboardRow(serializers.BTN_ADD_SITE),
		tgbotapi.NewInlineKeyboardRow(serializers.BTN_HOMEPAGE),
	)

	keyboard := tgbotapi.NewInlineKeyboardMarkup(rows...)

	return &keyboard, nil
}

// Клавиатура настройки сайта
func KeyboardSiteSettings(a *serializers.SiteSerializer) *tgbotapi.InlineKeyboardMarkup {
	monitorButtonText := "Мониторить"
	if a.Monitoring {
		monitorButtonText = "Мониторить ✅"
	}

	noMonitorButtonText := "Не мониторить"
	if !a.Monitoring {
		noMonitorButtonText = "Не мониторить ✅"
	}

	durationButton10Text := "10 минут"
	if a.DurationMinutes == 10 {
		durationButton10Text += " ✅"
	}

	durationButton15Text := "15 минут"
	if a.DurationMinutes == 15 {
		durationButton15Text += " ✅"
	}

	durationButton20Text := "20 минут"
	if a.DurationMinutes == 20 {
		durationButton20Text += " ✅"
	}

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(serializers.BTN_FEEDBACKS),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(monitorButtonText, serializers.DATA_BTN_ADD_SITE_MONITOR_YES),
			tgbotapi.NewInlineKeyboardButtonData(noMonitorButtonText, serializers.DATA_BTN_ADD_SITE_MONITOR_NO),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(durationButton10Text, serializers.DATA_BTN_ADD_SITE_MONITOR_DURATION_10),
			tgbotapi.NewInlineKeyboardButtonData(durationButton15Text, serializers.DATA_BTN_ADD_SITE_MONITOR_DURATION_15),
			tgbotapi.NewInlineKeyboardButtonData(durationButton20Text, serializers.DATA_BTN_ADD_SITE_MONITOR_DURATION_20),
		),
		tgbotapi.NewInlineKeyboardRow(serializers.BTN_DEL_SITE),
		tgbotapi.NewInlineKeyboardRow(serializers.BTN_HOMEPAGE),
	)
	return &keyboard
}
