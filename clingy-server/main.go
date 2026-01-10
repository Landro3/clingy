package main

import (
	"crypto/rand"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

var connMap *ConnectionMap

func generateUUID() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		log.Printf("Error generating UUID: %v", err)
		return ""
	}
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:16])
}

func main() {
	cert, err := tls.LoadX509KeyPair("server.crt", "server.key")
	if err != nil {
		log.Fatal(err)
	}

	// tlsConfig := &tls.Config{
	// 	MinVersion:   tls.VersionTLS13,
	// 	Certificates: []tls.Certificate{cert},
	// 	NextProtos:   []string{"h3"},
	// }

	tlsConfig := &tls.Config{
		MinVersion:   tls.VersionTLS12,
		Certificates: []tls.Certificate{cert},
		NextProtos:   []string{"h2", "http/1.1"},
	}

	mux := http.NewServeMux()
	mux.HandleFunc("POST /register", register)
	mux.HandleFunc("POST /chat", chat)
	// server := http3.Server{
	// 	Handler:    mux,
	// 	Addr:       "0.0.0.0:8443",
	// 	TLSConfig:  http3.ConfigureTLSConfig(tlsConfig),
	// 	QUICConfig: &quic.Config{},
	// }
	server := &http.Server{
		Addr:      "0.0.0.0:8443",
		Handler:   mux,
		TLSConfig: tlsConfig,
	}

	log.Println("Creating Connection Map")
	connMap = NewConnectionMap()

	// err = server.ListenAndServe()
	// if err != nil {
	// 	log.Fatal(err)
	// 	return
	// }
	// log.Println("Server listening on :8443")

	log.Println("Starting HTTP/2 server on :8443")
	err = server.ListenAndServeTLS("", "") // Uses certificates from TLSConfig
	if err != nil {
		log.Fatal(err)
	}
}

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

func register(w http.ResponseWriter, r *http.Request) {
	assignedID := generateUUID()

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
	connMap.Add(body.Username, w)
	log.Printf("SSE connection established for user: %s", body.Username)

	// Keep connection alive for incoming messages
	<-r.Context().Done()

	connMap.Remove(body.Username)
	log.Printf("SSE connection closed for user: %s", body.Username)
}

func chat(w http.ResponseWriter, r *http.Request) {
	var body chatMessage
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	writer, exists := connMap.Get(body.To)
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

	connMap.LogConnections()
}
