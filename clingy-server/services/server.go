package services

import (
	"clingy-server/util"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/quic-go/quic-go"
	"github.com/quic-go/quic-go/http3"
	"log"
	"net/http"
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

type ServerType int

const (
	Http2 ServerType = iota
	Http3
)

type ClingyServer struct {
	serverType ServerType
	mux        *http.ServeMux
	connMap    *ConnectionMap
	certFile	 string
	keyFile    string
}

func StartClingyServer(serverType ServerType) {
	cs := &ClingyServer{
		serverType: serverType,
		mux: http.NewServeMux(),
		connMap: NewConnectionMap(),
		certFile: "services/server.crt",
		keyFile: "services/server.key",
	}

	cs.mux.HandleFunc("POST /register", cs.register)
	cs.mux.HandleFunc("POST /chat", cs.chat)

	if serverType == Http2 {
		cs.startHttp2Server()
	} else {
		cs.startHttp3Server()
	}
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

	if err := cs.flushSSEMessage(w, response); err != nil {
		log.Printf("Error sending registration response: %v", err)
		return
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
		if err := cs.flushSSEMessage(writer, body); err != nil {
			log.Printf("STREAM: Error sending chat message: %v", err)
			return
		}
		log.Printf("✅ Sent message to %s", body.To)
	} else {
		log.Printf("User %s not connected", body.To)
	}

	cs.connMap.LogConnections()
}

func (cs *ClingyServer) flushSSEMessage(w http.ResponseWriter, data interface{}) error {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("error marshaling data: %v", err)
	}

	_, err = fmt.Fprintf(w, "data: %s\n\n", string(jsonBytes))
	if err != nil {
		return fmt.Errorf("error writing to response: %v", err)
	}

	if flusher, ok := w.(http.Flusher); ok {
		flusher.Flush()
	}

	return nil
}

func (cs *ClingyServer) startHttp2Server() {
	cert, err := tls.LoadX509KeyPair("services/server.crt", "services/server.key")
	if err != nil {
		log.Fatal("Error loading HTTP2 cert")
		log.Fatal(err)
	}

	tlsConfig := &tls.Config{
		MinVersion:   tls.VersionTLS12,
		Certificates: []tls.Certificate{cert},
		NextProtos:   []string{"h2", "http/1.1"},
	}

	server := &http.Server{
		Handler: cs.mux,
		Addr: "localhost:8443",
		TLSConfig: tlsConfig,
	}

	log.Printf("HTTP2 server listening at %s", server.Addr)
	server.ListenAndServe()
}

func (cs*ClingyServer) startHttp3Server() {
	cert, err := tls.LoadX509KeyPair(cs.certFile, cs.keyFile)
	if err != nil {
		log.Fatal("Error loading HTTP3 cert")
		log.Fatal(err)
		return
	}

	tlsConfig := &tls.Config{
		MinVersion:   tls.VersionTLS13,
		Certificates: []tls.Certificate{cert},
		NextProtos:   []string{"h3"},
	}

	server := &http3.Server{
		Handler:    cs.mux,
		Addr:       "localhost:8443",
		TLSConfig:  http3.ConfigureTLSConfig(tlsConfig),
		QUICConfig: &quic.Config{},
	}

	log.Printf("HTTP3 server listening at %s", server.Addr)
	server.ListenAndServe()
}
