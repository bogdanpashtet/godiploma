package app

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/bogdanpashtet/godiploma/internal/app/grpc"
	"github.com/bogdanpashtet/godiploma/internal/config"
	filev1 "github.com/bogdanpashtet/godiploma/internal/grpc/file/v1"
	"github.com/bogdanpashtet/godiploma/internal/log"
	healthgo "github.com/hellofresh/health-go/v5"

	"go.uber.org/automaxprocs/maxprocs"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
)

const appTimeout = 30 * time.Second

func NewApp() *fx.App {
	app := fx.New(
		DependenciesGraph(),
		fx.WithLogger(func(logger log.Logger) fxevent.Logger {
			return &log.FxLogger{Logger: logger}
		}),
		fx.StopTimeout(appTimeout),
		fx.Invoke(
			health,
			metrics,
			grpc.Register,
		),
	)

	return app
}

func DependenciesGraph() fx.Option {
	return fx.Options(
		fx.Provide(
			context.Background,
			config.New,
			logger,
			grpc.AsRegistrar(filev1.NewServer),
			grpc.New,
		),
	)
}

func logger(lc fx.Lifecycle) (log.Logger, error) {
	lg, err := zap.NewProduction()
	if err != nil {
		return nil, err
	}

	if _, err := maxprocs.Set(maxprocs.Logger(lg.Sugar().Infof)); err != nil {
		lg.Warn("Set maxprocs", zap.Error(err))
	}

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			return lg.Sync()
		},
	})

	return lg, nil
}

func health(lc fx.Lifecycle, l log.Logger, cfg *config.AppConfig) {
	h, _ := healthgo.New(healthgo.WithComponent(healthgo.Component{
		Name:    cfg.AppName,
		Version: cfg.Version,
	}))

	mux := http.NewServeMux()
	mux.Handle(cfg.HTTPHealth.LiveEndpoint, h.Handler())
	mux.Handle(cfg.HTTPHealth.ReadyEndpoint, h.Handler())

	srv := &http.Server{
		Addr:              fmt.Sprintf(":%d", cfg.HTTPHealth.Port),
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			l.Info("Starting health check server", zap.String("address", srv.Addr))
			go func() {
				if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
					l.Error("Health server error", zap.Error(err))
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			l.Info("Shutting down health check server")
			return srv.Shutdown(ctx)
		},
	})
}

func metrics(lc fx.Lifecycle, l log.Logger, cfg *config.AppConfig) {
	srv := newMetricsServer(cfg.Metrics)
	lc.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			l.Info("Starting metrics server", zap.String("address", srv.Addr))
			go func() {
				if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
					l.Error("Metrics server error", zap.Error(err))
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			l.Info("Shutting down metrics check server")
			return srv.Shutdown(ctx)
		},
	})
}
