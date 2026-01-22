package paystack

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"time"
)

type Client struct {
	secretKey string
	baseURL   string
	client    *http.Client
}

func NewClient(secretKey string) *Client {
	return &Client{
		secretKey: secretKey,
		baseURL:   "https://api.paystack.co",
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *Client) makeRequest(method, endpoint string, body interface{}) ([]byte, error) {
	var reqBody io.Reader
	
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		reqBody = bytes.NewBuffer(jsonData)
	}

	req, err := http.NewRequest(method, c.baseURL+endpoint, reqBody)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+c.secretKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}