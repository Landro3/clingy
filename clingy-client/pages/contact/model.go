package contact

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type Contact struct {
	userId   string
	username string
}

type Model struct {
	list     list.Model
	contacts []list.Item
	ready    bool
}

func InitialModel() Model {
	items := []list.Item{
		item{title: "Raspberry Pi’s", desc: "I have ’em all over my house"},
		item{title: "Nutella", desc: "It's good on toast"},
	}

	m := Model{
		list:     list.New(items, list.NewDefaultDelegate(), 0, 0),
		ready:    false,
		contacts: items,
	}
	m.list.Title = "My Fave Things"
	return m
}

func (m Model) Init() tea.Cmd {
	return nil
}
