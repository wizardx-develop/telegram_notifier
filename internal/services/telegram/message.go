package telegram

import (
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"wizardx/telegram_notifier/internal/config"
	"wizardx/telegram_notifier/internal/services/forgejo"

	"code.gitea.io/sdk/gitea"
	"github.com/PaulSonOfLars/gotgbot/v2"
)

type FormatMessage struct {
	Message        string
	InlineKeyboard *gotgbot.InlineKeyboardMarkup
}

func choiceEmoji(prType string) string {
	switch prType {
	case config.OpenedType:
		return "🆕"
	case config.ClosedType:
		return "❌"
	case config.MergedType:
		return "🔀"
	case config.SynchronizedType:
		return "🔄"
	case config.ReopenedType:
		return "🔝"
	default:
		return ""
	}
}
func isMerged(pr *forgejo.ActionPayload) string {
	if pr.Action == config.ClosedType {
		if pr.PullRequest.Merged {
			return config.MergedType
		}
	}
	return pr.Action
}
func CreateMessage(pr *forgejo.ActionPayload, commit *gitea.Commit) (*FormatMessage, error) {
	action := isMerged(pr)
	emoji := choiceEmoji(action)

	msg := fmt.Sprintf(
		"%s <b>Pull Request №%d:</b> <code>%s</code>\n"+
			"📝 <b>PR Title:</b> <a href=\"%s\">%s</a>\n\n"+
			"🧑‍💻 <b>Actor:</b> <a href=\"%s\">%s</a>\n"+
			"📦 <b>Repository:</b> <a href=\"%s\">%s</a>\n",
		emoji, pr.Number, action,
		pr.PullRequest.URL, pr.PullRequest.Title,
		pr.PullRequest.User.HTMLURL, pr.PullRequest.User.Username,
		pr.PullRequest.Base.Repo.HTMLURL, pr.PullRequest.Base.Repo.FullName,
	)
	if commit != nil && action == config.SynchronizedType {
		msg += fmt.Sprintf("✉️ <b>Commit message:</b> <a href=\"%s\">%s</a>\n",
			commit.HTMLURL, commit.RepoCommit.Message,
		)
	}
	message := &FormatMessage{
		Message:        msg,
		InlineKeyboard: CreateButtonWithLink(pr),
	}
	return message, nil
}

func SendMessage(msg *FormatMessage) error {
	msgOpts := &gotgbot.SendMessageOpts{
		ParseMode:          gotgbot.ParseModeHTML,
		ReplyMarkup:        msg.InlineKeyboard,
		LinkPreviewOptions: &gotgbot.LinkPreviewOptions{IsDisabled: true},
	}

	botToken := os.Getenv(config.TelegramBotToken.ToString())
	topicID := os.Getenv(config.TelegramTopicID.ToString())

	if topicID != "" {
		topicIDint, err := strconv.ParseInt(topicID, 10, 64)
		if err != nil {
			return err
		}
		msgOpts.MessageThreadId = topicIDint
	}

	chatID, err := strconv.ParseInt(os.Getenv(config.TelegramChatID.ToString()), 10, 64)
	if err != nil {
		return err
	}

	b, err := gotgbot.NewBot(botToken, nil)
	if err != nil {
		return err
	}

	if topicID != "" {

	}

	// Вызов метода SendMessage, передача текста с HTML-тегами и созданных опций.
	_, err = b.SendMessage(chatID, msg.Message, msgOpts)
	if err != nil {
		return err
	}

	slog.Info("the message has been sent")
	return nil
}

func CreateButtonWithLink(pr *forgejo.ActionPayload) *gotgbot.InlineKeyboardMarkup {
	btn := gotgbot.InlineKeyboardButton{
		Text: "↗️Link to Pull Request",
		Url:  pr.PullRequest.URL,
	}
	inlineButtons := [][]gotgbot.InlineKeyboardButton{
		[]gotgbot.InlineKeyboardButton{btn},
	}
	kb := gotgbot.InlineKeyboardMarkup{
		InlineKeyboard: inlineButtons,
	}
	return &kb
}
