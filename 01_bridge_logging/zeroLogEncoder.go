package main

import (
	"context"
	"sync"
	"time"

	"github.com/rs/zerolog"
	"go.uber.org/zap/buffer"
	"go.uber.org/zap/zapcore"
)

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
