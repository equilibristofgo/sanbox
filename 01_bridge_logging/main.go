// Sample app to play with logger and bridging between different implementations
package main

import (
	"context"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	REQUEST_ID = "request_id"
)

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
