package logger

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"syscall"
)

func Init(verbosity int) {
	initWithWriter(os.Stderr, verbosity)
}

func InitWithFile(verbosity int, path string) (*os.File, error) {
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, syscall.S_IRUSR|syscall.S_IWUSR)
	if err != nil {
		return nil, fmt.Errorf("error opening file for logging: %w", err)
	}

	initWithWriter(f, verbosity)

	return f, nil
}

//nolint:mnd // Verbosity is parsed from int to allow for a simple 1-5 scale.
func initWithWriter(w io.Writer, verbosity int) {
	opts := &slog.HandlerOptions{}

	switch verbosity {
	case 1:
		opts.Level = slog.LevelError
	case 2:
		opts.Level = slog.LevelWarn
	case 3:
		opts.Level = slog.LevelInfo
	case 4:
		opts.Level = slog.LevelDebug
	case 5:
		opts.Level = slog.LevelDebug
		opts.AddSource = true
	}

	handler := slog.NewTextHandler(w, opts)
	logger := slog.New(handler)
	slog.SetDefault(logger)
}
