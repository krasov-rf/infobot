package infobotgrpc

import (
	"context"
	"errors"

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
			s.errErrorChan <- err
			return nil, err
		}
	} else {
		err := errors.New("не найдено значение сайта в контексте")
		s.errErrorChan <- err
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
		s.errErrorChan <- err
		return nil, err
	}

	msg := tgbotapi.NewMessage(b.config.TG_SUPER_ADMIN, "Bot started")
	_, err = s.BotAPI.Send(msg)
	if err != nil {
		s.errErrorChan <- err
		return nil, err
	}
	return &emptypb.Empty{}, nil
}
