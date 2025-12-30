package handlers

import (
	"clingy-client/services"
	"clingy-client/util"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type ChatHandler struct {
	quicService   *services.Quic
	configService *services.Config
}

func NewChatHandler(
	quicService *services.Quic,
	configService *services.Config,
) *ChatHandler {
	return &ChatHandler{
		quicService:   quicService,
		configService: configService,
	}
}

type SendChatBody struct {
	To      string `json:"to"`
	Message string `json:"message"`
}

func (h *ChatHandler) SendChatMessage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	log.Printf("API: POST /api/chat - sending message")

	var body SendChatBody
	if !decodeJSONBody(w, r, &body) {
		return
	}

	chatMessage := services.ChatMessage{
		To:      body.To,
		Message: body.Message,
		From:    h.configService.Username,
	}
	msgBytes, _ := json.Marshal(chatMessage)
	h.quicService.SendMessage(msgBytes)
}

func (h *ChatHandler) GetMessageStream(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("X-Accel-Buffering", "no")
	w.WriteHeader(http.StatusOK)

	// Establishing initial connection
	if flusher, ok := w.(http.Flusher); ok {
		flusher.Flush()
	}

	messageChan := services.GetMessageChannel()
	util.Log("🔥 STREAM: Got message channel, starting loop")

	for {
		select {
		case msg, ok := <-messageChan:
			if !ok {
				util.Log("🔥 STREAM: Channel closed")
				return
			}
			util.Log(fmt.Sprintf("🔥 STREAM: Received message from channel - From: %s, To: %s, Message: %s", msg.From, msg.To, msg.Message))
			data, err := json.Marshal(msg)
			if err != nil {
				util.Log(fmt.Sprintf("🔥 STREAM: Error marshaling message: %v", err))
				continue
			}
			_, err = fmt.Fprintf(w, "data: %s\n\n", string(data))
			if err != nil {
				util.Log(fmt.Sprintf("🔥 STREAM: Error writing to response: %v", err))
				return
			}
			w.(http.Flusher).Flush()
			util.Log(fmt.Sprintf("🔥 STREAM: ✅ Successfully wrote to stream: %s", string(data)))
		case <-r.Context().Done():
			util.Log("🔥 STREAM: Client disconnected")
			return
		}
	}
}
