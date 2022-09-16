package main

import (
	"context"

	"github.com/rs/zerolog"
	"go.uber.org/zap"
)

type ContextualLoggerBridge struct {
	zerolog.Logger
	zapLogger      zap.Logger
	sampled        zerolog.Logger
	logContextFunc func(context.Context, *zerolog.Event) *zerolog.Event
}

func (l ContextualLoggerBridge) genericLogWithContext(ctx context.Context, logFunc func() *zerolog.Event) *zerolog.Event {
	l.zapLogger.Info("")

	if l.logContextFunc != nil {
		return l.logContextFunc(ctx, logFunc())
	}
	return logFunc()
}

func (l ContextualLoggerBridge) TraceWithContext(ctx context.Context) *zerolog.Event {
	return l.genericLogWithContext(ctx, l.Trace)
}

func (l ContextualLoggerBridge) DebugWithContext(ctx context.Context) *zerolog.Event {
	return l.genericLogWithContext(ctx, l.Debug)
}

func (l ContextualLoggerBridge) InfoWithContext(ctx context.Context) *zerolog.Event {
	return l.genericLogWithContext(ctx, l.Info)
}

func (l ContextualLoggerBridge) WarnWithContext(ctx context.Context) *zerolog.Event {
	return l.genericLogWithContext(ctx, l.Warn)
}

func (l ContextualLoggerBridge) ErrorWithContext(ctx context.Context) *zerolog.Event {
	return l.genericLogWithContext(ctx, l.Error)
}

func (l ContextualLoggerBridge) FatalWithContext(ctx context.Context) *zerolog.Event {
	return l.genericLogWithContext(ctx, l.Fatal)
}

func (l ContextualLoggerBridge) SampledTrace(ctx context.Context) *zerolog.Event {
	return l.genericLogWithContext(ctx, l.sampled.Trace)
}

func (l ContextualLoggerBridge) SampledDebug(ctx context.Context) *zerolog.Event {
	return l.genericLogWithContext(ctx, l.sampled.Debug)
}

func (l ContextualLoggerBridge) SampledInfo(ctx context.Context) *zerolog.Event {
	return l.genericLogWithContext(ctx, l.sampled.Info)
}

func (l ContextualLoggerBridge) SampledWarn(ctx context.Context) *zerolog.Event {
	return l.genericLogWithContext(ctx, l.sampled.Warn)
}

func (l ContextualLoggerBridge) SampledError(ctx context.Context) *zerolog.Event {
	return l.genericLogWithContext(ctx, l.sampled.Error)
}
