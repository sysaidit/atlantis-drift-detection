package notification

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type SlackWebhook struct {
	WebhookURL string
	HTTPClient *http.Client
}

func (s *SlackWebhook) TemporaryError(ctx context.Context, dir string, workspace string, err error) error {
	return s.sendSlackMessage(ctx, fmt.Sprintf("Unknown error in remote\nDirectory: %s\nWorkspace: %s\nError: %s", dir, workspace, err.Error()))
}

func NewSlackWebhook(webhookURL string, HTTPClient *http.Client) *SlackWebhook {
	if webhookURL == "" {
		return nil
	}
	return &SlackWebhook{
		WebhookURL: webhookURL,
		HTTPClient: HTTPClient,
	}
}

type SlackWebhookMessage struct {
	Text string `json:"text"`
}

func (s *SlackWebhook) sendSlackMessage(ctx context.Context, msg string) error {
	body := SlackWebhookMessage{
		Text: msg,
	}
	b, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("failed to marshal slack webhook message: %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, s.WebhookURL, bytes.NewReader(b))
	if err != nil {
		return fmt.Errorf("failed to create slack webhook request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := s.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send slack webhook request: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to send slack webhook request: %w", err)
	}
	return nil
}

func formatRepoDirWorkspaceMsg(msgContext string, dir string, workspace string) string {
	gh_repo := os.Getenv("GITHUB_REPOSITORY")
	gh_repo_url := fmt.Sprintf("%s/%s", os.Getenv("GITHUB_SERVER_URL"), gh_repo)
	gh_action_id := os.Getenv("GITHUB_RUN_ID")
	gh_action_url := fmt.Sprintf("%s/actions/runs/%s", gh_repo_url, gh_action_id)

	msg := `%s in <%s|%s> detected by action <%s|%s>:
	Directory: %s Workspace: %s
	`
	slackMsg := fmt.Sprintf(msg, msgContext, gh_repo_url, gh_repo, gh_action_url, gh_action_id, dir, workspace)
	return slackMsg
}

func (s *SlackWebhook) ExtraWorkspaceInRemote(ctx context.Context, dir string, workspace string) error {
	//return s.sendSlackMessage(ctx, fmt.Sprintf("Extra workspace in remote\nDirectory: %s\nWorkspace: %s", dir, workspace))
	slackMsg := formatRepoDirWorkspaceMsg("Extra workspace", dir, workspace)
	return s.sendSlackMessage(ctx, slackMsg)
}

func (s *SlackWebhook) MissingWorkspaceInRemote(ctx context.Context, dir string, workspace string) error {
	//return s.sendSlackMessage(ctx, fmt.Sprintf("Missing workspace in remote\nDirectory: %s\nWorkspace: %s", dir, workspace))
	slackMsg := formatRepoDirWorkspaceMsg("Missing workspace", dir, workspace)
	return s.sendSlackMessage(ctx, slackMsg)
}

func (s *SlackWebhook) PlanDrift(ctx context.Context, dir string, workspace string) error {
	slackMsg := formatRepoDirWorkspaceMsg("Terraform Drift", dir, workspace)
	return s.sendSlackMessage(ctx, slackMsg)
}

var _ Notification = &SlackWebhook{}
