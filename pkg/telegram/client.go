package telegram

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Client interface {
	SendMessage(ctx context.Context, chatID string, text string) error
}

type client struct {
	token      string
	httpClient *http.Client
}

func NewClient(token string) Client {
	return &client{
		token: token,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

type sendMessageReq struct {
	ChatID string `json:"chat_id"`
	Text   string `json:"text"`
}

func (c *client) SendMessage(ctx context.Context, chatID string, text string) error {
	if c.token == "" {
		return fmt.Errorf("telegram bot token is empty")
	}

	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", c.token)

	reqBody := sendMessageReq{
		ChatID: chatID,
		Text:   text,
	}

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("failed to marshal telegram request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return fmt.Errorf("failed to create telegram request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send telegram request: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("telegram api returned status: %d", resp.StatusCode)
	}

	return nil
}
