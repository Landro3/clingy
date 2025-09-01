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

var (
	connection  *quic.Conn
	isConnected bool
)

func IsConnected() bool { return isConnected }

// func IsRegistered() bool { return isRegistered }

type registerMessage struct {
	Username string
	ID       string
}

type IncomingMessage struct {
	Content string
	Author  string
	Time    string
}

var messageChannel = make(chan IncomingMessage, 100)

func ConnectToServer(serverAddr string) error {
	tlsConfig := &tls.Config{
		MinVersion:         tls.VersionTLS13,
		NextProtos:         []string{"clingy-v1"},
		InsecureSkipVerify: true, // Only for self-signed certs
	}

	// TODO: Config value
	conn, err := quic.DialAddr(context.Background(), serverAddr, tlsConfig, nil)
	if err != nil {
		return fmt.Errorf("failed to connect:\n%s", err)
	}

	connection = conn

	select {
	case <-conn.HandshakeComplete():
		util.Log("✅ Handshake completed")
		isConnected = true
	case <-time.After(5 * time.Second):
		util.Log("❌ Handshake timeout")
		isConnected = false
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
	bytes, err := json.Marshal(message)
	if err != nil {
		util.Log(fmt.Sprintf("%s", err))
	}

	payload := string(bytes)
	err = SendMessage(payload)
	if err != nil {
		// TODO: registered?
		isConnected = false
		util.Log(err.Error())
	}

	config.UpdateConfig(config)

	StartMessageListener()
}

func SendMessage(msg string) error {
	stream, err := connection.OpenStreamSync(context.Background())
	if err != nil {
		return fmt.Errorf("failed to open stream:\n%s", err)
	}
	defer stream.Close()

	n, err := stream.Write([]byte(msg))
	if err != nil {
		return fmt.Errorf("failed to send message:\n%s", err)
	}

	util.Log(fmt.Sprintf("✅ Sent %d bytes: %s", n, msg))
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

	msg := IncomingMessage{
		Content: string(buffer[:n]),
		Author:  "Server",
		Time:    time.Now().Format("15:04"),
	}

	select {
	case messageChannel <- msg:
		util.Log(fmt.Sprintf("📨 Received message: %s", msg.Content))
	default:
		util.Log("Message channel full, dropping message")
	}
}

func GetMessageChannel() <-chan IncomingMessage {
	return messageChannel
}
