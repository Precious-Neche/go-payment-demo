# Go Payment Demo

A complete Go web application with an HTML frontend for integrating the Paystack payment gateway.

## Features

- Paystack payment integration (initialize, verify, webhooks)
- Bootstrap 5 responsive frontend
- Real-time amount and product synchronization
- Mobile-friendly design
- Secure configuration management

## Project Structure

```text
paystack-go-app/
├── frontend/           # HTML, CSS, JavaScript
├── internal/
│   ├── config/         # Configuration loader
│   ├── paystack/       # Paystack API client
│   ├── handlers/       # HTTP request handlers
│   └── server/         # Server setup
├── main.go             # Application entry point
└── .env                # Environment variables

## Setup

-Get your Paystack API keys from https://dashboard.paystack.com

-Clone the repository and install dependencies

-Configure your .env file with your Paystack keys

-Run the application: go run main.go

Technologies

Go 1.22+

Chi Router

Bootstrap 5

Paystack API
