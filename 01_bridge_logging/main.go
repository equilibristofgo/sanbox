// Sample app to play with logger and bridging between different implementations
package main

import (
	"context"
	"os"
	"sync"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"go.uber.org/zap"
	"go.uber.org/zap/buffer"
	"go.uber.org/zap/zapcore"
)

const (
	REQUEST_ID = "request_id"
)

// With this, try to hook every print with zerolog
type SeverityHook struct{}

// Need to implement this
func (h SeverityHook) Run(e *zerolog.Event, level zerolog.Level, msg string) {
	if level != zerolog.NoLevel {
		e.Str("severity", level.String())
	}
}

// With this, try to create a encoder for zap logger that access to the fields that are injected
type zeroLogEncoder struct {
	zapcore.Encoder

	ctx           context.Context `json:"ctx,omitempty"`
	zeroLogLogger *zerolog.Logger `json:"zero_log_logger,omitempty"`

	event *zerolog.Event `json:"event,omitempty"`
}

// A simple implementation
func (e *zeroLogEncoder) Clone() zapcore.Encoder {
	clone := sync.Pool{New: func() interface{} {
		return &zeroLogEncoder{e.Encoder, e.ctx, e.zeroLogLogger, e.event}
	}}

	return clone.Get().(*zeroLogEncoder)
}

// Needed for the fields added in initialization
func (e *zeroLogEncoder) AddString(key, value string) {
	e.event = e.event.Str(key, value)
}

func (e *zeroLogEncoder) EncodeEntry(entry zapcore.Entry, fields []zapcore.Field) (*buffer.Buffer, error) {
	for _, f := range fields {
		// TODO Check types
		switch f.Type {
		case zapcore.StringType:
			// intCtx := context.WithValue(*e.ctx, propagateKey, &Values{Key: f.Key, Value: f.String}))
			// e.ctx = &intCtx

			e.event = e.event.Str(f.Key, f.String)
		case zapcore.Int16Type:
		case zapcore.Int32Type:
		case zapcore.Int64Type:
			e.event = e.event.Int64(f.Key, f.Integer)

		case zapcore.DurationType:
			e.event = e.event.Dur(f.Key, time.Duration(f.Integer))
		default:
		}
	}

	e.event.Msg("ZEROLOG LOGGER - Generic message")

	return e.Encoder.EncodeEntry(entry, fields)
}

type (
	// contextKey is an unexported type used as key for items stored in the
	// Context object
	contextKey struct{}

	// propagator implements the custom context propagator
	propagator struct{}

	// Values is a struct holding values
	Values struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	}
)

var propagateKey = contextKey{}

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

func GetKeyFromCtx(ctx context.Context, key string) string {
	if val := ctx.Value(propagateKey); val != nil {
		vals, ok := val.(Values)
		if !ok {
			valsPtr, ok := val.(*Values)
			if !ok {
				return ""
			}
			if valsPtr.Key == key {
				return valsPtr.Value
			}
		}
		if vals.Key == key {
			return vals.Value
		}
	}
	return ""
}

func logRechargeEngineContextFieldsFromActivity(ctx context.Context, logEvent *zerolog.Event) *zerolog.Event {
	id := GetKeyFromCtx(ctx, REQUEST_ID)
	if id != "" {
		logEvent = logEvent.Str(REQUEST_ID, id)
	}
	return logEvent
}

// equilibristofgo/sandbox/01_bridge_logging$ gomarkdoc --output doc.md .
func main() {

	ctx := context.Background()

	// Inject generated id from ctx in original request
	ctx = context.WithValue(ctx, propagateKey, &Values{Key: REQUEST_ID, Value: "1234567890"})

	zeroLogLogger := log.With().Str("module", "RechargeFrontModule").Logger()

	zapLogger, _ := zap.NewProduction()

	defer zapLogger.Sync()
	zapLogger.Info("failed to fetch URL",
		zap.String("url", "http://example.com"),
		zap.Int("attempt", 3),
		zap.Duration("backoff", time.Second),
	)

	enc := &zeroLogEncoder{
		Encoder:       zapcore.NewJSONEncoder(zap.NewProductionConfig().EncoderConfig), // NewConsoleEncoder
		ctx:           ctx,
		zeroLogLogger: &zeroLogLogger,
		event:         zeroLogLogger.Info(),
	}

	f, _ := os.OpenFile(os.DevNull, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	defer f.Close()

	zapCore := zapcore.NewCore(
		enc,
		f,
		zapcore.InfoLevel,
	)

	zapLoggerCustomEncoder := zap.New(
		zapCore,
		// this mimics the behavior of NewProductionConfig.Build
		zap.ErrorOutput(os.Stderr),
	)

	zapLoggerCustomEncoder = zapLoggerCustomEncoder.With(
		zapcore.Field{Key: "tagWorkflowType", Type: zapcore.StringType, String: "XXX"},
		zapcore.Field{Key: "tagWorkflowID", Type: zapcore.StringType, String: "XXX"},
		zapcore.Field{Key: "tagRunID", Type: zapcore.StringType, String: "XXX"},
	).WithOptions(zap.WrapCore(func(c zapcore.Core) zapcore.Core { return zapCore }))

	zapLoggerCustomEncoder.Info("ZAP LOGGER - failed to fetch URL",
		zap.String("url", "http://example.com"),
		zap.Int("attempt", 3),
		zap.Duration("backoff", time.Second),
	)

	zeroLogLogger.Info().Msg("ZEROLOG LOGGER - failed to fetch URL")

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	hooked := log.Hook(SeverityHook{})
	hooked.Warn().Msg("")

	log.Debug().
		Str("Scale", "833 cents").
		Float64("Interval", 833.09).
		Msg("Fibonacci is everywhere")

	// // https://github.com/busimus/gocutelog/blob/master/gocutelog.go
	// w := gocutelog.NewWriter("localhost:19996", "json")
	// l := zerolog.New(w)
	// l.Info().Msg("Hello world from zerolog!")

	loggerContextualFromZAP := ContextualLoggerBridge{
		zeroLogLogger,
		*zapLoggerCustomEncoder,
		zeroLogLogger.Sample(&zerolog.BasicSampler{N: 10}),
		logRechargeEngineContextFieldsFromActivity,
	}

	loggerContextualFromZAP.InfoWithContext(enc.ctx).Bool("Result", false).Str("Tenant", "lyca").Msg("This was a sample with ZEROLOG")

	// Adding again ... simulating cadence internals?
	zapLoggerCustomEncoder = zapLoggerCustomEncoder.With(
		zapcore.Field{Key: "tagWorkflowType", Type: zapcore.StringType, String: "YYY"},
		zapcore.Field{Key: "tagWorkflowID", Type: zapcore.StringType, String: "YYY"},
		zapcore.Field{Key: "tagRunID", Type: zapcore.StringType, String: "YYY"},
	).WithOptions(zap.WrapCore(func(c zapcore.Core) zapcore.Core { return zapCore }))

	loggerContextualFromZAP.InfoWithContext(enc.ctx).Bool("Result", false).Str("Tenant", "lyca").Msg("This was a sample with ZAP")

}
