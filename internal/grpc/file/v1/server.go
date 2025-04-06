package v1

import (
	"context"

	filev1 "github.com/bogdanpashtet/godiploma/protos/gen/go/client/godiploma/file/v1"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Params struct {
	fx.In

	Logger *zap.Logger
	Svc    Service
}
type Server struct {
	l   *zap.Logger
	svc Service

	filev1.UnimplementedFileServiceServer
}

func NewServer(params Params) *Server {
	return &Server{
		l:   params.Logger,
		svc: params.Svc,
	}
}

func (s *Server) Register(gRPCServer *grpc.Server) {
	filev1.RegisterFileServiceServer(gRPCServer, s)
}

func (s *Server) UploadFiles(ctx context.Context, req *filev1.UploadFilesRequest) (*filev1.UploadFilesResponse, error) {
	s.l = s.l.With(zap.String("request_id", req.GetRequestId()))

	if err := s.svc.UploadFiles(ctx, convertUploadRequestToDomain(req)); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to upload documents: %s", err.Error())
	}

	return &filev1.UploadFilesResponse{}, nil
}
