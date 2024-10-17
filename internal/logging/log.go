package logging

import (
	"io"
	"log/slog"

	"github.com/lmittmann/tint"
)

type Options struct {
	Debug  bool
	Writer io.Writer
}

func NewLogger(opts Options) *slog.Logger {

	logLevel := slog.LevelInfo
	if opts.Debug {
		logLevel = slog.LevelDebug
	}

	return slog.New(tint.NewHandler(opts.Writer, &tint.Options{Level: logLevel}))
}
