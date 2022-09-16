package main

import "github.com/rs/zerolog"

// With this, try to hook every print with zerolog
type SeverityHook struct{}

// Need to implement this
func (h SeverityHook) Run(e *zerolog.Event, level zerolog.Level, msg string) {
	if level != zerolog.NoLevel {
		e.Str("severity", level.String())
	}
}
