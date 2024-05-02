package log

import (
	"io"
	"log/slog"
)

type Option interface {
	apply(*Options)
}

type optionFunc func(*Options)

func (f optionFunc) apply(o *Options) {
	f(o)
}

func WithAppName(app string) Option {
	return optionFunc(func(o *Options) {
		o.appName = app
	})
}

func WithLevelDebug() Option {
	return optionFunc(func(o *Options) {
		o.level = slog.LevelDebug
	})
}

func WithLevelInfo() Option {
	return optionFunc(func(o *Options) {
		o.level = slog.LevelInfo
	})
}

func WithLevelWarn() Option {
	return optionFunc(func(o *Options) {
		o.level = slog.LevelWarn
	})
}

func WithLevelError() Option {
	return optionFunc(func(o *Options) {
		o.level = slog.LevelError
	})
}

func WithServiceName(service string) Option {
	return optionFunc(func(o *Options) {
		o.serviceName = service
	})
}

func WithWriter(writer io.Writer) Option {
	return optionFunc(func(o *Options) {
		o.writer = writer
	})
}

type Options struct {
	appName     string
	serviceName string
	level       slog.Leveler
	writer      io.Writer
}

func (opts *Options) Apply(opt Option) {
	opt.apply(opts)
}

func (opts *Options) AppName() string {
	return opts.appName
}

func (opts *Options) ServiceName() string {
	return opts.serviceName
}

func (opts *Options) Level() slog.Leveler {
	return opts.level
}

func (opts *Options) Writer() io.Writer {
	return opts.writer
}

func Parse(opts ...Option) *Options {
	o := &Options{}
	for _, option := range opts {
		option.apply(o)
	}
	return o
}
