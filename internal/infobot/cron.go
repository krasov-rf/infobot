package infobot

func (b *Bot) InitializeCron() error {
	_, err := b.cron.AddFunc("@every 1m", func() {

	})
	if err != nil {
		return err
	}

	return nil
}
