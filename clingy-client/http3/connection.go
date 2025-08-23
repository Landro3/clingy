package http3

import (
	"clingy-client/util"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"time"

	quic "github.com/quic-go/quic-go"
)

var connection *quic.Conn

type registerMessage struct {
	Username string
	ID       string
}

func ConnectToServer() {
	tlsConfig := &tls.Config{
		MinVersion:         tls.VersionTLS13,
		NextProtos:         []string{"clingy-v1"},
		InsecureSkipVerify: true, // Only for self-signed certs
	}
	conn, err := quic.DialAddr(context.Background(), "localhost:8443", tlsConfig, nil)
	if err != nil {
		util.Log(fmt.Sprintf("Failed to connect:\n%s", err))
	}

	connection = conn

	// Wait for handshake to complete
	select {
	case <-conn.HandshakeComplete():
		util.Log("✅ Handshake completed")
	case <-time.After(5 * time.Second):
		util.Log("❌ Handshake timeout")
	}

	registerInServer()
}

func registerInServer() {
	message := registerMessage{Username: "Testy", ID: "jfp9q8uf89ahflkj3r32o98g89"}

	bytes, err := json.Marshal(message)
	if err != nil {
		util.Log(fmt.Sprintf("%s", err))
	}

	payload := string(bytes)
	SendMessage(payload)
}

func SendMessage(msg string) {
	stream, err := connection.OpenStreamSync(context.Background())
	if err != nil {
		util.Log(fmt.Sprintf("Failed to open stream:\n%s", err))
	}
	defer stream.Close()

	n, err := stream.Write([]byte(msg))
	if err != nil {
		util.Log(fmt.Sprintf("Failed to send message:\n%s", err))
	}

	util.Log(fmt.Sprintf("✅ Sent %d bytes: %s", n, msg))

	// Optional: Read response from server
	buffer := make([]byte, 1024)
	n, err = stream.Read(buffer)
	if err != nil {
		util.Log(fmt.Sprintf("No response or error reading: %v", err))
	} else {
		util.Log(fmt.Sprintf("📨 Server response: %s", buffer[:n]))
	}
}
