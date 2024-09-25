package notifier

import (
	"encoding/json"
	"fmt"
	"net/http"
	"project-wraith/src/pkg/req"
)

type TelegramBot struct {
	botToken string
	chatID   string
}

func NewTelegramBot(botToken string, chatID string) TelegramBot {
	return TelegramBot{
		botToken: botToken,
		chatID:   chatID,
	}
}

func (tb TelegramBot) SendChatNotification(text string) (string, error) {
	message := map[string]string{
		"chat_id": tb.chatID,
		"text":    text,
	}

	apiURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", tb.botToken)

	jsonData, err := json.Marshal(message)
	if err != nil {
		return "", fmt.Errorf("failed to marshal message: %w", err)
	}

	content := req.HTTPRequest{
		Method: http.MethodPost,
		URL:    apiURL,
		Body:   jsonData,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}
	res, err := req.SendRequest(content)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}

	return res, nil
}
