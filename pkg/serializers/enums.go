package serializers

import (
	"fmt"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type ACTION_TYPE int

const (
	ACTION_NONE ACTION_TYPE = iota

	ACTION_SITE_LIST
	ACTION_SITE_ADD
	ACTION_SITE_DEL
	ACTION_SITE_UPD
	ACTION_SITE_UPD_URL
	ACTION_SITE_UPD_MONITOR
	ACTION_SITE_UPD_MONITOR_DURATION

	ACTION_FEEDBACK_LIST
	ACTION_FEEDBACK
)

type ACTION_BUTTON string

var (
	DATA_BTN_NONE ACTION_BUTTON = ""

	DATA_BTN_FEEDBACK = strconv.Itoa(int(ACTION_FEEDBACK_LIST))

	DATA_BTN_ADD_SITE = strconv.Itoa(int(ACTION_SITE_ADD))
	DATA_BTN_UPD_SITE = strconv.Itoa(int(ACTION_SITE_UPD))
	DATA_BTN_DEL_SITE = strconv.Itoa(int(ACTION_SITE_DEL))

	DATA_BTN_ADD_SITE_MONITOR_YES = fmt.Sprintf("%d_1", ACTION_SITE_UPD_MONITOR)
	DATA_BTN_ADD_SITE_MONITOR_NO  = fmt.Sprintf("%d_2", ACTION_SITE_UPD_MONITOR)

	DATA_BTN_ADD_SITE_MONITOR_DURATION_10 = fmt.Sprintf("%d_10", ACTION_SITE_UPD_MONITOR)
	DATA_BTN_ADD_SITE_MONITOR_DURATION_15 = fmt.Sprintf("%d_15", ACTION_SITE_UPD_MONITOR)
	DATA_BTN_ADD_SITE_MONITOR_DURATION_20 = fmt.Sprintf("%d_20", ACTION_SITE_UPD_MONITOR)

	DATA_BTN_SITES     = "sites"
	DATA_BTN_FEEDBACKS = "feedbacks"
	DATA_BTN_HELP      = "help"
	DATA_BTN_HOMEPAGE  = "homepage"
	DATA_BTN_TG_ID     = "tg_id"

	DATA_BTN_NEXT     = "next"
	DATA_BTN_PREVIOUS = "previous"
)

var (
	BTN_SITES     = tgbotapi.NewInlineKeyboardButtonData("🖥️ Сайты", DATA_BTN_SITES)
	BTN_FEEDBACKS = tgbotapi.NewInlineKeyboardButtonData("📄 Обратные обращения", DATA_BTN_FEEDBACKS)
	BTN_TG_ID     = tgbotapi.NewInlineKeyboardButtonData("🫥 Telegram ID", DATA_BTN_TG_ID)
	BTN_ADD_SITE  = tgbotapi.NewInlineKeyboardButtonData("➕ Добавить сайт", DATA_BTN_ADD_SITE)
	BTN_DEL_SITE  = tgbotapi.NewInlineKeyboardButtonData("❌ Удалить сайт", DATA_BTN_DEL_SITE)
	BTN_HOMEPAGE  = tgbotapi.NewInlineKeyboardButtonData("🏠 Главное меню", DATA_BTN_HOMEPAGE)
	BTN_HELP      = tgbotapi.NewInlineKeyboardButtonData("🆘 Помощь", DATA_BTN_HELP)

	BTN_NEXT     = tgbotapi.NewInlineKeyboardButtonData("Далее ➡️", DATA_BTN_NEXT)
	BTN_PREVIOUS = tgbotapi.NewInlineKeyboardButtonData("⬅️ Предыдущее", DATA_BTN_PREVIOUS)
)
