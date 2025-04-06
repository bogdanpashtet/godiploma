package grpc

import (
	"context"
	"fmt"
	"net"

	"github.com/bogdanpashtet/godiploma/internal/config"
	grpcmw "github.com/grpc-ecosystem/go-grpc-middleware"
	loggingmw "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	tagsmw "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	authmw "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/auth"
	recoverymw "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	validatormw "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/validator"
	prometheusmw "github.com/grpc-ecosystem/go-grpc-prometheus"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type Registrar interface {
	Register(srv *grpc.Server)
}
type Params struct {
	fx.In

	Logger     *zap.Logger
	Config     *config.AppConfig
	Registrars []Registrar `group:"grpcRegistrars"`
}

type App struct {
	logger     *zap.Logger
	gRPCServer *grpc.Server
	cfg        config.GRPCBase
	isReady    bool
}

func New(params Params) *App {
	gRPCServer := grpc.NewServer(grpc.UnaryInterceptor(grpcmw.ChainUnaryServer(
		recoverymw.UnaryServerInterceptor(),
		tagsmw.UnaryServerInterceptor(tagsmw.WithFieldExtractor(tagsmw.CodeGenRequestFieldExtractor)),
		prometheusmw.UnaryServerInterceptor,
		loggingmw.UnaryServerInterceptor(params.Logger),
		authmw.UnaryServerInterceptor(func(ctx context.Context) (context.Context, error) { return ctx, nil }), // TODO: add auth
		validatormw.UnaryServerInterceptor(),
	)))

	for _, reg := range params.Registrars {
		reg.Register(gRPCServer)
	}

	prometheusmw.EnableClientHandlingTimeHistogram()
	prometheusmw.EnableHandlingTimeHistogram()
	prometheusmw.Register(gRPCServer)

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
	a.logger.Sugar().Infow("starting gRPC server", zap.Int("port", a.cfg.Port))

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", a.cfg.Port))
	if err != nil {
		return err
	}

	go func() {
		if err = a.gRPCServer.Serve(listener); err != nil {
			a.logger.Sugar().Errorw("grpc start error", zap.Error(err))
		}

		a.isReady = false
	}()

	a.isReady = true

	return nil
}

func (a *App) onStop(_ context.Context) error {
	a.logger.Sugar().Infow("stopping gRPC server", zap.Int("port", a.cfg.Port))

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
