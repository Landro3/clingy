package shared

import (
	tea "github.com/charmbracelet/bubbletea"
)

type PageType int

const (
	ChatPage PageType = iota
	ContactPage
	ConnectPage
)

type NavigateMsg struct {
	Page PageType
}

func NavigateCmd(page PageType) func() tea.Msg {
	return func() tea.Msg {
		return NavigateMsg{Page: page}
	}
}
