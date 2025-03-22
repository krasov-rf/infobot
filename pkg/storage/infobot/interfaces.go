package infobotdb

import (
	"context"

	"github.com/krasov-rf/infobot/pkg/serializers"
	"github.com/krasov-rf/infobot/pkg/storage"
)

// клиентское api
type IInfoBotDB interface {
	storage.IStorage

	// добавить нового телеграм пользователя
	TelegramUserRegister(ctx context.Context, user *serializers.UserSerializer) error
	// получить телеграм пользователя
	TelegramUserGet(ctx context.Context, user_id int64) (*serializers.UserSerializer, error)

	// Вывести пользовательские обращения
	Feedbacks(ctx context.Context, opt *OptionsInfoBot) ([]*serializers.FeedbackSerializer, int, error)
	// Добавить пользовательское обращение
	FeedbackInsert(ctx context.Context, user *serializers.FeedbackSerializer) error

	// вывести сайты в мониторинге
	MonitoringSites(ctx context.Context, opt *OptionsInfoBot) ([]*serializers.SiteSerializer, int, error)
	// удалить сайт из мониторинга
	MonitoringSiteDelete(ctx context.Context, user_id int64, site_id int) error
	// добавить сайт в мониторинг
	MonitoringSiteAdd(ctx context.Context, user_id int64, site_url string, working bool, status_code int) (*serializers.SiteSerializer, error)
	// обновить сайт
	MonitoringSiteUpdate(ctx context.Context, user_id int64, site *serializers.SiteSerializer) (*serializers.SiteSerializer, error)
}
