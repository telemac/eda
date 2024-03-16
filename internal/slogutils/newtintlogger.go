package slogutils

import (
	"github.com/lmittmann/tint"
	"log/slog"
	"os"
	"path"
)

func NewTintLogger(Level slog.Leveler, addSource, addProcess, setDefaultLogger bool) *slog.Logger {
	thitHandler := tint.NewHandler(os.Stderr, &tint.Options{
		AddSource:  addSource,
		Level:      Level,
		TimeFormat: "2006-01-02 15:04:05.000",
		//TimeFormat:  time.Kitchen,
		ReplaceAttr: CleanSourceAttr,
	})
	logger := slog.New(thitHandler)

	if addProcess {
		exe, err := os.Executable()
		if err != nil {
			logger.Error("os.Executable", "error", err)
		} else {
			exe = path.Base(exe)
			logger = logger.With(slog.String("process", exe))
		}
	}
	if setDefaultLogger {
		slog.SetDefault(logger)
	}
	return logger
}
