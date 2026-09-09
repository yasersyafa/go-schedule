package notifier

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type LineNotifier struct {
	channelToken string
	userID string
	client *http.Client
}

func NewLineNotifier(channelToken, userID string) *LineNotifier {
	return &LineNotifier{
		channelToken: channelToken,
		userID: userID,
		client: &http.Client{},
	}
}

func (n *LineNotifier) Send(ctx context.Context, message string) error {
	payload, err := json.Marshal(map[string]interface{}{
		"to": n.userID,
		"messages": []map[string]string{
			{"type": "text", "text": message},
		},
	})

	if err != nil {
		return fmt.Errorf("marshal line payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://api.line.me/v2/bot/message/push", bytes.NewReader(payload))
	if err != nil {
		return fmt.Errorf("create line request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer " + n.channelToken)

	res, err := n.client.Do(req)
	if err != nil {
		return fmt.Errorf("send line request: %w", err)
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("line api returned status code: %d", res.StatusCode)
	}

	return nil
}