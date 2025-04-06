package log

import (
	"strings"

	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
)

type FxLogger struct {
	*zap.Logger
}

func (l *FxLogger) LogEvent(event fxevent.Event) { //nolint:funlen,gocyclo
	switch e := event.(type) {
	case *fxevent.OnStartExecuting:
		l.Sugar().Infow("OnStart hook executing",
			zap.String("callee", e.FunctionName),
			zap.String("caller", e.CallerName),
		)
	case *fxevent.OnStartExecuted:
		if e.Err != nil {
			l.Sugar().Errorw("OnStart hook failed",
				zap.String("callee", e.FunctionName),
				zap.String("caller", e.CallerName),
				zap.Error(e.Err),
			)
		} else {
			l.Sugar().Infow("OnStart hook executed",
				zap.String("callee", e.FunctionName),
				zap.String("caller", e.CallerName),
				zap.String("runtime", e.Runtime.String()),
			)
		}
	case *fxevent.OnStopExecuting:
		l.Sugar().Infow("OnStop hook executing",
			zap.String("callee", e.FunctionName),
			zap.String("caller", e.CallerName),
		)
	case *fxevent.OnStopExecuted:
		if e.Err != nil {
			l.Sugar().Errorw("OnStop hook failed",
				zap.String("callee", e.FunctionName),
				zap.String("caller", e.CallerName),
				zap.Error(e.Err),
			)
		} else {
			l.Sugar().Infow("OnStop hook executed",
				zap.String("callee", e.FunctionName),
				zap.String("caller", e.CallerName),
				zap.String("runtime", e.Runtime.String()),
			)
		}
	case *fxevent.Supplied:
		if e.Err != nil {
			l.Sugar().Errorw("error encountered while applying options",
				zap.String("type", e.TypeName),
				zap.Strings("stacktrace", e.StackTrace),
				moduleField(e.ModuleName),
				zap.Error(e.Err))
		} else {
			l.Sugar().Infow("supplied",
				zap.String("type", e.TypeName),
				zap.Strings("stacktrace", e.StackTrace),
				moduleField(e.ModuleName),
			)
		}
	case *fxevent.Provided:
		for _, rtype := range e.OutputTypeNames {
			l.Sugar().Infow("provided",
				zap.String("constructor", e.ConstructorName),
				zap.Strings("stacktrace", e.StackTrace),
				moduleField(e.ModuleName),
				zap.String("type", rtype),
				maybeBool("private", e.Private),
			)
		}
		if e.Err != nil {
			l.Sugar().Errorw("error encountered while applying options",
				moduleField(e.ModuleName),
				zap.Strings("stacktrace", e.StackTrace),
				zap.Error(e.Err))
		}
	case *fxevent.Replaced:
		for _, rtype := range e.OutputTypeNames {
			l.Sugar().Infow("replaced",
				zap.Strings("stacktrace", e.StackTrace),
				moduleField(e.ModuleName),
				zap.String("type", rtype),
			)
		}
		if e.Err != nil {
			l.Sugar().Errorw("error encountered while replacing",
				zap.Strings("stacktrace", e.StackTrace),
				moduleField(e.ModuleName),
				zap.Error(e.Err))
		}
	case *fxevent.Decorated:
		for _, rtype := range e.OutputTypeNames {
			l.Sugar().Infow("decorated",
				zap.String("decorator", e.DecoratorName),
				zap.Strings("stacktrace", e.StackTrace),
				moduleField(e.ModuleName),
				zap.String("type", rtype),
			)
		}
		if e.Err != nil {
			l.Sugar().Errorw("error encountered while applying options",
				zap.Strings("stacktrace", e.StackTrace),
				moduleField(e.ModuleName),
				zap.Error(e.Err))
		}
	case *fxevent.Run:
		if e.Err != nil {
			l.Sugar().Errorw("error returned",
				zap.String("name", e.Name),
				zap.String("kind", e.Kind),
				moduleField(e.ModuleName),
				zap.Error(e.Err),
			)
		} else {
			l.Sugar().Infow("run",
				zap.String("name", e.Name),
				zap.String("kind", e.Kind),
				moduleField(e.ModuleName),
			)
		}
	case *fxevent.Invoking:
		// Do not log stack as it will make logs hard to read.
		l.Sugar().Infow("invoking",
			zap.String("function", e.FunctionName),
			moduleField(e.ModuleName),
		)
	case *fxevent.Invoked:
		if e.Err != nil {
			l.Sugar().Errorw("invoke failed",
				zap.Error(e.Err),
				zap.String("stack", e.Trace),
				zap.String("function", e.FunctionName),
				moduleField(e.ModuleName),
			)
		}
	case *fxevent.Stopping:
		l.Sugar().Infow("received signal",
			zap.String("signal", strings.ToUpper(e.Signal.String())))
	case *fxevent.Stopped:
		if e.Err != nil {
			l.Sugar().Errorw("stop failed", zap.Error(e.Err))
		}
	case *fxevent.RollingBack:
		l.Sugar().Errorw("start failed, rolling back", zap.Error(e.StartErr))
	case *fxevent.RolledBack:
		if e.Err != nil {
			l.Sugar().Errorw("rollback failed", zap.Error(e.Err))
		}
	case *fxevent.Started:
		if e.Err != nil {
			l.Sugar().Errorw("start failed", zap.Error(e.Err))
		} else {
			l.Sugar().Infow("started")
		}
	case *fxevent.LoggerInitialized:
		if e.Err != nil {
			l.Sugar().Errorw("custom logger initialization failed", zap.Error(e.Err))
		} else {
			l.Sugar().Infow("initialized custom fxevent.Logger", zap.String("function", e.ConstructorName))
		}
	}
}

func moduleField(name string) zap.Field {
	if name == "" {
		return zap.Skip()
	}

	return zap.String("module", name)
}

func maybeBool(name string, b bool) zap.Field {
	if b {
		return zap.Bool(name, true)
	}

	return zap.Skip()
}
