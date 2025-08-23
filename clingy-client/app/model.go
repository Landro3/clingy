package app

import (
	"clingy-client/http3"
	"clingy-client/pages/chat"
	"clingy-client/pages/config"
	"clingy-client/pages/contact"
	"clingy-client/shared"

	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	currentPage  shared.PageType
	chatModel    chat.Model
	contactModel contact.Model
	configModel  config.Model
	windowHeight int
	windowWidth  int
}

func (m Model) Init() tea.Cmd { return nil }

func InitialModel() Model {
	http3.ConnectToServer()
	return Model{
		currentPage:  shared.ChatPage,
		chatModel:    chat.InitialModel(),
		contactModel: contact.InitialModel(),
		configModel:  config.InitialModel(),
		windowHeight: 0,
		windowWidth:  0,
	}
}
