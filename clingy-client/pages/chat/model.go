package chat

import (
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

type ChatMessage struct {
	Author  string
	Content string
	Time    string
	IsOwn   bool
}

type Model struct {
	content   string
	ready     bool
	viewport  viewport.Model
	messages  []ChatMessage
	chatInput textinput.Model
}

func InitialModel() Model {
	ti := textinput.New()
	ti.Placeholder = "The message will go here..."
	ti.Focus()

	messages := []ChatMessage{
		{Author: "Alice", Content: "Hello everyone!", Time: "14:30", IsOwn: false},
		{Author: "You", Content: "Hey Alice!", Time: "14:31", IsOwn: true},
		{Author: "Bob", Content: "Hey Alice, how's it going?", Time: "14:32", IsOwn: false},
		{Author: "You", Content: "Pretty good, thanks!", Time: "14:33", IsOwn: true},
		{Author: "Alice", Content: "Working on some Go code.", Time: "14:34", IsOwn: false},
		{Author: "You", Content: "Nice! What kind of project?", Time: "14:35", IsOwn: true},
	}

	m := Model{
		messages:  messages,
		chatInput: ti,
		ready:     false,
	}

	return m
}

func (m Model) Init() tea.Cmd {
	return nil
}
