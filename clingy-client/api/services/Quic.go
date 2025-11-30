package services

import (
	"clingy-client/util"
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	quic "github.com/quic-go/quic-go"
)

type registerMessage struct {
	Username string `json:"username"`
}

type registrationResponse struct {
	Success    bool   `json:"success"`
	AssignedID string `json:"assignedId"`
	Username   string `json:"username"`
	Message    string `json:"message,omitempty"`
}

type ChatMessage struct {
	To      string
	From    string
	Message string
}

var messageChannel = make(chan ChatMessage, 100)

type Quic struct {
	conn *quic.Conn
	mu   sync.Mutex
}

func NewQuic(config *Config) *Quic {
	c := &Quic{}
	c.Connect(config.ServerAddr)
	c.StartMessageListener()
	return c
}

func (c *Quic) Connect(serverAddr string) error {
	tlsConfig := &tls.Config{
		MinVersion:         tls.VersionTLS13,
		NextProtos:         []string{"clingy-v1"},
		InsecureSkipVerify: true, // Only for self-signed certs
	}

	conn, err := quic.DialAddr(
		context.Background(),
		serverAddr,
		tlsConfig,
		&quic.Config{
			KeepAlivePeriod: 30 * time.Second,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to connect:\n%s", err)
	}

	select {
	case <-conn.HandshakeComplete():
		util.Log("✅ Handshake completed")
	case <-time.After(5 * time.Second):
		util.Log("❌ Handshake timeout")
		return errors.New("could not connect to server")
	}

	c.conn = conn
	return nil
}

func (c *Quic) Register(username string) (string, error) {
	util.Log(fmt.Sprintf("🔄 Starting registration process for user: %s", username))

	message := registerMessage{Username: username}
	bytes, err := json.Marshal(message)
	if err != nil {
		return "", fmt.Errorf("failed to marshal registration message: %s", err)
	}
	util.Log(fmt.Sprintf("📤 Sending registration message: %s", string(bytes)))

	stream, err := c.conn.OpenStreamSync(context.Background())
	if err != nil {
		return "", fmt.Errorf("failed to open registration stream: %s", err)
	}

	bytesWritten, err := stream.Write(bytes)
	if err != nil {
		util.Log(fmt.Sprintf("❌ Failed to write to stream: %v", err))
		return "", fmt.Errorf("failed to send registration message: %s", err)
	}
	util.Log(fmt.Sprintf("✅ Sent %d bytes to server", bytesWritten))

	timeout := 5 * time.Second
	deadline := time.Now().Add(timeout)
	stream.SetReadDeadline(deadline)
	util.Log(fmt.Sprintf("⏱️ Set read deadline: %v (timeout: %v)", deadline, timeout))

	buffer := make([]byte, 1024)
	util.Log("📥 Waiting for server response...")
	n, err := stream.Read(buffer)
	if err != nil {
		util.Log(fmt.Sprintf("❌ Failed to read from stream: %v", err))
		return "", fmt.Errorf("failed to read registration response: %s", err)
	}
	util.Log(fmt.Sprintf("✅ Received %d bytes from server", n))
	util.Log(fmt.Sprintf("📄 Raw response: %s", string(buffer[:n])))

	var response registrationResponse
	if err := json.Unmarshal(buffer[:n], &response); err != nil {
		util.Log(fmt.Sprintf("❌ Failed to parse JSON response: %v", err))
		util.Log(fmt.Sprintf("📄 Attempting to parse: %s", string(buffer[:n])))
		return "", fmt.Errorf("failed to parse registration response: %s", err)
	}
	util.Log(fmt.Sprintf("✅ Parsed response: Success=%t, AssignedID=%s, Username=%s, Message=%s",
		response.Success, response.AssignedID, response.Username, response.Message))

	if !response.Success {
		util.Log(fmt.Sprintf("❌ Registration failed: %s", response.Message))
		return "", fmt.Errorf("registration failed: %s", response.Message)
	}

	util.Log(fmt.Sprintf("✅ Registration successful. Assigned UUID: %s", response.AssignedID))

	return response.AssignedID, nil
}

// SendMessage we'll see what to do with this
func (c *Quic) SendMessage(bytes []byte) error {
	stream, err := c.conn.OpenStreamSync(context.Background())
	if err != nil {
		return fmt.Errorf("failed to open stream:\n%s", err)
	}

	n, err := stream.Write(bytes)
	if err != nil {
		return fmt.Errorf("failed to send message:\n%s", err)
	}

	util.Log(fmt.Sprintf("✅ Sent %d bytes", n))
	return nil
}

func (c *Quic) StartMessageListener() {
	go func() {
		for {
			stream, err := c.conn.AcceptStream(context.Background())
			if err != nil {
				util.Log(fmt.Sprintf("Error accepting stream: %v", err))
				return
			}
			go handleIncomingStream(stream)
		}
	}()
}

func handleIncomingStream(stream *quic.Stream) {
	defer stream.Close()

	buffer := make([]byte, 1024)
	n, err := stream.Read(buffer)
	if err != nil {
		util.Log(fmt.Sprintf("Error reading incoming message: %v", err))
		return
	}

	var serverChatMsg ChatMessage
	if err := json.Unmarshal(buffer[:n], &serverChatMsg); err == nil {
		msg := ChatMessage{
			Message: serverChatMsg.Message,
			To:      serverChatMsg.To,
			From:    serverChatMsg.From,
		}

		select {
		case messageChannel <- msg:
			util.Log(fmt.Sprintf("📨 Received structured message from %s: %s", msg.From, msg.Message))
		default:
			util.Log("Message channel full, dropping message")
		}
	}

	util.Log(fmt.Sprintf("Received response: %s", string(buffer[:n])))
}

func GetMessageChannel() <-chan ChatMessage {
	return messageChannel
}
