package v1

import (
	"context"
	"fmt"
	"time"

	filed "github.com/bogdanpashtet/godiploma/internal/domain/file"
	filev1 "github.com/bogdanpashtet/godiploma/protos/gen/go/client/godiploma/file/v1"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

const kind = "FileServer"

type Params struct {
	fx.In

	Logger *zap.Logger
}
type Server struct {
	l *zap.Logger

	filev1.Un
}

func NewServer(params Params) *Server {
	return &Server{
		l: params.Logger,
	}
}

func (s *Server) Register(gRPCServer *grpc.Server) {
	filev1.RegisterFileServiceServer(gRPCServer, s)
}

func (s *Server) UploadDocuments(ctx context.Context, req *filev1.UploadDocumentsRequest) (*filev1.UploadDocumentsResponse, error) {
	s.l = s.l.With(zap.String("request_id", req.GetRequestId()))

	return
}
