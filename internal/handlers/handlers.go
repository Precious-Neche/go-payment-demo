package handlers

import (
	"encoding/json"
	"fmt"  
	"net/http"
	"paystack-go-integration/internal/paystack"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"  
)

type PaymentHandler struct {
	PaystackClient *paystack.Client
}

func NewPaymentHandler(client *paystack.Client) *PaymentHandler {
	return &PaymentHandler{
		PaystackClient: client,
	}
}

func (h *PaymentHandler) InitializePayment(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Amount   int64                  `json:"amount"`
		Email    string                 `json:"email"`
		Metadata map[string]interface{} `json:"metadata"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	paystackreq := paystack.TransactionRequest{
		Amount:   req.Amount * 100,
		Email:    req.Email,
		Metadata: req.Metadata,
		Currency: "NGN",
	}

	
	resp, err := h.PaystackClient.InitializeTransaction(paystackreq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h *PaymentHandler) VerifyPayment(w http.ResponseWriter, r *http.Request) {
	reference := chi.URLParam(r, "reference")
	resp, err := h.PaystackClient.VerifyTransaction(reference)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// FIX: Method name should match paystack package
// If paystack has ParseWebhook, use that instead
func (h *PaymentHandler) HandleWebHook(w http.ResponseWriter, r *http.Request) {
	// CHANGE HandleWebHook to ParseWebhook if that's what paystack has
	event, err := h.PaystackClient.ParseWebhook(r)  // CHANGED HERE
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	switch event.Event {
	case "charge.success":
		fmt.Printf("Payment successful for reference: %s\n", event.Data.Reference)

	case "transfer.success":
		// Handle transfer success

	case "charge.failed":
		// Handle failed payment

	default:
		fmt.Printf("Unhandled event type: %s\n", event.Event)
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Webhook Received"))
}

// FIX: Either implement ListTransactions or remove this method
func (h *PaymentHandler) ListTransactions(w http.ResponseWriter, r *http.Request) {
	// Temporarily return empty or remove this method
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "ListTransactions not implemented yet",
		"data":    []interface{}{},
	})
}

func (h *PaymentHandler) SetupRoutes() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Post("/initialize", h.InitializePayment)
	r.Get("/verify/{reference}", h.VerifyPayment)
	r.Post("/webhook", h.HandleWebHook)
	r.Get("/transactions", h.ListTransactions)

	return r
}