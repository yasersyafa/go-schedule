package notifier

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type TelegramNotifier struct {
	botToken string
	chatID string
	client *http.Client
}

func NewTelegramNotifier(botToken, chatID string) *TelegramNotifier {
	return &TelegramNotifier{
		botToken: botToken, 
		chatID: chatID,
		client: &http.Client{},
	}
}

func (n *TelegramNotifier) Send(ctx context.Context, message string) error {
	apiURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", n.botToken)

	payload := map[string]string {
		"chat_id": n.chatID,
		"message": message,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal telegram payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, apiURL, nil)
	if err != nil {
		return fmt.Errorf("create telegram request: %w", err)
	}

	req.Body = io.NopCloser(bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	res, err := n.client.Do(req)
	if err != nil {
		return fmt.Errorf("send telegram request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("telegram api return status code: %d", res.StatusCode)
	}

	return nil
}