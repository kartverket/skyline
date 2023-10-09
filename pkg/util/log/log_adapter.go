package log

import (
	"context"
	"fmt"
	"log/slog"
)

// CommonLogger represents common logging methods found across multiple libraries and stdlib.
type CommonLogger interface {
	Printf(format string, v ...interface{})
	Println(v ...interface{})
}

type logToSlogAdapter struct {
	ctx            context.Context
	logLevel       slog.Level
	defaultBaggage []slog.Attr
}

func NewLogAdapter(ctx context.Context, logLevel slog.Level, defaultBaggage map[string]string) CommonLogger {
	return &logToSlogAdapter{
		ctx:            ctx,
		logLevel:       logLevel,
		defaultBaggage: convertToAttrs(defaultBaggage),
	}
}

func (l *logToSlogAdapter) Printf(format string, v ...interface{}) {
	slog.Default().LogAttrs(l.ctx, l.logLevel, fmt.Sprintf(format, v...), l.defaultBaggage...)
}

func (l *logToSlogAdapter) Println(v ...interface{}) {
	for _, msg := range v {
		if msg != nil {
			slog.Default().LogAttrs(l.ctx, l.logLevel, fmt.Sprint(msg), l.defaultBaggage...)
		}
	}
}
