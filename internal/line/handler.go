package line

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	channelSecret string
}

func NewHandler(channelSecret string) *Handler {
	return &Handler{channelSecret: channelSecret}
}

type webhookPayload struct {
	Events []struct {
		Type   string `json:"type"`
		Source struct {
			UserID string `json:"userId"`
		} `json:"source"`
	} `json:"events"`
}

func (h *Handler) Webhook(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)

	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	if !h.validSignature(body, c.GetHeader("X-Line-Signature")) {
		c.Status(http.StatusUnauthorized)
	}

	var payload webhookPayload
	if err := json.Unmarshal(body, &payload); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	for _, event := range payload.Events {
		if event.Type == "follow" {
			log.Printf("LINE follow event - User ID kamu: %s", event.Source.UserID)
		}
	}

	c.Status(http.StatusOK)
}

func (h *Handler) validSignature(body []byte, signature string) bool {
	mac := hmac.New(sha256.New, []byte(h.channelSecret))
	mac.Write(body)
	expected := base64.StdEncoding.EncodeToString(mac.Sum(nil))
	return hmac.Equal([]byte(expected), []byte(signature))
}
