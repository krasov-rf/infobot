package infobotgrpc

import (
	"context"
	"errors"
	"fmt"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/krasov-rf/infobot/internal/infobot"
	"github.com/krasov-rf/infobot/pkg/grpc/infobotpb"
	"github.com/krasov-rf/infobot/pkg/serializers"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *Server) Feedback(ctx context.Context, in *infobotpb.FeedbackRequest) (*emptypb.Empty, error) {
	var site serializers.SiteSerializer
	if val := ctx.Value(infobot.CTX_KEY_SITE); val != nil {
		var ok bool
		site, ok = val.(serializers.SiteSerializer)
		if !ok {
			err := errors.New("ошибка преобразование в serializers.SiteSerializer")
			s.errChan <- err
			return nil, err
		}
	} else {
		err := errors.New("не найдено значение сайта в контексте")
		s.errChan <- err
		return nil, err
	}

	err := s.DB.FeedbackInsert(ctx, &serializers.FeedbackSerializer{
		SiteId:      site.Id,
		Name:        in.Name,
		Contact:     in.Contact,
		Message:     in.Message,
		FeedbackUrl: in.FeedbackUrl,
	})
	if err != nil {
		s.errChan <- err
		return nil, err
	}

	created_at := time.Now()

	d, err := s.DB.RelatedUsersBySites(ctx, int64(site.Id))
	if err != nil {
		s.errChan <- err
		return nil, err
	}

	for _, user := range d[site.Id] {
		_, err = s.BotAPI.Send(tgbotapi.NewMessage(
			int64(user),
			fmt.Sprintf(
				`Новое обращение\n\n
				Имя: %s\n
				Контактная информация: %s\n
				Сообщение: %s\n
				Прилетело со страницы: %s\n
				Дата обращения: %s`,
				in.Name,
				in.Contact,
				in.Message,
				in.FeedbackUrl,
				created_at,
			),
		))
		if err != nil {
			s.errChan <- err
			return nil, err
		}
	}

	return &emptypb.Empty{}, nil
}
