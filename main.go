package main

import (
	"log"
	"paystack-go-integration/config"
	"paystack-go-integration/internal/paystack"
	"paystack-go-integration/internal/server"
)

func main() {
	// Load configuration
	cfg := config.Load()
	
	// Validate configuration
	if cfg.PaystackSecretKey == "" {
		log.Fatal("PAYSTACK_SECRET_KEY is required. Set it in .env file")
	}
	
	// Initialize Paystack client
	paystackClient := paystack.NewClient(cfg.PaystackSecretKey)
	
	// Create and start server
	srv, err := server.New(cfg.ServerPort, paystackClient)
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}
	
	// Start the server
	if err := srv.Start(); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}