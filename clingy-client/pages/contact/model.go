package contact

import (
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

type Model struct {
	list       list.Model
	contacts   []list.Item
	inputs     []textinput.Model
	focusIndex int
	mode       ModeType
	ready      bool
}

func InitialModel() Model {
	items := []list.Item{
		user{name: "John Nathan", id: "fjw38fhfn9v87fga9jifnga0"},
		user{name: "Alex Zander", id: "0q9fhva73098ang0b9uahfa"},
	}

	keys := list.DefaultKeyMap()
	keys.GoToStart.SetEnabled(false)
	keys.GoToEnd.SetEnabled(false)
	keys.NextPage.SetEnabled(false)
	keys.PrevPage.SetEnabled(false)
	keys.Quit.SetEnabled(false)

	contactList := list.New(items, list.NewDefaultDelegate(), 0, 0)
	contactList.KeyMap = keys
	contactList.SetShowHelp(false)

	usernameInput := usernameInputModel()
	uuidInput := uuidInputModel()
	inputs := []textinput.Model{usernameInput, uuidInput}

	m := Model{
		list:     contactList,
		ready:    false,
		contacts: items,
		inputs:   inputs,
		mode:     None,
	}

	m.list.Title = "Contacts"
	return m
}

func (m Model) Init() tea.Cmd {
	return nil
}
