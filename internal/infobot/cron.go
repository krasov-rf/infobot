package infobot

import (
	"fmt"
	"net/http"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *Bot) InitializeCron() error {
	_, err := b.cron.AddFunc("@every 1m", func() {
		sites, err := b.DB.MonitoringSitesForCheck(b.ctx)
		if err != nil {
			b.errChan <- err
			return
		}

		for _, site := range sites {
			resp, err := http.Get(site.Url)
			if err != nil {
				b.errChan <- err
				return
			}
			defer resp.Body.Close()

			err = b.DB.MonitoringSiteStatusUpdate(b.ctx, site.Id, resp.StatusCode)
			if err != nil {
				b.errChan <- err
				return
			}

			if site.StatusCode == resp.StatusCode {
				continue
			}

			for _, user_id := range site.TgUsers {
				_, err = b.Send(tgbotapi.NewMessage(
					user_id,
					fmt.Sprintf("Новый код %d у страницы %s", resp.StatusCode, site.Url),
				))
				if err != nil {
					b.errChan <- err
				}
			}
		}
	})
	if err != nil {
		return err
	}

	return nil
}
