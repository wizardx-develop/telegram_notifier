package main

import (
	"fmt"
	"log/slog"
	"os"

	"wizardx/telegram_notifier/internal/services/forgejo"
	"wizardx/telegram_notifier/internal/services/telegram"
)

// LoggerInit настройка логирования
func LoggerInit() {
	opts := slog.HandlerOptions{
		Level:     slog.LevelInfo,
		AddSource: false,
	}
	logger := slog.New(slog.NewJSONHandler(os.Stderr, &opts))
	slog.SetDefault(logger)
}

func main() {
	LoggerInit()

	pr := &forgejo.PullRequestAction{}

	err := pr.Parse()
	if err != nil {
		os.Exit(1)
	}
	msg, err := telegram.CreateMessage(pr)
	if err != nil {
		os.Exit(1)
	}

	fmt.Println(msg.Message)
	err = telegram.SendMessage(msg)
	if err != nil {
		slog.Error("failed send telegram message", "error", err)
	}
}
