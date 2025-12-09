package config

import "errors"

type env string

var (
	ErrEventFilePathEmpty = errors.New("error get event file path: path is empty")
	ErrUnknownEvent       = errors.New("event is unknown")
	ErrEmptyEvent         = errors.New("event is empty")
)

const (
	ForgejoToken     env    = "INPUT_FORGEJO_TOKEN"
	ForgejoServerURL env    = "FORGEJO_SERVER_URL"
	ForgejoEventName env    = "FORGEJO_EVENT_NAME"
	ForgejoEventPath env    = "FORGEJO_EVENT_PATH"
	TelegramBotToken env    = "INPUT_TOKEN"
	TelegramChatID   env    = "INPUT_CHAT_ID"
	TelegramTopicID  env    = "INPUT_THREAD_ID"
	ForgejoEnvFile   string = ".env"
	FromFile         bool   = true
	FromEnv          bool   = false
	ActionEvent      string = "pull_request"

	OpenedType       string = "opened"
	SynchronizedType string = "synchronized"
	MergedType       string = "merged"
	ClosedType       string = "closed"
	ReopenedType     string = "reopened"
)

func (e env) ToString() string {
	return string(e)
}
