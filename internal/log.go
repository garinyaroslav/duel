package internal

import (
	"log/slog"
	"os"
)

func InitLogger() {
	debugEnvValue := os.Getenv("DEBUG")

	var logLevel slog.Level
	if debugEnvValue == "1" {
		logLevel = slog.LevelDebug
	} else {
		logLevel = slog.LevelInfo
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel}))

	slog.SetDefault(logger)
}
