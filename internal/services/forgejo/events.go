package forgejo

import (
	"encoding/json"
	"log/slog"
	"os"
	"wizardx/telegram_notifier/internal/config"
)

const eventType = "pull_request"

// PullRequestAction структура для работы с данными эвента из файла
type PullRequestAction struct {
	Action      string      `json:"action"`
	Number      uint        `json:"number"`
	PullRequest pullRequest `json:"pull_request"`
}

type pullRequest struct {
	Base      prBase `json:"base"`
	State     string `json:"state"`
	Title     string `json:"title"`
	UpdatedAt string `json:"updated_at"`
	URL       string `json:"url"`
	User      prUser `json:"user"`
	Merged    bool   `json:"merged"`
}

type prUser struct {
	Username string `json:"username"`
	HTMLURL  string `json:"html_url"`
}

type prBase struct {
	Repo prRepo `json:"repo"`
}

type prRepo struct {
	FullName string `json:"full_name"`
	HTMLURL  string `json:"html_url"`
}

// ReadEventFile читает файл эвента из ОС
func ReadEventFile() ([]byte, error) {
	path, err := GetEventFilePath()
	if err != nil {
		return nil, err
	}
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return content, nil
}

// GetEventFilePath получает путь к файлу с описанием эвента из переменной окружения
func GetEventFilePath() (string, error) {
	path := os.Getenv(config.ForgejoEventPath.ToString())
	if path == "" {
		return "", config.ErrEventFilePathEmpty
	}
	return path, nil
}

// Parse парсит jsonчик для получения нужной структуры сообщения
func (pr *PullRequestAction) Parse() error {
	err := CheckEventType()
	if err != nil {
		slog.Error("check event type failed", slog.Any("error", err))
		return err
	}

	content, err := ReadEventFile()
	if err != nil {
		slog.Error("read event file failed", slog.Any("error", err))
		return err
	}

	err = json.Unmarshal(content, pr)
	if err != nil {
		slog.Error("unmarshal event file failed", slog.Any("error", err))
		return err
	}

	return nil
}

// CheckEventType проверяет тип эвента (pull_request)
func CheckEventType() error {
	env := os.Getenv(string(config.ForgejoEventName))
	if env != eventType {
		return config.ErrUnknownEvent
	} else if env == "" {
		return config.ErrEmptyEvent
	}
	return nil
}
