package app

import (
	"clingy-client/pages/chat"

	tea "github.com/charmbracelet/bubbletea"
)

type PageType int

const (
	ChatPage PageType = iota
	ContactPage
	ConnectPage
)

type Model struct {
	currentPage PageType
	chatModel   chat.Model
}

func (m Model) Init() tea.Cmd { return nil }

func InitialModel() Model {
	return Model{
		currentPage: ChatPage,
		chatModel:   chat.InitialModel(),
	}
}
