package chat

import (
	"clingy-client/http3"
	"clingy-client/services"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

type FocusedElement int

const (
	ChatInput FocusedElement = iota
	Contact
	Config
	Max
)

type Model struct {
	ready       bool
	viewport    viewport.Model
	messages    []http3.ChatMessage
	chatInput   textinput.Model
	focus       FocusedElement
	currentChat *services.ContactInfo
	config      *services.Config
}

func InitialModel(configService *services.Config) Model {
	ti := textinput.New()
	ti.Placeholder = "The message will go here..."
	ti.Focus()

	messages := []http3.ChatMessage{}

	m := Model{
		messages:    messages,
		chatInput:   ti,
		ready:       false,
		currentChat: nil,
		config:      configService,
	}

	return m
}

type IncomingMessageMsg struct {
	Message http3.ChatMessage
}

type NoMessageMsg struct{}

func WaitForMessage() tea.Cmd {
	return func() tea.Msg {
		select {
		case msg := <-http3.GetMessageChannel():
			return IncomingMessageMsg{Message: msg}
		case <-time.After(1000 * time.Millisecond):
			return NoMessageMsg{}
		}
	}
}

func (m Model) Init() tea.Cmd {
	return WaitForMessage()
}
