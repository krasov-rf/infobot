package infobotgrpc

import (
	"context"

	"github.com/krasov-rf/infobot/pkg/grpc/infobotpb"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *Server) Feedback(ctx context.Context, in *infobotpb.FeedbackRequest) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}
