package main

import (
	"fmt"
	"log/slog"
	"os"
	"strings"

	"wizardx/telegram_notifier/internal/config"
	"wizardx/telegram_notifier/internal/services/forgejo"
	"wizardx/telegram_notifier/internal/services/telegram"

	"code.gitea.io/sdk/gitea"
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

	pr := &forgejo.ActionPayload{}

	err := pr.Parse()
	if err != nil {
		os.Exit(1)
	}

	separator := strings.Repeat("*", 40)

	fmt.Println(separator)
	client, err := forgejo.CreateClient(
		os.Getenv(config.ForgejoServerURL.ToString()),
		os.Getenv(config.ForgejoToken.ToString()),
		[]gitea.ClientOption{},
	)
	if err != nil {
		slog.Error("create client error", "error", err)
		os.Exit(1)
	}
	commitMsg, err := forgejo.GetCommitMsg(client, pr)
	if err != nil {
		slog.Error("get commit msg error", "error", err)
		os.Exit(1)
	}
	msg, err := telegram.CreateMessage(pr, commitMsg)
	if err != nil {
		slog.Error("create message failed", "error", err)
		os.Exit(1)
	}

	err = telegram.SendMessage(msg)
	if err != nil {
		slog.Error("send telegram message failed", "error", err)
		os.Exit(1)
	}
}
