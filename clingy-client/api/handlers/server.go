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
	http3Service   *services.Http3
	chatChannel    <-chan services.ChatMessage
}

func NewServer(
	configService *services.Config,
	contactService *services.Contact,
	http3Service *services.Http3,
	chatChannel <-chan services.ChatMessage,
) *Server {
	return &Server{
		configService:  configService,
		contactService: contactService,
		http3Service:   http3Service,
		chatChannel:    chatChannel,
	}
}

func (s *Server) Start(addr string) error {
	s.server = http.NewServeMux()

	log.Println("=== API Server Ready ===")
	log.Println("Endpoints:")

	// Contacts
	contactHandler := NewContactHandler(s.contactService, s.configService)
	s.registerRoute("GET", "/api/contacts", contactHandler.GetContacts)
	s.registerRoute("POST", "/api/contacts", contactHandler.CreateContact)
	s.registerRoute("PUT", "/api/contacts", contactHandler.UpdateContact)
	s.registerRoute("DELETE", "/api/contacts", contactHandler.DeleteContact)

	// Config
	configHandler := NewConfigHandler(s.configService, s.http3Service)
	s.registerRoute("GET", "/api/config/server", configHandler.GetServerConfig)
	s.registerRoute("POST", "/api/config/server", configHandler.SetServerConfig)
	// TODO: move to new handler
	s.registerRoute("POST", "/api/register", configHandler.RegisterWithServer)

	// Chat
	chatHandler := NewChatHandler(s.configService, s.http3Service, s.chatChannel)
	s.registerRoute("POST", "/api/chat", chatHandler.SendChatMessage)
	s.registerRoute("GET", "/api/chat/stream", chatHandler.GetMessageStream)

	// Health
	s.registerRoute("GET", "/health", s.handleHealth)

	log.Println("========================")
	log.Printf("Server: Starting HTTP server on %s", addr)

	return http.ListenAndServe(addr, s.server)
}

func (s *Server) registerRoute(method string, route string, handler http.HandlerFunc) {
	s.server.HandleFunc(method+" "+route, handler)
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
	_, err := w.Write([]byte(`{"status":"healthy"}`))
	if err != nil {
		log.Printf("Error writing response: %v", err)
	}
}
