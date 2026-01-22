package paystack

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"io"
	"fmt"
	
)

func (c *Client) VerifyWebhookSignature(r *http.Request, body []byte) bool {
	signature := r.Header.Get("x-paystack-signature")
	if signature == "" {
		return false
	}

	mac := hmac.New(sha512.New, []byte(c.secretKey))
	mac.Write(body)
	expectedSignature := hex.EncodeToString(mac.Sum(nil))

	return hmac.Equal([]byte(signature), []byte(expectedSignature))
}

func (c *Client) ParseWebhook(r *http.Request) (*WebhookEvent, error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	if !c.VerifyWebhookSignature(r, body) {
		return nil, fmt.Errorf("invalid webhook signature")
	}

	var event WebhookEvent
	if err := json.Unmarshal(body, &event); err != nil {
		return nil, err
	}

	return &event, nil
}