package main

import (
	"context"
	"encoding/json"
	"log"

	"github.com/quic-go/quic-go"
)

type messageBody struct {
	Message string
}

func sendMessage(conn *quic.Conn, msg messageBody) {
	stream, err := conn.OpenStreamSync(context.Background())
	if err != nil {
		log.Printf("Failed to open stream:\n%s", err)
	}
	defer stream.Close()

	bytes, err := json.Marshal(msg)
	if err != nil {
		log.Printf("%s", err)
	}

	payload := string(bytes)
	stream.Write([]byte(payload))
}
