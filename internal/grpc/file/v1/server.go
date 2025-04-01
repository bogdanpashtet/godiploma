package v1

import (
	"context"

	"github.com/bogdanpashtet/godiploma/internal/log"
	"go.uber.org/fx"
	"google.golang.org/grpc"

	filev1 "github.com/bogdanpashtet/godiploma/protos/gen/go/client/godiploma/file/v1"
)

const kind = "FileServer"

type Params struct {
	fx.In

	//Service Service
	Logger log.Logger
}
type Server struct {
	//svc Service
	l log.Logger

	filev1.UnimplementedFileServiceServer
}

func NewServer(params Params) *Server {
	return &Server{
		//svc: params.Service,
		l: params.Logger,
	}
}

func (s *Server) Register(gRPCServer *grpc.Server) {
	filev1.RegisterFileServiceServer(gRPCServer, s)
}

func (s *Server) UploadDocuments(_ context.Context, _ *filev1.UploadDocumentsRequest) (*filev1.UploadDocumentsResponse, error) {
	s.l.Info("hello from UploadDocuments")
	return &filev1.UploadDocumentsResponse{}, nil
}
