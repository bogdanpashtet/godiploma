package grpc

import (
	"context"
	"fmt"
	"net"

	"github.com/bogdanpashtet/godiploma/internal/config"
	"github.com/bogdanpashtet/godiploma/internal/log"

	grpcmiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpcctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	grpcprometheus "github.com/grpc-ecosystem/go-grpc-prometheus"

	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

const HealthCheckName = "grpc-server"

type Registrar interface {
	Register(srv *grpc.Server)
}
type Params struct {
	fx.In

	Logger     log.Logger
	Config     *config.AppConfig
	Registrars []Registrar `group:"grpcRegistrars"`
}

type App struct {
	logger     log.Logger
	gRPCServer *grpc.Server
	cfg        config.GRPCBase
	isReady    bool
}

func New(params Params) *App {
	gRPCServer := grpc.NewServer(grpc.UnaryInterceptor(grpcmiddleware.ChainUnaryServer(
		grpcprometheus.UnaryServerInterceptor,
		grpcctxtags.UnaryServerInterceptor(grpcctxtags.WithFieldExtractor(grpcctxtags.CodeGenRequestFieldExtractor)),
	)))

	// init servers
	for _, reg := range params.Registrars {
		reg.Register(gRPCServer)
	}

	grpcprometheus.EnableClientHandlingTimeHistogram()
	grpcprometheus.EnableHandlingTimeHistogram()
	grpcprometheus.Register(gRPCServer)

	app := &App{
		logger:     params.Logger,
		gRPCServer: gRPCServer,
		cfg:        params.Config.Grpc,
	}

	return app
}

func Register(lc fx.Lifecycle, app *App) {
	lc.Append(fx.Hook{
		OnStart: app.onStart,
		OnStop:  app.onStop,
	})
}

func (a *App) onStart(_ context.Context) error {
	a.logger.Info("starting gRPC server", zap.Int("port", a.cfg.Port))

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", a.cfg.Port))
	if err != nil {
		return err
	}

	go func() {
		if err := a.gRPCServer.Serve(listener); err != nil {
			a.logger.Error("grpc start error", zap.Error(err))
		}

		a.isReady = false
	}()

	a.isReady = true

	return nil
}

func (a *App) onStop(_ context.Context) error {
	a.logger.Info("stopping gRPC server", zap.Int("port", a.cfg.Port))

	a.gRPCServer.GracefulStop()

	return nil
}

func (a *App) HealthCheck(_ context.Context) error {
	if !a.isReady {
		return fmt.Errorf("grpc server is not ready")
	}

	return nil
}

func AsRegistrar(f any) any {
	return fx.Annotate(
		f,
		fx.As(new(Registrar)),
		fx.ResultTags(`group:"grpcRegistrars"`),
	)
}
