package main

import (
	"log/slog"
	"os"

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

	serverURL := forgejo.GetServerURL()
	if serverURL == "" {
		slog.Error("Server URL is empty")
		os.Exit(1)
	}
	forgejoToken := forgejo.GetForgejoToken()
	if forgejoToken == "" {
		slog.Error("Forgejo token is empty")
		os.Exit(1)
	}

	client, err := forgejo.CreateClient(
		serverURL,
		forgejoToken,
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
