package services

import (
	"bufio"
	"bytes"
	"clingy-client/util"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/quic-go/quic-go"
	"github.com/quic-go/quic-go/http3"
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
	To      string `json:"to"`
	From    string `json:"from"`
	Message string `json:"message"`
}

type Http3 struct {
	client    *http.Client
	sseCancel func()
	config    *Config
	sendMessage func(ChatMessage)
}

func NewHttp3(config *Config, sendMessage func(ChatMessage)) *Http3 {
	tlsConfig := &tls.Config{
		MinVersion:         tls.VersionTLS13,
		NextProtos:         []string{"h3"},
		InsecureSkipVerify: true, // Only for self-signed certs
	}

	transport := &http3.Transport{
		TLSClientConfig: tlsConfig,
		QUICConfig: &quic.Config{
			KeepAlivePeriod: 15 * time.Minute, // More frequent keep-alives
			MaxIdleTimeout:  10 * time.Minute, // Allow long-lived SSE connections
		},
	}

	client := &http.Client{
		Transport: transport,
		Timeout:   30 * time.Second,
	}

	return &Http3{
		client: client,
		config: config,
		sendMessage: sendMessage,
	}
}

func (h *Http3) Register(username string) (string, error) {
	util.Log(fmt.Sprintf("🔄 Starting HTTP/3 registration for user: %s", username))

	message := registerMessage{Username: username}
	jsonData, err := json.Marshal(message)
	if err != nil {
		return "", fmt.Errorf("failed to marshal registration message: %s", err)
	}

	url := h.config.ServerAddr + "/register"
	util.Log(fmt.Sprintf("📤 Sending POST to: %s", url))

	// Create client without timeout for SSE connection
	transport := h.client.Transport
	sseClient := &http.Client{
		Transport: transport,
		// No timeout for persistent SSE connection
	}

	resp, err := sseClient.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to send registration request: %s", err)
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return "", fmt.Errorf("registration failed with status: %d", resp.StatusCode)
	}

	util.Log("✅ Registration request sent, starting SSE connection...")

	// Start processing SSE stream in background
	go h.establishSSE(resp)

	return username, nil
}

func (h *Http3) SendMessage(bytes []byte) error {
	var chatMsg ChatMessage
	if err := json.Unmarshal(bytes, &chatMsg); err != nil {
		return fmt.Errorf("failed to unmarshal chat message: %s", err)
	}

	url := h.config.ServerAddr + "/chat"
	util.Log(fmt.Sprintf("📤 Sending chat message to: %s", url))

	resp, err := h.client.Post(url, "application/json", strings.NewReader(string(bytes)))
	if err != nil {
		return fmt.Errorf("failed to send chat message: %s", err)
	}
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			util.Log(fmt.Sprintf("error closing response body: %s", err))
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("chat message failed with status: %d", resp.StatusCode)
	}

	util.Log("✅ Chat message sent successfully")
	return nil
}

func (h *Http3) establishSSE(resp *http.Response) {
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			util.Log(fmt.Sprintf("error closing response body: %s", err))
		}
	}()

	util.Log("🔗 Setting up SSE stream")

	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasPrefix(line, "data: ") {
			data := strings.TrimPrefix(line, "data: ")

			// Try to parse as chat message
			var chatMsg ChatMessage
			if err := json.Unmarshal([]byte(data), &chatMsg); err == nil && chatMsg.From != "" {
				h.sendMessage(chatMsg)
			} else {
				// Log other events (like registration confirmation)
				util.Log(fmt.Sprintf("📄 SSE event: %s", data))
			}
		}
	}

	if err := scanner.Err(); err != nil {
		util.Log(fmt.Sprintf("❌ SSE connection error: %v", err))

		// Simple retry for H3_MESSAGE_ERROR
		if strings.Contains(err.Error(), "H3_MESSAGE_ERROR") {
			time.Sleep(5 * time.Second)
			go h.Register(h.config.Username)
		}
	}

	util.Log("🔌 SSE connection closed")
}

