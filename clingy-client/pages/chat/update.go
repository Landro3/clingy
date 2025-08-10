package chat

import (
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
		// if k := msg.String(); k == "ctrl+c" || k == "q" || k == "esc" {
		// 	return m, tea.Quit
		// }

		if m.chatInput.Focused() {
			m.chatInput, cmd = m.chatInput.Update(msg)
			return m, cmd
		}

	case tea.WindowSizeMsg:
		headerHeight := lipgloss.Height(m.headerView())
		footerHeight := lipgloss.Height(m.footerView())
		verticalMarginHeight := headerHeight + footerHeight

		if !m.ready {
			// Since this program is using the full size of the viewport we
			// need to wait until we've received the window dimensions before
			// we can initialize the viewport. The initial dimensions come in
			// quickly, though asynchronously, which is why we wait for them
			// here.
			m.viewport = viewport.New(msg.Width, msg.Height-lipgloss.Height(m.chatInputView())-verticalMarginHeight)
			m.viewport.Style = lipgloss.NewStyle().Border(lipgloss.RoundedBorder())
			m.viewport.YPosition = headerHeight

			m.content = m.formatMessages()
			m.viewport.SetContent(m.content)

			// Chat input sizing
			if msg.Width < 100 {
				m.chatInput.Width = msg.Width - 12
			} else {
				m.chatInput.Width = msg.Width / 2
			}
			m.ready = true
		} else {
			// Viewport sizing
			m.viewport.Width = msg.Width
			m.viewport.Height = msg.Height - verticalMarginHeight

			// Add messages (m.content) to viewport
			m.content = m.formatMessages()
			m.viewport.SetContent(m.content)

			// Chat input sizing
			if msg.Width < 100 {
				m.chatInput.Width = msg.Width - 12
			} else {
				m.chatInput.Width = msg.Width / 2
			}
		}
	}

	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}
