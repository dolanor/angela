package main

import (
	"errors"
	"log/slog"
	"os"
)

var ErrWrongLogFormat = errors.New("wrong log handler format")

func getLogger(logFormat string) (*slog.Logger, error) {
	var logger *slog.Logger
	var opts *slog.HandlerOptions

	opts = &slog.HandlerOptions{Level: slog.LevelDebug}

	switch logFormat {
	case "text":
		logger = slog.New(slog.NewTextHandler(os.Stdout, opts))
	case "json":
		logger = slog.New(slog.NewJSONHandler(os.Stdout, opts))
	default:
		return nil, ErrWrongLogFormat
	}

	return logger, nil
}
