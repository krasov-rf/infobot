package infobotgrpc

import (
	"context"
	"log"
	"net"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/krasov-rf/infobot/internal/infobot"
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
	config  *settings.Config
	ctx     context.Context
	errChan chan error

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
		ctx:    context.Background(),
	}, nil
}

func (s *Server) Run() {
	opts := []grpc.ServerOption{}

	s.errChan = make(chan error, 10)
	defer close(s.errChan)
	go s.errorListener()

	grpc.UnaryInterceptor(s.ensureValidBasicCredentials)
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

func (s *Server) ensureValidBasicCredentials(
	ctx context.Context, req any,
	info *grpc.UnaryServerInfo, handler grpc.UnaryHandler,
) (any, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.InvalidArgument, "отсутствует метадата")
	}

	notAuthorized := status.Errorf(codes.Unauthenticated, "авторизация не пройдена")

	auth := md["authorization"]
	if len(auth) != 2 {
		return nil, notAuthorized
	}

	sites, _, err := s.DB.MonitoringSites(ctx, infobotdb.NewInfoBotOptions(
		infobotdb.WithDomain(auth[0]),
	))
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "авторизация не пройдена, ошибка: %v", err)
	}
	if len(sites) == 0 {
		return nil, notAuthorized
	}

	site := sites[0]

	if site.SecretKey != auth[1] {
		return nil, notAuthorized
	}

	ctx = context.WithValue(ctx, infobot.CTX_KEY_SITE, *site)

	return handler(ctx, req)
}
