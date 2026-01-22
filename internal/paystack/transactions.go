package paystack

import (
	"encoding/json"
	"fmt"
)


func (c *Client) InitializeTransaction(req TransactionRequest) (*TransactionResponse, error) {
	if req.Amount < 10000 {
		return nil, fmt.Errorf("amount must be at least ₦100 (10000 kobo)")
	}

	body, err := c.makeRequest("POST", "/transaction/initialize", req)
	if err != nil {
		return nil, err
	}

	var resp TransactionResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}

	if !resp.Status {
		return nil, fmt.Errorf("paystack error: %s", resp.Message)
	}

	return &resp, nil
}

func (c *Client) VerifyTransaction(reference string) (*Transaction, error) {
	body, err := c.makeRequest("GET", "/transaction/verify/"+reference, nil)
	if err != nil {
		return nil, err
	}

	var resp struct {
		Status  bool        `json:"status"`
		Message string      `json:"message"`
		Data    Transaction `json:"data"`
	}

	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}

	if !resp.Status {
		return nil, fmt.Errorf("verification failed: %s", resp.Message)
	}

	return &resp.Data, nil
}

func (c *Client) ListTransactions(perPage, page int) ([]Transaction, error) {
	endpoint := fmt.Sprintf("/transaction?perPage=%d&page=%d", perPage, page)
	
	body, err := c.makeRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	var resp struct {
		Status  bool          `json:"status"`
		Message string        `json:"message"`
		Data    []Transaction `json:"data"`
		Meta    struct {
			Total     int `json:"total"`
			PerPage   int `json:"perPage"`
			Page      int `json:"page"`
			PageCount int `json:"pageCount"`
		} `json:"meta"`
	}

	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}

	if !resp.Status {
		return nil, fmt.Errorf("failed to list transactions: %s", resp.Message)
	}

	return resp.Data, nil
}