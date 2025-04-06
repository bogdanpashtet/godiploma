package v1

import (
	"context"

	filed "github.com/bogdanpashtet/godiploma/internal/domain/file"
	server "github.com/bogdanpashtet/godiploma/internal/grpc"
	cipherv1 "github.com/bogdanpashtet/godiploma/protos/gen/go/client/godiploma/cipher/v1"
	"github.com/samber/lo"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type Params struct {
	fx.In

	Logger *zap.Logger
	Svc    Service
}
type Server struct {
	l   *zap.Logger
	svc Service

	cipherv1.UnimplementedCipherServiceServer
}

func NewServer(params Params) *Server {
	return &Server{
		l:   params.Logger,
		svc: params.Svc,
	}
}

func (s *Server) Register(gRPCServer *grpc.Server) {
	cipherv1.RegisterCipherServiceServer(gRPCServer, s)
}

func (s *Server) CreateStegoImage(ctx context.Context, req *cipherv1.CreateStegoImageRequest) (*cipherv1.CreateStegoImageResponse, error) {
	s.l = s.l.With(zap.String("request_id", req.GetRequestId()))

	res, err := s.svc.CreateStegoImage(ctx, convertCreateStegoImageToDomain(req))
	if err != nil {
		return nil, server.ErrorFromDomain(err)
	}

	return &cipherv1.CreateStegoImageResponse{
		Files: lo.Map(res, func(file filed.File, _ int) *cipherv1.File {
			return &cipherv1.File{
				Metadata: &cipherv1.Metadata{
					Type: file.Metadata.Type.FromDomain(),
				},
				DocumentData: file.File,
			}
		}),
	}, nil
}
