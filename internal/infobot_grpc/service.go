package infobotgrpc

import (
	"context"
	"encoding/base64"
	"log"
	"net"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	settings "github.com/krasov-rf/infobot/pkg/settings/infobot"
	infobotdb "github.com/krasov-rf/infobot/pkg/storage/infobot"
	infobotdb_pg "github.com/krasov-rf/infobot/pkg/storage/infobot/postgres"

	"github.com/krasov-rf/infobot/pkg/grpc/infobotpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type Server struct {
	*tgbotapi.BotAPI
	infobotpb.UnimplementedInfoBotServiceServer
	config *settings.Config

	DB infobotdb.IInfoBotDB
}

func New(c *settings.Config) (*Server, error) {
	bot, err := tgbotapi.NewBotAPI(c.TG_TOKEN)
	if err != nil {
		return nil, err
	}

	db, err := infobotdb_pg.New(c)
	if err != nil {
		return nil, err
	}

	return &Server{
		BotAPI: bot,
		DB:     db,
		config: c,
	}, nil
}

func (s *Server) Run() {
	opts := []grpc.ServerOption{}

	grpc.UnaryInterceptor(ensureValidBasicCredentials)
	grpc_server := grpc.NewServer(opts...)
	infobotpb.RegisterInfoBotServiceServer(grpc_server, s)

	lis, err := net.Listen("tcp", "5001")
	if err != nil {
		log.Fatalf("ошибка прослушивания порта: %v", err)
	}

	if err := grpc_server.Serve(lis); err != nil {
		log.Fatalf("ошибка сервера: %v", err)
	}
}

func ensureValidBasicCredentials(
	ctx context.Context, req any,
	info *grpc.UnaryServerInfo, handler grpc.UnaryHandler,
) (any, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.InvalidArgument, "отсутствует метадата")
	}
	if !valid(md["authorization"]) {
		return nil, status.Errorf(codes.Unauthenticated, "авторизация не пройдена")
	}

	return handler(ctx, req)
}

func valid(authorization []string) bool {
	if len(authorization) < 1 {
		return false
	}
	token := strings.TrimPrefix(authorization[0], "Basic ")
	return token == base64.StdEncoding.EncodeToString([]byte("admin:admin"))
}
