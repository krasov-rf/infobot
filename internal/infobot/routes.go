package infobot

import "github.com/krasov-rf/infobot/pkg/serializers"

func (b *Bot) InitializeRoutes() {
	b.Use(b.RegisterUserMiddleware)

	b.RouteCallback(serializers.DATA_BTN_NEXT, b.HB_UpdateOffsetNext)
	b.RouteCallback(serializers.DATA_BTN_PREVIOUS, b.HB_UpdateOffsetPrevious)
	b.RouteCallback(serializers.DATA_BTN_HOMEPAGE, b.HB_HomePage)

	b.RouteCallback(serializers.DATA_BTN_SITES, b.HB_Sites)
	b.RouteCallback(serializers.DATA_BTN_ADD_SITE, b.HB_SiteAdd)
	b.RouteCallback(serializers.DATA_BTN_UPD_SITE, b.HB_SiteUpdate)
	b.RouteCallback(serializers.DATA_BTN_DEL_SITE, b.HB_DelSite)
	b.RouteCallback(serializers.DATA_BTN_ADD_SITE_MONITOR_YES, b.HB_UpdateSiteMonitorYes)
	b.RouteCallback(serializers.DATA_BTN_ADD_SITE_MONITOR_NO, b.HB_UpdateSiteMonitorNo)
	b.RouteCallback(serializers.DATA_BTN_ADD_SITE_MONITOR_DURATION_10, b.HB_UpdateSiteMonitorDuration10)
	b.RouteCallback(serializers.DATA_BTN_ADD_SITE_MONITOR_DURATION_15, b.HB_UpdateSiteMonitorDuration15)
	b.RouteCallback(serializers.DATA_BTN_ADD_SITE_MONITOR_DURATION_20, b.HB_UpdateSiteMonitorDuration20)

	b.RouteCallback(serializers.DATA_BTN_FEEDBACKS, b.HB_Feedbacks)
	b.RouteCallback(serializers.DATA_BTN_FEEDBACK, b.HB_Feedback)

	b.RouteCallback(serializers.DATA_BTN_TG_ID, b.HB_TelegramId)
	b.RouteCallback(serializers.DATA_BTN_HELP, b.HB_Help)

	b.RouteMessage("start", b.MSG_HomePage)

	b.RouteRawMessage(serializers.ACTION_SITE_ADD_URL, b.HB_SiteAddUrl)
}
