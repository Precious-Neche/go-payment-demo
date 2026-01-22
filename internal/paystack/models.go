package paystack

import "time"

type TransactionRequest struct {
	Email     string                 `json:"email"`
	Amount    int64                  `json:"amount"`
	Currency  string                 `json:"currency,omitempty"`
	Reference string                 `json:"reference,omitempty"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
	Channels  []string               `json:"channels,omitempty"`
}

type TransactionResponse struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
	Data    struct {
		AuthorizationURL string `json:"authorization_url"`
		AccessCode       string `json:"access_code"`
		Reference        string `json:"reference"`
	} `json:"data"`
}

type Transaction struct {
	ID        int64                  `json:"id"`
	Reference string                 `json:"reference"`
	Amount    int64                  `json:"amount"`
	Currency  string                 `json:"currency"`
	Status    string                 `json:"status"`
	CreatedAt time.Time              `json:"created_at"`
	PaidAt    *time.Time             `json:"paid_at"`
	Metadata  map[string]interface{} `json:"metadata"`
	Customer  struct {
		Email string `json:"email"`
	} `json:"customer"`
}

type WebhookEvent struct {
	Event string `json:"event"`
	Data  struct {
		Reference string                 `json:"reference"`
		Amount    int64                  `json:"amount"`
		Status    string                 `json:"status"`
		Metadata  map[string]interface{} `json:"metadata"`
	} `json:"data"`
}