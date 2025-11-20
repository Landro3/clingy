package app

import (
	"clingy-client/http3"
	"clingy-client/pages/chat"
	"clingy-client/pages/config"
	"clingy-client/pages/contact"
	"clingy-client/services"
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

func (m Model) Init() tea.Cmd {
	switch m.currentPage {
	case shared.ChatPage:
		return m.chatModel.Init()
	case shared.ContactPage:
		return m.contactModel.Init()
	case shared.ConfigPage:
		return m.configModel.Init()
	default:
		return nil
	}
}

func InitialModel() Model {
	configService := services.NewConfig()
	if configService.Username != "" && configService.UniqueID != "" {
		http3.RegisterInServer(configService)
	}

	contactService := services.NewContact(configService)
	return Model{
		currentPage:  shared.ChatPage,
		chatModel:    chat.InitialModel(configService),
		contactModel: contact.InitialModel(contactService),
		configModel:  config.InitialModel(configService),
		windowHeight: 0,
		windowWidth:  0,
	}
}
