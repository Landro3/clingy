package services

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"github.com/quic-go/quic-go"
	"github.com/quic-go/quic-go/http3"
	"context"
	"clingy-server/util"
)

type registerBody struct {
	Username string `json:"username"`
}

type registrationResponse struct {
	Success    bool   `json:"success"`
	AssignedID string `json:"assignedId"`
	Username   string `json:"username"`
	Message    string `json:"message,omitempty"`
}

type chatMessage struct {
	To      string
	From    string
	Message string
}

type Server interface {
	ListenAndServeTLS(certFile, keyFile string) error
	Shutdown(ctx context.Context) error
}

type ServerType int
const (
	Http2 ServerType = iota
	Http3
)

type ClingyServer struct {
	serverType ServerType
	mux        *http.ServeMux
	connMap    *ConnectionMap
	server     Server 
}


func NewClingyServer(serverType ServerType) Server {
  cs := &ClingyServer{
		serverType: serverType,
		mux:        http.NewServeMux(),
		connMap:    NewConnectionMap(),
	}

	cs.mux.HandleFunc("POST /register", cs.register)
	cs.mux.HandleFunc("POST /chat", cs.chat)

	tlsConfig := cs.getTlsConfig()

	if serverType == Http2 {
		server := &http.Server{
			Handler:   cs.mux,
			Addr:      "0.0.0.0:8443",
			TLSConfig: tlsConfig,
		}

		return server
	}

	server := &http3.Server{
		Handler:    cs.mux,
		Addr:       "0.0.0.0:8443",
		TLSConfig:  http3.ConfigureTLSConfig(tlsConfig),
		QUICConfig: &quic.Config{},
	}

	return server
}

func (cs *ClingyServer) register(w http.ResponseWriter, r *http.Request) {
	assignedID := util.GenerateUUID()

	var body registerBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	log.Printf("Registered\nUsername: %s\nAssigned ID: %s", body.Username, assignedID)

	// Set SSE headers for streaming response
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)

	// Send registration response as SSE event
	response := registrationResponse{
		Success:    true,
		AssignedID: assignedID,
		Username:   body.Username,
		Message:    "Registration successful",
	}

	responseBytes, err := json.Marshal(response)
	if err != nil {
		log.Printf("Error marshaling registration response: %v", err)
		return
	}

	fmt.Fprintf(w, "data: %s\n\n", string(responseBytes))
	if flusher, ok := w.(http.Flusher); ok {
		flusher.Flush()
	}

	// TODO: change to UUID
	cs.connMap.Add(body.Username, w)
	log.Printf("SSE connection established for user: %s", body.Username)

	// Keep connection alive for incoming messages
	<-r.Context().Done()

	cs.connMap.Remove(body.Username)
	log.Printf("SSE connection closed for user: %s", body.Username)
}

func (cs *ClingyServer) chat(w http.ResponseWriter, r *http.Request) {
	var body chatMessage
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	writer, exists := cs.connMap.Get(body.To)
	if exists {
		jsonBytes, _ := json.Marshal(body)
		_, err := fmt.Fprintf(writer, "data: %s\n\n", string(jsonBytes))
		if err != nil {
			log.Printf("STREAM: Error writing to response: %v", err)
			return
		}
		if flusher, ok := writer.(http.Flusher); ok {
			flusher.Flush()
		}
		log.Printf("✅ Sent message to %s", body.To)
	} else {
		log.Printf("User %s not connected", body.To)
	}

	w.Header().Set("Content-Type", "application/json")

	cs.connMap.LogConnections()
}

func (cs *ClingyServer) getTlsConfig() *tls.Config {
	cert, err := tls.LoadX509KeyPair("server.crt", "server.key")
	if err != nil {
		log.Fatal(err)
	}

	if (cs.serverType == Http2) {
		return &tls.Config{
			MinVersion:   tls.VersionTLS12,
			Certificates: []tls.Certificate{cert},
			NextProtos:   []string{"h2", "http/1.1"},
		}
	}

	return &tls.Config{
		MinVersion:   tls.VersionTLS13,
		Certificates: []tls.Certificate{cert},
		NextProtos:   []string{"h3"},
	}
}

