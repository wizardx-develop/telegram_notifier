package forgejo

import (
	"fmt"
	"log/slog"
	"os"
	"wizardx/telegram_notifier/internal/config"

	"code.gitea.io/sdk/gitea"
)

func CreateClient(url, gitToken string, options []gitea.ClientOption) (*gitea.Client, error) {
	tokenOpt := gitea.SetToken(gitToken)
	options = append(options, tokenOpt)

	client, err := gitea.NewClient(url, options...)
	if err != nil {
		return nil, err
	}

	return client, err
}

func GetCommitMsg(client *gitea.Client, pr *ActionPayload) (*gitea.Commit, error) {
	commit, responce, err := client.GetSingleCommit(
		pr.PullRequest.User.Username,
		pr.PullRequest.Base.Repo.Name,
		pr.PullRequest.Head.SHA,
	)
	if err != nil {
		slog.Debug("create client error", "error", err)
		err = fmt.Errorf("error while getting commit: %w; status: %d", err, responce.StatusCode)
		return nil, err
	}
	if responce.StatusCode >= 400 {
		return nil, fmt.Errorf("request error %d", responce.StatusCode)
	}
	return commit, nil
}

func GetForgejoToken() string {
	token := os.Getenv(config.ForgejoToken.ToString())
	return token
}

func GetServerURL() string {
	url := os.Getenv(config.ForgejoServerURL.ToString())
	return url
}
