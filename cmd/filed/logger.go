package main

import (
	"errors"
	"log/slog"
	"os"
)

var ErrWrongLogFormat = errors.New("wrong log handler format")

func getLogger(logFormat string) (*slog.Logger, error) {
	var logger *slog.Logger
	switch logFormat {
	case "text":
		logger = slog.New(slog.NewTextHandler(os.Stdout, nil))
	case "json":
		logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))
	default:
		return nil, ErrWrongLogFormat
	}

	return logger, nil
}
