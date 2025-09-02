package http3

import (
	"clingy-client/services"
	"clingy-client/util"
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	quic "github.com/quic-go/quic-go"
)

var connection *quic.Conn

type registerMessage struct {
	Username string
	ID       string
}

type ChatMessage struct {
	To      string
	From    string
	Message string
}

var messageChannel = make(chan ChatMessage, 100)

func ConnectToServer(serverAddr string) error {
	tlsConfig := &tls.Config{
		MinVersion:         tls.VersionTLS13,
		NextProtos:         []string{"clingy-v1"},
		InsecureSkipVerify: true, // Only for self-signed certs
	}

	conn, err := quic.DialAddr(context.Background(), serverAddr, tlsConfig, nil)
	if err != nil {
		return fmt.Errorf("failed to connect:\n%s", err)
	}

	connection = conn

	select {
	case <-conn.HandshakeComplete():
		util.Log("✅ Handshake completed")
		go keepConnectionAlive()
	case <-time.After(5 * time.Second):
		util.Log("❌ Handshake timeout")
		return errors.New("could not connect to server")
	}

	return nil
}

func RegisterInServer(config *services.Config) {
	err := ConnectToServer(config.ServerAddr)
	if err != nil {
		util.Log(err.Error())
		return
	}

	message := registerMessage{Username: config.Username, ID: config.UniqueID}
	bytes, _ := json.Marshal(message)
	err = SendMessage(bytes)
	if err != nil {
		util.Log(err.Error())
	}

	config.UpdateConfig(config)

	StartMessageListener()
}

func SendMessage(bytes []byte) error {
	stream, err := connection.OpenStreamSync(context.Background())
	if err != nil {
		return fmt.Errorf("failed to open stream:\n%s", err)
	}
	defer stream.Close()

	n, err := stream.Write(bytes)
	if err != nil {
		return fmt.Errorf("failed to send message:\n%s", err)
	}

	util.Log(fmt.Sprintf("✅ Sent %d bytes", n))
	return nil
}

func StartMessageListener() {
	go func() {
		for {
			stream, err := connection.AcceptStream(context.Background())
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

func keepConnectionAlive() error {
	ticker := time.NewTicker(20 * time.Second)
	defer ticker.Stop()

	ctx := context.Background()

	for {
		select {
		case <-ticker.C:
			err := SendMessage([]byte("ping"))
			if err != nil {
				return fmt.Errorf("failed to write keep-alive ping: %w", err)
			}

		case <-ctx.Done():
			return ctx.Err()
		}
	}
}
