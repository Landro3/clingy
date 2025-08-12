package chat

import (
	"clingy-client/shared"
	"clingy-client/util"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "tab" {
			switch m.focus {
			case ChatInput:
				m.chatInput.Blur()
				m.focus = Contact
			case Contact:
				m.chatInput.Focus()
				m.focus = ChatInput
			}
		}

		if msg.String() == "enter" {
			switch m.focus {
			case ChatInput:
				util.Log("...sending message")
			case Contact:
				navigateCmd := func() tea.Msg {
					return shared.NavigateMsg{Page: shared.ContactPage}
				}
				cmds = append(cmds, navigateCmd)
			}
		}

		m.chatInput, cmd = m.chatInput.Update(msg)
		cmds = append(cmds, cmd)
	case tea.WindowSizeMsg:
		headerHeight := lipgloss.Height(m.headerView())
		footerHeight := lipgloss.Height(m.footerView())
		verticalMarginHeight := headerHeight + footerHeight
		chatInputHeight := lipgloss.Height(m.chatInputView())

		m.resizeChatInput(msg.Width)
		if !m.ready {
			// Since this program is using the full size of the viewport we
			// need to wait until we've received the window dimensions before
			// we can initialize the viewport. The initial dimensions come in
			// quickly, though asynchronously, which is why we wait for them
			// here.
			m.viewport = viewport.New(msg.Width, msg.Height-chatInputHeight-verticalMarginHeight)
			m.viewport.Style = lipgloss.NewStyle().Border(lipgloss.RoundedBorder())
			m.viewport.YPosition = headerHeight

			m.ready = true
		} else {
			// Viewport sizing
			m.viewport.Width = msg.Width
			m.viewport.Height = msg.Height - verticalMarginHeight - chatInputHeight
		}

		m.setChatContent()
	}

	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m *Model) resizeChatInput(msgWidth int) {
	if msgWidth < 100 {
		m.chatInput.Width = msgWidth - 12
	} else {
		m.chatInput.Width = msgWidth / 2
	}
}

func (m *Model) setChatContent() {
	m.content = m.formatMessages()
	m.viewport.SetContent(m.content)
}
