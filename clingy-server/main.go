package main

import (
	"context"
	"crypto/rand"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"

	quic "github.com/quic-go/quic-go"
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

	tlsConfig := &tls.Config{
		MinVersion:   tls.VersionTLS13,
		Certificates: []tls.Certificate{cert},
		NextProtos:   []string{"clingy-v1"},
	}

	listener, err := quic.ListenAddr(":8443", tlsConfig, nil)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Server listening on :8443")

	log.Println("Creating Connection Map")
	connMap = NewConnectionMap()

	for {
		conn, err := listener.Accept(context.Background())
		if err != nil {
			log.Printf("Accept error: %v", err)
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn *quic.Conn) {
	defer conn.CloseWithError(0, "")

	for {
		stream, err := conn.AcceptStream(context.Background())
		if err != nil {
			return
		}
		go handleStream(stream, conn)
	}
}

type registerMessage struct {
	Username string
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

func handleStream(stream *quic.Stream, conn *quic.Conn) {
	buffer := make([]byte, 1024)
	n, err := stream.Read(buffer)
	if err != nil {
		if err.Error() == "EOF" {
			log.Printf("Client closed stream normally")
			return
		}
		log.Printf("Read error (%d bytes): %v", n, err)
		return
	}

	var regMsg registerMessage
	if err := json.Unmarshal(buffer[:n], &regMsg); err == nil && regMsg.Username != "" {
		assignedID := generateUUID()

		log.Printf("Registered\nUsername: %s\nAssigned ID: %s", regMsg.Username, assignedID)
		// TODO: base this off of UUID
		connMap.Add(regMsg.Username, conn)

		response := registrationResponse{
			Success:    true,
			AssignedID: assignedID,
			Username:   regMsg.Username,
			Message:    "Registration successful",
		}

		responseBytes, err := json.Marshal(response)
		if err != nil {
			log.Printf("Error marshaling registration response: %v", err)
			return
		}

		stream.Write(responseBytes)
	}

	var chatMsg chatMessage
	if err := json.Unmarshal(buffer[:n], &chatMsg); err == nil && chatMsg.To != "" {
		conn, exists := connMap.Get(chatMsg.To)
		if exists {
			stream, err := conn.OpenStreamSync(context.Background())
			if err != nil {
				log.Fatal(err)
				return
			}

			bytes, _ := json.Marshal(chatMsg)
			n, err := stream.Write(bytes)
			if err != nil {
				log.Printf("Failed to send message:\n%s", err)
			}

			log.Printf("✅ Sent %d bytes to %s: %s", n, chatMsg.To, chatMsg.Message)
		}
	}

	connMap.LogConnections()

	log.Printf("Message: %s", buffer[:n])
}
