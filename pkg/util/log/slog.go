package log

import (
	"context"
	"io"
	"log/slog"
)

// defaultNoopFunc returns the provided input
var defaultNoopFunc = func(input string) string {
	return input
}

// slogWriter provides implements [io.Writer] for [slog].
type slogWriter struct {
	ctx              context.Context
	logLevel         slog.Level
	logLineProcessor func(string) string
	defaultBaggage   []slog.Attr
}

// NewSlogWriter constructs a new instance for translating between [io.Writer] and [slog].
func NewSlogWriter(ctx context.Context, logLevel slog.Level, defaultBaggage map[string]string, processor func(string) string) io.Writer {
	var fn = defaultNoopFunc
	if processor != nil {
		fn = processor
	}

	return &slogWriter{
		ctx:              ctx,
		logLevel:         logLevel,
		logLineProcessor: fn,
		defaultBaggage:   convertToAttrs(defaultBaggage),
	}
}

func (s *slogWriter) Write(b []byte) (int, error) {
	n := len(b)
	if n > 0 && b[n-1] == '\n' {
		b = b[:n-1]
	}

	slog.Default().LogAttrs(s.ctx, s.logLevel, s.logLineProcessor(string(b)), s.defaultBaggage...)
	return n, nil
}

func convertToAttrs(baggage map[string]string) []slog.Attr {
	var converted = make([]slog.Attr, len(baggage))
	for k, v := range baggage {
		converted = append(converted, slog.String(k, v))
	}

	return converted
}
