package handlers

import (
	"clingy-client/services"
	"log"
	"net/http"
)

type Server struct {
	server         *http.ServeMux
	configService  *services.Config
	contactService *services.Contact
	quicService    *services.Quic
}

func NewServer(
	configService *services.Config,
	contactService *services.Contact,
	quicService *services.Quic,
) *Server {
	return &Server{
		configService:  configService,
		contactService: contactService,
		quicService:    quicService,
	}
}

func (s *Server) Start(addr string) error {
	s.server = http.NewServeMux()

	log.Println("=== API Server Ready ===")
	log.Println("Endpoints:")

	// Contacts
	contactHandler := NewContactHandler(s.contactService)
	s.registerRoute("GET", "/api/contacts", contactHandler.GetContacts)
	s.registerRoute("POST", "/api/contacts", contactHandler.CreateContact)
	s.registerRoute("PUT", "/api/contacts", contactHandler.UpdateContact)
	s.registerRoute("DELETE", "/api/contacts", contactHandler.DeleteContact)

	// Config
	configHandler := NewConfigHandler(s.configService, s.quicService)
	s.registerRoute("GET", "/api/config/server", configHandler.GetServerConfig)
	s.registerRoute("POST", "/api/config/server", configHandler.SetServerConfig)

	// Health
	s.registerRoute("GET", "/health", s.handleHealth)

	log.Println("========================")
	log.Printf("Server: Starting HTTP server on %s", addr)

	return http.ListenAndServe(addr, s.server)
}

func (s *Server) registerRoute(method string, route string, handler http.HandlerFunc) {
	s.server.HandleFunc(method+" "+route, handler)
	// Pad method to 6 characters for alignment (longest verb is DELETE = 6 chars)
	paddedMethod := method
	for len(paddedMethod) < 6 {
		paddedMethod += " "
	}
	log.Printf("  %s  %s", paddedMethod, route)
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	log.Printf("API: GET /health - Health check requested")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"healthy"}`))
}
