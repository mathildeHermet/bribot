package log

import (
	"context"
	"log/slog"
	"os"
	"runtime"
	"sync/atomic"
)

var defaultLogger atomic.Pointer[Logger]

type ctxKey string

const (
	ContextKeyTraceID ctxKey = "trace_id"
)

type contextWrapper struct {
	slog.Handler
}

func (w contextWrapper) Handle(ctx context.Context, r slog.Record) error {
	if attr := ctx.Value(ContextKeyTraceID); attr != nil {
		r.AddAttrs(slog.Attr{Key: string(ContextKeyTraceID), Value: slog.AnyValue(attr)})
	}
	return w.Handler.Handle(ctx, r)
}

type Logger interface {
	Debug(string, ...any)
	DebugContext(context.Context, string, ...any)
	Info(string, ...any)
	InfoContext(context.Context, string, ...any)
	Warn(string, ...any)
	WarnContext(context.Context, string, ...any)
	Error(string, error, ...any)
	ErrorContext(context.Context, string, error, ...any)
}

type logger struct {
	sl *slog.Logger
}

func NewLogger(opts ...Option) Logger {
	loggerOpts := Parse(opts...)
	logHandlerOpts := slog.HandlerOptions{
		Level: loggerOpts.Level(), // if unspecified fallback to info
	}
	writer := loggerOpts.Writer()
	if writer == nil {
		writer = os.Stdout
	}
	attrs := make([]slog.Attr, 0, 2)
	if len(loggerOpts.ServiceName()) > 0 {
		attrs = append(attrs, slog.Attr{Key: "service", Value: slog.AnyValue(loggerOpts.ServiceName())})
	}
	if len(loggerOpts.AppName()) > 0 {
		attrs = append(attrs, slog.Attr{Key: "app", Value: slog.AnyValue(loggerOpts.AppName())})
	}
	handler := &contextWrapper{slog.NewJSONHandler(writer, &logHandlerOpts).WithAttrs(attrs)}
	return &logger{sl: slog.New(handler)}
}

func (l *logger) Debug(msg string, args ...any) {
	l.sl.Debug(msg, addSourceToArgs(args...)...)
}

func (l *logger) DebugContext(ctx context.Context, msg string, args ...any) {
	l.sl.DebugContext(ctx, msg, addSourceToArgs(args...)...)
}

func (l *logger) Info(msg string, args ...any) {
	l.sl.Info(msg, addSourceToArgs(args...)...)
}

func (l *logger) InfoContext(ctx context.Context, msg string, args ...any) {
	l.sl.InfoContext(ctx, msg, addSourceToArgs(args...)...)
}

func (l *logger) Warn(msg string, args ...any) {
	l.sl.Warn(msg, addSourceToArgs(args...)...)
}

func (l *logger) WarnContext(ctx context.Context, msg string, args ...any) {
	l.sl.WarnContext(ctx, msg, addSourceToArgs(args...)...)
}

func (l *logger) Error(msg string, err error, args ...any) {
	l.sl.Error(msg, prependErrorToArgs(err, addSourceToArgs(args...)...)...)
}

func (l *logger) ErrorContext(ctx context.Context, msg string, err error, args ...any) {
	l.sl.ErrorContext(ctx, msg, prependErrorToArgs(err, addSourceToArgs(args...)...)...)
}

func prependErrorToArgs(err error, args ...any) []any {
	return append([]any{"error", err.Error()}, args...)
}

func SetDefaultLogger(opts ...Option) {
	l := NewLogger(opts...)
	defaultLogger.Store(&l)
}

func Debug(msg string, args ...any) {
	if l := defaultLogger.Load(); l != nil {
		(*l).(*logger).sl.Log(context.TODO(), slog.LevelDebug, msg, addSourceToArgs(args...)...)
	}
}

func DebugContext(ctx context.Context, msg string, args ...any) {
	if l := defaultLogger.Load(); l != nil {
		(*l).(*logger).sl.Log(ctx, slog.LevelDebug, msg, addSourceToArgs(args...)...)
	}
}

func Info(msg string, args ...any) {
	if l := defaultLogger.Load(); l != nil {
		(*l).(*logger).sl.Log(context.TODO(), slog.LevelInfo, msg, addSourceToArgs(args...)...)
	}
}

func InfoContext(ctx context.Context, msg string, args ...any) {
	if l := defaultLogger.Load(); l != nil {
		(*l).(*logger).sl.Log(ctx, slog.LevelInfo, msg, addSourceToArgs(args...)...)
	}
}

func Warn(msg string, args ...any) {
	if l := defaultLogger.Load(); l != nil {
		(*l).(*logger).sl.Log(context.TODO(), slog.LevelWarn, msg, addSourceToArgs(args...)...)
	}
}

func WarnContext(ctx context.Context, msg string, args ...any) {
	if l := defaultLogger.Load(); l != nil {
		(*l).(*logger).sl.Log(ctx, slog.LevelWarn, msg, addSourceToArgs(args...)...)
	}
}

func Error(msg string, err error, args ...any) {
	if l := defaultLogger.Load(); l != nil {
		(*l).(*logger).sl.Log(context.TODO(), slog.LevelError, msg, addSourceToArgs(prependErrorToArgs(err, args...)...)...)
	}
}

func ErrorContext(ctx context.Context, msg string, err error, args ...any) {
	if l := defaultLogger.Load(); l != nil {
		(*l).(*logger).sl.Log(ctx, slog.LevelError, msg, addSourceToArgs(prependErrorToArgs(err, args...)...)...)
	}
}

// https://github.com/golang/go/blob/master/src/log/slog/logger.go#L249
func addSourceToArgs(args ...any) []any {
	_, f, l, _ := runtime.Caller(2)
	group := slog.Group(
		"source",
		slog.Attr{
			Key:   "filename",
			Value: slog.AnyValue(f),
		},
		slog.Attr{
			Key:   "line",
			Value: slog.AnyValue(l),
		},
	)
	return append(args, group)
}
