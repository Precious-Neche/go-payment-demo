package server

import (
	
	"log"
	"net/http"
	"paystack-go-integration/internal/handlers"
	"paystack-go-integration/internal/paystack"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Server struct {
	port           string
	paystackClient *paystack.Client
	paymentHandler *handlers.PaymentHandler
}

func New(port string, paystackClient *paystack.Client) (*Server, error) {
	// Create payment handler
	paymentHandler := handlers.NewPaymentHandler(paystackClient)
	
	return &Server{
		port:           port,
		paystackClient: paystackClient,
		paymentHandler: paymentHandler,
	}, nil
}

func (s *Server) Start() error {
	r := chi.NewRouter()
	
	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	
	// Setup routes
	s.setupRoutes(r)
	
	// Start server
	log.Printf("🚀 Server starting on http://localhost:%s", s.port)
	return http.ListenAndServe(":"+s.port, r)
}

func (s *Server) setupRoutes(r chi.Router) {
	// Serve static files (CSS, JS)
	fs := http.FileServer(http.Dir("frontend/static"))
	r.Handle("/static/*", http.StripPrefix("/static/", fs))
	
	// Serve HTML pages
	r.Get("/", s.serveHomePage)
	r.Get("/payment", s.servePaymentPage)
	r.Get("/success", s.serveSuccessPage)
	r.Get("/failed", s.serveFailedPage)
	
	// API routes - Mount your existing handler
	r.Mount("/api", s.paymentHandler.SetupRoutes())
}

// HTML page handlers
func (s *Server) serveHomePage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "frontend/templates/index.html")
}

func (s *Server) servePaymentPage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "frontend/templates/payment.html")
}

func (s *Server) serveSuccessPage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "frontend/templates/success.html")
}

func (s *Server) serveFailedPage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "frontend/templates/failed.html")
}