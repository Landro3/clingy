package contact

import (
	"clingy-client/services"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type ModeType int

const (
	None ModeType = iota
	Add
	Edit
	Remove
)

type SelectChatMsg struct {
	Chat services.ContactInfo
}

func SelectChatCmd(contactInfo services.ContactInfo) func() tea.Msg {
	return func() tea.Msg {
		return SelectChatMsg{contactInfo}
	}
}

type Model struct {
	list       list.Model
	inputs     []textinput.Model
	focusIndex int
	mode       ModeType
	ready      bool
	contact    *services.Contact
}

func InitialModel(contactService *services.Contact) Model {
	keys := list.DefaultKeyMap()
	keys.GoToStart.SetEnabled(false)
	keys.GoToEnd.SetEnabled(false)
	keys.NextPage.SetEnabled(false)
	keys.PrevPage.SetEnabled(false)
	keys.Quit.SetEnabled(false)

	listItems := contactService.ToListItems()
	contactList := list.New(listItems, list.NewDefaultDelegate(), 0, 0)
	contactList.KeyMap = keys
	contactList.SetShowHelp(false)

	usernameInput := usernameInputModel()
	uuidInput := uuidInputModel()
	inputs := []textinput.Model{usernameInput, uuidInput}

	m := Model{
		list:    contactList,
		ready:   false,
		inputs:  inputs,
		mode:    None,
		contact: contactService,
	}

	m.list.Title = "Contacts"
	return m
}

func (m Model) Init() tea.Cmd {
	return nil
}
