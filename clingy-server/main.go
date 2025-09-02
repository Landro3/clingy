package main

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"log"

	quic "github.com/quic-go/quic-go"
)

var connMap *ConnectionMap

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
	ID       string
}

type chatMessage struct {
	To      string
	From    string
	Message string
}

func handleStream(stream *quic.Stream, conn *quic.Conn) {
	defer stream.Close()

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
		log.Printf("Registered\nUsername: %s\nID: %s", regMsg.Username, regMsg.ID)
		connMap.Add(regMsg.ID, conn)
		return
	}

	var chatMsg chatMessage
	if err := json.Unmarshal(buffer[:n], &chatMsg); err == nil && chatMsg.To != "" {
		log.Printf("Chat: %s", chatMsg.Message)
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
		return
	}

	connMap.LogConnections()

	log.Printf("Message: %s", buffer[:n])
}
